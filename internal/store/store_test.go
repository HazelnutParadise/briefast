package store

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/HazelnutParadise/briefast/internal/report"
)

func openTestStore(t *testing.T) *Store {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "briefast.db"))
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func testReport(date, headline string) report.Report {
	return report.Report{
		Date: date, Headline: headline, OverviewMD: "總覽", WatchMD: "觀察",
		Calls: report.Calls{}, GeneratedAt: date + "T07:50:00+08:00",
	}
}

func putReport(t *testing.T, s *Store, r report.Report) {
	t.Helper()
	payload, err := json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertReport(context.Background(), r, payload); err != nil {
		t.Fatalf("UpsertReport() error = %v", err)
	}
}

func TestMigrationsAndPragmas(t *testing.T) {
	s := openTestStore(t)
	var version int
	if err := s.DB().QueryRow("SELECT MAX(version) FROM schema_migrations").Scan(&version); err != nil || version != 1 {
		t.Fatalf("migration version = %d, err = %v", version, err)
	}
	for name, want := range map[string]int{"foreign_keys": 1, "busy_timeout": 5000} {
		var got int
		if err := s.DB().QueryRow("PRAGMA " + name).Scan(&got); err != nil || got != want {
			t.Errorf("PRAGMA %s = %d, err = %v, want %d", name, got, err, want)
		}
	}
	var mode string
	if err := s.DB().QueryRow("PRAGMA journal_mode").Scan(&mode); err != nil || mode != "wal" {
		t.Errorf("journal_mode = %q, err = %v, want wal", mode, err)
	}
}

func TestReportsEmptyWriteOverwriteAndSort(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	if _, err := s.LatestReport(ctx); !errors.Is(err, ErrNotFound) {
		t.Fatalf("LatestReport() error = %v, want ErrNotFound", err)
	}

	putReport(t, s, testReport("2026-08-05", "五日"))
	putReport(t, s, testReport("2026-08-07", "舊標題"))
	putReport(t, s, testReport("2026-08-06", "六日"))
	putReport(t, s, testReport("2026-08-07", "新標題"))

	latest, err := s.LatestReport(ctx)
	if err != nil || latest.Date != "2026-08-07" || latest.Headline != "新標題" {
		t.Fatalf("LatestReport() = %+v, err = %v", latest, err)
	}
	count, err := s.CountReports(ctx)
	if err != nil || count != 3 {
		t.Fatalf("CountReports() = %d, err = %v, want 3", count, err)
	}
	rows, err := s.ListReports(ctx, 1, 2)
	if err != nil || len(rows) != 2 || rows[0].Date != "2026-08-07" || rows[1].Date != "2026-08-06" {
		t.Fatalf("ListReports() = %+v, err = %v", rows, err)
	}
	page2, err := s.ListReports(ctx, 2, 2)
	if err != nil || len(page2) != 1 || page2[0].Date != "2026-08-05" {
		t.Fatalf("ListReports(page 2) = %+v, err = %v", page2, err)
	}
}

func TestAPIKeysRemainViewableAndRevocationPersists(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	created, err := s.CreateAPIKey(ctx, "cowork-daily", "plain-secret-token")
	if err != nil {
		t.Fatal(err)
	}
	found, err := s.APIKeyByToken(ctx, created.Token)
	if err != nil || found.Token != "plain-secret-token" || !found.Active() {
		t.Fatalf("APIKeyByToken() = %+v, err = %v", found, err)
	}
	if err := s.RevokeAPIKey(ctx, created.ID); err != nil {
		t.Fatal(err)
	}
	revoked, err := s.APIKeyByToken(ctx, created.Token)
	if err != nil || revoked.Active() || revoked.RevokedAt == nil {
		t.Fatalf("revoked APIKeyByToken() = %+v, err = %v", revoked, err)
	}
	keys, err := s.ListAPIKeys(ctx)
	if err != nil || len(keys) != 1 || keys[0].Token != created.Token || keys[0].RevokedAt == nil {
		t.Fatalf("ListAPIKeys() = %+v, err = %v", keys, err)
	}
}

func TestUpdateLogNewestFirst(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	for i := 0; i < 2; i++ {
		if err := s.AddUpdateLog(ctx, UpdateLog{KeyName: "key", ReportDate: fmt.Sprintf("2026-08-0%d", i+6), Action: "ingest_ok"}); err != nil {
			t.Fatal(err)
		}
	}
	logs, err := s.ListUpdateLogs(ctx, 10)
	if err != nil || len(logs) != 2 || logs[0].ReportDate != "2026-08-07" {
		t.Fatalf("ListUpdateLogs() = %+v, err = %v", logs, err)
	}
}

// Fractional seconds must not change the stored width, or ORDER BY at DESC puts
// older rows first: RFC3339Nano renders 100ms as ".1Z" and 120ms as ".12Z",
// and ".1Z" sorts after ".12Z".
func TestUpdateLogOrdersAcrossFractionalSeconds(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	base := time.Date(2026, 8, 7, 7, 50, 0, 0, time.UTC)
	for i, offset := range []time.Duration{100 * time.Millisecond, 120 * time.Millisecond} {
		at := base.Add(offset)
		entry := UpdateLog{At: at, KeyName: "key", ReportDate: fmt.Sprintf("2026-08-0%d", i+6), Action: "ingest_ok"}
		if err := s.AddUpdateLog(ctx, entry); err != nil {
			t.Fatal(err)
		}
	}
	logs, err := s.ListUpdateLogs(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(logs) != 2 || logs[0].ReportDate != "2026-08-07" {
		t.Fatalf("newest entry should come first, got %+v", logs)
	}
}
