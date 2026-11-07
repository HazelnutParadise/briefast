package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	SeenArticleTTL = 90 * 24 * time.Hour
	SeenEventTTL   = 30 * 24 * time.Hour
)

// SeenArticle is one processed article in the daily-brief read log.
type SeenArticle struct {
	SourceID     string
	URL          string
	URLHash      string
	Title        string
	TitleHash    string
	ContentHash  string
	LeadHash     string
	Sketch       []string
	SourceDomain string
	Symbols      []string
	EventKey     string
	PublishedAt  string
	Decision     string
	Call         string
	FirstSeenAt  time.Time
}

// SeenEvent is a cross-day event key such as "2330|月營收|2026-07".
type SeenEvent struct {
	Key         string    `json:"event_key"`
	Note        string    `json:"note"`
	FirstSeenAt time.Time `json:"first_seen_at"`
}

// PruneSeen deletes articles and events past their retention period.
func (s *Store) PruneSeen(ctx context.Context) error {
	now := s.now()
	if _, err := s.db.ExecContext(ctx, "DELETE FROM seen_articles WHERE created_at < ?", formatTime(now.Add(-SeenArticleTTL))); err != nil {
		return fmt.Errorf("prune seen articles: %w", err)
	}
	if _, err := s.db.ExecContext(ctx, "DELETE FROM seen_events WHERE created_at < ?", formatTime(now.Add(-SeenEventTTL))); err != nil {
		return fmt.Errorf("prune seen events: %w", err)
	}
	return nil
}

// SeenMatch reports which fingerprint, if any, matches a stored article,
// checked in the order source ID, URL hash, title hash. Empty values never
// match. It returns "" for an unseen article.
func (s *Store) SeenMatch(ctx context.Context, sourceID, urlHash, titleHash string) (string, error) {
	checks := []struct{ reason, column, value string }{
		{"source_id", "source_id", sourceID},
		{"url", "url_hash", urlHash},
		{"title", "title_hash", titleHash},
	}
	for _, check := range checks {
		if check.value == "" {
			continue
		}
		var found int
		err := s.db.QueryRowContext(ctx, "SELECT 1 FROM seen_articles WHERE "+check.column+" = ? LIMIT 1", check.value).Scan(&found)
		if errors.Is(err, sql.ErrNoRows) {
			continue
		}
		if err != nil {
			return "", fmt.Errorf("match seen article: %w", err)
		}
		return check.reason, nil
	}
	return "", nil
}

// RecordSeenArticle stores the article unless its source ID or URL hash is
// already stored, and reports whether it was added.
func (s *Store) RecordSeenArticle(ctx context.Context, a SeenArticle) (bool, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("begin record seen: %w", err)
	}
	defer tx.Rollback()
	var exists int
	err = tx.QueryRowContext(ctx, `SELECT 1 FROM seen_articles
        WHERE (? <> '' AND source_id = ?) OR (? <> '' AND url_hash = ?) LIMIT 1`,
		a.SourceID, a.SourceID, a.URLHash, a.URLHash).Scan(&exists)
	if err == nil {
		return false, nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return false, fmt.Errorf("check seen article: %w", err)
	}
	sketch, err := json.Marshal(nonNil(a.Sketch))
	if err != nil {
		return false, err
	}
	symbols, err := json.Marshal(nonNil(a.Symbols))
	if err != nil {
		return false, err
	}
	now := formatTime(s.now())
	if _, err := tx.ExecContext(ctx, `INSERT INTO seen_articles(
        source_id, url, url_hash, title, title_hash, content_hash, lead_hash, sketch,
        source_domain, symbols, event_key, published_at, decision, stock_call, created_at, updated_at)
        VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.SourceID, a.URL, a.URLHash, a.Title, a.TitleHash, a.ContentHash, a.LeadHash, string(sketch),
		a.SourceDomain, string(symbols), a.EventKey, a.PublishedAt, a.Decision, a.Call, now, now); err != nil {
		return false, fmt.Errorf("insert seen article: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("commit seen article: %w", err)
	}
	return true, nil
}

// SeenSketches returns every stored article that has a content sketch.
func (s *Store) SeenSketches(ctx context.Context) ([]SeenArticle, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT title, url, event_key, sketch FROM seen_articles WHERE sketch <> '[]'")
	if err != nil {
		return nil, fmt.Errorf("list seen sketches: %w", err)
	}
	defer rows.Close()
	var out []SeenArticle
	for rows.Next() {
		var a SeenArticle
		var sketch string
		if err := rows.Scan(&a.Title, &a.URL, &a.EventKey, &sketch); err != nil {
			return nil, fmt.Errorf("scan seen sketch: %w", err)
		}
		if err := json.Unmarshal([]byte(sketch), &a.Sketch); err != nil {
			return nil, fmt.Errorf("decode seen sketch: %w", err)
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// AddSeenEvent registers an event key and reports whether it was new; an
// existing key is left unchanged.
func (s *Store) AddSeenEvent(ctx context.Context, key, note string) (bool, error) {
	now := formatTime(s.now())
	result, err := s.db.ExecContext(ctx, `INSERT INTO seen_events(event_key, note, created_at, updated_at)
        VALUES(?, ?, ?, ?) ON CONFLICT(event_key) DO NOTHING`, key, note, now, now)
	if err != nil {
		return false, fmt.Errorf("add seen event: %w", err)
	}
	n, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("add seen event: %w", err)
	}
	return n == 1, nil
}

// SeenEventByKey returns the registered event or ErrNotFound.
func (s *Store) SeenEventByKey(ctx context.Context, key string) (SeenEvent, error) {
	var event SeenEvent
	var created string
	err := s.db.QueryRowContext(ctx, "SELECT event_key, note, created_at FROM seen_events WHERE event_key = ?", key).Scan(&event.Key, &event.Note, &created)
	if errors.Is(err, sql.ErrNoRows) {
		return SeenEvent{}, ErrNotFound
	}
	if err != nil {
		return SeenEvent{}, fmt.Errorf("read seen event: %w", err)
	}
	event.FirstSeenAt, err = parseTime(created)
	if err != nil {
		return SeenEvent{}, err
	}
	return event, nil
}

func nonNil(values []string) []string {
	if values == nil {
		return []string{}
	}
	return values
}
