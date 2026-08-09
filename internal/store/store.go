package store

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/HazelnutParadise/briefast/internal/report"
	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var migrationV1 string

var ErrNotFound = errors.New("not found")

type Store struct {
	db  *sql.DB
	now func() time.Time
}

type ReportSummary struct {
	Date        string
	Headline    string
	GeneratedAt string
}

type APIKey struct {
	ID        int64
	Name      string
	Token     string
	CreatedAt time.Time
	RevokedAt *time.Time
}

func (k APIKey) Active() bool { return k.RevokedAt == nil }

type UpdateLog struct {
	ID         int64
	At         time.Time
	APIKeyID   *int64
	KeyName    string
	ReportDate string
	Action     string
	Detail     string
}

func Open(path string) (*Store, error) {
	dsn, err := sqliteDSN(path)
	if err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	s := &Store{db: db, now: time.Now}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}
	if err := s.migrate(context.Background()); err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}

func sqliteDSN(path string) (string, error) {
	const pragmas = "_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"
	if path == ":memory:" {
		return "file::memory:?cache=shared&" + pragmas, nil
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve database path: %w", err)
	}
	if err := ensureParent(abs); err != nil {
		return "", err
	}
	u := url.URL{Scheme: "file", Path: filepath.ToSlash(abs)}
	return u.String() + "?" + pragmas, nil
}

func ensureParent(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return fmt.Errorf("create database directory: %w", err)
	}
	return nil
}

func (s *Store) Close() error { return s.db.Close() }

func (s *Store) DB() *sql.DB { return s.db }

func (s *Store) migrate(ctx context.Context) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin migration: %w", err)
	}
	defer tx.Rollback()
	if _, err := tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (
        version INTEGER PRIMARY KEY,
        applied_at TEXT NOT NULL
    )`); err != nil {
		return fmt.Errorf("create schema_migrations: %w", err)
	}
	var applied int
	if err := tx.QueryRowContext(ctx, "SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&applied); err != nil {
		return fmt.Errorf("read schema version: %w", err)
	}
	if applied < 1 {
		if _, err := tx.ExecContext(ctx, migrationV1); err != nil {
			return fmt.Errorf("apply migration 1: %w", err)
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO schema_migrations(version, applied_at) VALUES(1, ?)", formatTime(s.now())); err != nil {
			return fmt.Errorf("record migration 1: %w", err)
		}
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit migration: %w", err)
	}
	return nil
}

func (s *Store) UpsertReport(ctx context.Context, r report.Report, payload []byte) error {
	return upsertReport(ctx, s.db, r, payload, formatTime(s.now()))
}

func upsertReport(ctx context.Context, exec interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}, r report.Report, payload []byte, now string) error {
	_, err := exec.ExecContext(ctx, `INSERT INTO reports(date, headline, payload, generated_at, created_at, updated_at)
        VALUES(?, ?, ?, ?, ?, ?)
        ON CONFLICT(date) DO UPDATE SET
            headline = excluded.headline,
            payload = excluded.payload,
            generated_at = excluded.generated_at,
            updated_at = excluded.updated_at`,
		r.Date, r.Headline, payload, r.GeneratedAt, now, now)
	if err != nil {
		return fmt.Errorf("upsert report: %w", err)
	}
	return nil
}

func (s *Store) IngestReport(ctx context.Context, r report.Report, payload []byte, key APIKey) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin ingest: %w", err)
	}
	defer tx.Rollback()
	now := formatTime(s.now())
	if err := upsertReport(ctx, tx, r, payload, now); err != nil {
		return err
	}
	if err := addUpdateLog(ctx, tx, UpdateLog{
		At: s.now(), APIKeyID: &key.ID, KeyName: key.Name,
		ReportDate: r.Date, Action: "ingest_ok",
	}); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit ingest: %w", err)
	}
	return nil
}

func (s *Store) LatestReport(ctx context.Context) (*report.Report, error) {
	return scanReport(s.db.QueryRowContext(ctx, "SELECT payload FROM reports ORDER BY date DESC LIMIT 1"))
}

func (s *Store) ReportByDate(ctx context.Context, date string) (*report.Report, error) {
	return scanReport(s.db.QueryRowContext(ctx, "SELECT payload FROM reports WHERE date = ?", date))
}

func scanReport(row *sql.Row) (*report.Report, error) {
	var payload []byte
	if err := row.Scan(&payload); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("scan report: %w", err)
	}
	var r report.Report
	if err := json.Unmarshal(payload, &r); err != nil {
		return nil, fmt.Errorf("decode stored report: %w", err)
	}
	return &r, nil
}

func (s *Store) ListReports(ctx context.Context, page, pageSize int) ([]ReportSummary, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	rows, err := s.db.QueryContext(ctx, `SELECT date, headline, generated_at
        FROM reports ORDER BY date DESC LIMIT ? OFFSET ?`, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, fmt.Errorf("list reports: %w", err)
	}
	defer rows.Close()
	var out []ReportSummary
	for rows.Next() {
		var item ReportSummary
		if err := rows.Scan(&item.Date, &item.Headline, &item.GeneratedAt); err != nil {
			return nil, fmt.Errorf("scan report summary: %w", err)
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) CountReports(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM reports").Scan(&count)
	return count, err
}

func (s *Store) CreateAPIKey(ctx context.Context, name, token string) (APIKey, error) {
	name = strings.TrimSpace(name)
	if name == "" || token == "" {
		return APIKey{}, errors.New("key name and token are required")
	}
	now := s.now().UTC()
	res, err := s.db.ExecContext(ctx, "INSERT INTO api_keys(name, token, created_at) VALUES(?, ?, ?)", name, token, formatTime(now))
	if err != nil {
		return APIKey{}, fmt.Errorf("create api key: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return APIKey{}, fmt.Errorf("read api key id: %w", err)
	}
	return APIKey{ID: id, Name: name, Token: token, CreatedAt: now}, nil
}

func (s *Store) APIKeyByToken(ctx context.Context, token string) (APIKey, error) {
	return scanAPIKey(s.db.QueryRowContext(ctx, `SELECT id, name, token, created_at, revoked_at
        FROM api_keys WHERE token = ?`, token))
}

func scanAPIKey(row *sql.Row) (APIKey, error) {
	var key APIKey
	var created string
	var revoked sql.NullString
	if err := row.Scan(&key.ID, &key.Name, &key.Token, &created, &revoked); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return APIKey{}, ErrNotFound
		}
		return APIKey{}, fmt.Errorf("scan api key: %w", err)
	}
	var err error
	key.CreatedAt, err = parseTime(created)
	if err != nil {
		return APIKey{}, err
	}
	if revoked.Valid {
		t, err := parseTime(revoked.String)
		if err != nil {
			return APIKey{}, err
		}
		key.RevokedAt = &t
	}
	return key, nil
}

func (s *Store) ListAPIKeys(ctx context.Context) ([]APIKey, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT id, name, token, created_at, revoked_at
        FROM api_keys ORDER BY created_at DESC, id DESC`)
	if err != nil {
		return nil, fmt.Errorf("list api keys: %w", err)
	}
	defer rows.Close()
	var out []APIKey
	for rows.Next() {
		var key APIKey
		var created string
		var revoked sql.NullString
		if err := rows.Scan(&key.ID, &key.Name, &key.Token, &created, &revoked); err != nil {
			return nil, err
		}
		key.CreatedAt, err = parseTime(created)
		if err != nil {
			return nil, err
		}
		if revoked.Valid {
			t, err := parseTime(revoked.String)
			if err != nil {
				return nil, err
			}
			key.RevokedAt = &t
		}
		out = append(out, key)
	}
	return out, rows.Err()
}

func (s *Store) RevokeAPIKey(ctx context.Context, id int64) error {
	res, err := s.db.ExecContext(ctx, `UPDATE api_keys SET revoked_at = COALESCE(revoked_at, ?) WHERE id = ?`, formatTime(s.now()), id)
	if err != nil {
		return fmt.Errorf("revoke api key: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func (s *Store) AddUpdateLog(ctx context.Context, entry UpdateLog) error {
	if entry.At.IsZero() {
		entry.At = s.now()
	}
	return addUpdateLog(ctx, s.db, entry)
}

func addUpdateLog(ctx context.Context, exec interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}, entry UpdateLog) error {
	_, err := exec.ExecContext(ctx, `INSERT INTO update_log(at, api_key_id, key_name, report_date, action, detail)
        VALUES(?, ?, ?, ?, ?, ?)`, formatTime(entry.At), entry.APIKeyID, entry.KeyName, entry.ReportDate, entry.Action, entry.Detail)
	if err != nil {
		return fmt.Errorf("add update log: %w", err)
	}
	return nil
}

func (s *Store) ListUpdateLogs(ctx context.Context, limit int) ([]UpdateLog, error) {
	if limit < 1 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx, `SELECT id, at, api_key_id, key_name, report_date, action, detail
        FROM update_log ORDER BY at DESC, id DESC LIMIT ?`, limit)
	if err != nil {
		return nil, fmt.Errorf("list update log: %w", err)
	}
	defer rows.Close()
	var out []UpdateLog
	for rows.Next() {
		var item UpdateLog
		var at string
		var keyID sql.NullInt64
		if err := rows.Scan(&item.ID, &at, &keyID, &item.KeyName, &item.ReportDate, &item.Action, &item.Detail); err != nil {
			return nil, err
		}
		item.At, err = parseTime(at)
		if err != nil {
			return nil, err
		}
		if keyID.Valid {
			item.APIKeyID = &keyID.Int64
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

func (s *Store) CountUpdateLogs(ctx context.Context) (int, error) {
	var count int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM update_log").Scan(&count)
	return count, err
}

// storedTimeLayout keeps a fixed width so stored timestamps sort lexicographically,
// which ORDER BY relies on. RFC3339Nano trims trailing zeros and would break that.
const storedTimeLayout = "2006-01-02T15:04:05.000000000Z07:00"

func formatTime(t time.Time) string { return t.UTC().Format(storedTimeLayout) }

func parseTime(value string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse stored time %q: %w", value, err)
	}
	return t, nil
}
