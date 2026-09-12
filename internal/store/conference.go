package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/HazelnutParadise/briefast/internal/report"
)

// Settlement holds the two closes a brief is judged against. Outcome is
// computed by the store from the stored prediction, never supplied by callers.
type Settlement struct {
	PreCloseDate  string
	PreClose      float64
	PostCloseDate string
	PostClose     float64
	Outcome       string
	SettledAt     time.Time
}

// ConferenceSummary carries the indexed columns only, enough for listings and
// the accuracy tally without decoding the brief payload.
type ConferenceSummary struct {
	Symbol      string
	Name        string
	Market      string
	HeldOn      string
	HeldAt      string
	Venue       string
	AnnouncedOn string
	Prediction  string
	Settlement  *Settlement
	UpdatedAt   time.Time
}

// Conference is a summary plus the decoded brief.
type Conference struct {
	ConferenceSummary
	Brief report.Conference
}

const conferenceColumns = `symbol, name, market, held_on, held_at, venue, announced_on, prediction,
        pre_close_date, pre_close, post_close_date, post_close, outcome, settled_at, updated_at`

func (s *Store) UpsertConference(ctx context.Context, c report.Conference, payload []byte) error {
	return upsertConference(ctx, s.db, c, payload, formatTime(s.now()))
}

// upsertConference rewrites the brief columns only; settlement columns are
// left alone so re-posting a brief never erases a recorded outcome.
func upsertConference(ctx context.Context, exec interface {
	ExecContext(context.Context, string, ...any) (sql.Result, error)
}, c report.Conference, payload []byte, now string) error {
	_, err := exec.ExecContext(ctx, `INSERT INTO conferences(symbol, held_on, name, market, held_at, venue, announced_on, prediction, payload, created_at, updated_at)
        VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(symbol, held_on) DO UPDATE SET
            name = excluded.name,
            market = excluded.market,
            held_at = excluded.held_at,
            venue = excluded.venue,
            announced_on = excluded.announced_on,
            prediction = excluded.prediction,
            payload = excluded.payload,
            updated_at = excluded.updated_at`,
		c.Symbol, c.HeldOn, c.Name, c.Market, c.HeldAt, c.Venue, c.AnnouncedOn, c.Prediction, payload, now, now)
	if err != nil {
		return fmt.Errorf("upsert conference: %w", err)
	}
	return nil
}

// IngestConference stores the brief and its update-log row in one transaction.
func (s *Store) IngestConference(ctx context.Context, c report.Conference, payload []byte, key APIKey) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin conference ingest: %w", err)
	}
	defer tx.Rollback()
	if err := upsertConference(ctx, tx, c, payload, formatTime(s.now())); err != nil {
		return err
	}
	if err := addUpdateLog(ctx, tx, UpdateLog{
		At: s.now(), APIKeyID: &key.ID, KeyName: key.Name,
		ReportDate: c.HeldOn, Action: "conference_ok", Detail: c.Symbol,
	}); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit conference ingest: %w", err)
	}
	return nil
}

// SettleConference records both closes, derives the outcome from the stored
// prediction, and logs the settlement, all in one transaction. It returns
// ErrNotFound when no brief exists for the pair.
func (s *Store) SettleConference(ctx context.Context, symbol, heldOn string, st Settlement, key APIKey) (string, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("begin settlement: %w", err)
	}
	defer tx.Rollback()
	var prediction string
	if err := tx.QueryRowContext(ctx, "SELECT prediction FROM conferences WHERE symbol = ? AND held_on = ?", symbol, heldOn).Scan(&prediction); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("read conference prediction: %w", err)
	}
	outcome := report.Outcome(prediction, st.PreClose, st.PostClose)
	now := s.now()
	if _, err := tx.ExecContext(ctx, `UPDATE conferences SET pre_close_date = ?, pre_close = ?, post_close_date = ?, post_close = ?,
        outcome = ?, settled_at = ?, updated_at = ? WHERE symbol = ? AND held_on = ?`,
		st.PreCloseDate, st.PreClose, st.PostCloseDate, st.PostClose, outcome, formatTime(now), formatTime(now), symbol, heldOn); err != nil {
		return "", fmt.Errorf("settle conference: %w", err)
	}
	if err := addUpdateLog(ctx, tx, UpdateLog{
		At: now, APIKeyID: &key.ID, KeyName: key.Name,
		ReportDate: heldOn, Action: "settle_ok", Detail: symbol + " " + outcome,
	}); err != nil {
		return "", err
	}
	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit settlement: %w", err)
	}
	return outcome, nil
}

func (s *Store) ConferenceByKey(ctx context.Context, symbol, heldOn string) (*Conference, error) {
	row := s.db.QueryRowContext(ctx, "SELECT "+conferenceColumns+", payload FROM conferences WHERE symbol = ? AND held_on = ?", symbol, heldOn)
	var out Conference
	var payload []byte
	if err := scanConferenceSummary(row, &out.ConferenceSummary, &payload); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(payload, &out.Brief); err != nil {
		return nil, fmt.Errorf("decode stored conference: %w", err)
	}
	return &out, nil
}

// ListUpcomingConferences returns briefs held on or after today, soonest first.
func (s *Store) ListUpcomingConferences(ctx context.Context, today string) ([]ConferenceSummary, error) {
	return s.listConferences(ctx, "WHERE held_on >= ? ORDER BY held_on ASC, symbol ASC", today)
}

// ListPendingSettlement returns past briefs that have no settlement yet.
func (s *Store) ListPendingSettlement(ctx context.Context, today string) ([]ConferenceSummary, error) {
	return s.listConferences(ctx, "WHERE held_on < ? AND settled_at IS NULL ORDER BY held_on ASC, symbol ASC", today)
}

// ListSettledConferences returns the most recently held settled briefs.
func (s *Store) ListSettledConferences(ctx context.Context, limit int) ([]ConferenceSummary, error) {
	if limit < 1 {
		limit = 20
	}
	return s.listConferences(ctx, "WHERE settled_at IS NOT NULL ORDER BY held_on DESC, symbol ASC LIMIT ?", limit)
}

// ListConferences returns every stored brief, newest event first.
func (s *Store) ListConferences(ctx context.Context) ([]ConferenceSummary, error) {
	return s.listConferences(ctx, "ORDER BY held_on DESC, symbol ASC")
}

func (s *Store) listConferences(ctx context.Context, clause string, args ...any) ([]ConferenceSummary, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT "+conferenceColumns+" FROM conferences "+clause, args...)
	if err != nil {
		return nil, fmt.Errorf("list conferences: %w", err)
	}
	defer rows.Close()
	var out []ConferenceSummary
	for rows.Next() {
		var item ConferenceSummary
		if err := scanConferenceSummary(rows, &item, nil); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

// scanConferenceSummary reads the indexed columns; when payload is non-nil the
// query also selected the payload column and it is filled in.
func scanConferenceSummary(row scanner, item *ConferenceSummary, payload *[]byte) error {
	var preDate, postDate, outcome, settledAt sql.NullString
	var preClose, postClose sql.NullFloat64
	var updated string
	dest := []any{&item.Symbol, &item.Name, &item.Market, &item.HeldOn, &item.HeldAt, &item.Venue, &item.AnnouncedOn, &item.Prediction,
		&preDate, &preClose, &postDate, &postClose, &outcome, &settledAt, &updated}
	if payload != nil {
		dest = append(dest, payload)
	}
	if err := row.Scan(dest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotFound
		}
		return fmt.Errorf("scan conference: %w", err)
	}
	var err error
	if item.UpdatedAt, err = parseTime(updated); err != nil {
		return err
	}
	if settledAt.Valid {
		at, err := parseTime(settledAt.String)
		if err != nil {
			return err
		}
		item.Settlement = &Settlement{
			PreCloseDate: preDate.String, PreClose: preClose.Float64,
			PostCloseDate: postDate.String, PostClose: postClose.Float64,
			Outcome: outcome.String, SettledAt: at,
		}
	}
	return nil
}
