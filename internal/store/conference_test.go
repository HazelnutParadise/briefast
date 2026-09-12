package store

import (
	"context"
	"encoding/json"
	"errors"
	"path/filepath"
	"testing"

	"github.com/HazelnutParadise/briefast/internal/report"
)

func testConference(symbol, heldOn, prediction string) report.Conference {
	return report.Conference{
		Symbol: symbol, Name: "公司" + symbol, Market: report.MarketTWSE,
		HeldOn: heldOn, HeldAt: "14:00", Venue: "線上", AnnouncedOn: "2026-09-10",
		Prediction: prediction, Headline: "標題", SummaryMD: "依據", WatchMD: "- 觀察",
		Sources: []report.Source{{Title: "來源", URL: "https://example.com"}}, GeneratedAt: "2026-09-10T07:50:00+08:00",
	}
}

func putConference(t *testing.T, s *Store, c report.Conference) {
	t.Helper()
	payload, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertConference(context.Background(), c, payload); err != nil {
		t.Fatalf("UpsertConference() error = %v", err)
	}
}

func TestFreshDatabaseReachesVersionTwoWithConferences(t *testing.T) {
	s := openTestStore(t)
	var version int
	if err := s.DB().QueryRow("SELECT MAX(version) FROM schema_migrations").Scan(&version); err != nil || version != 2 {
		t.Fatalf("migration version = %d, err = %v, want 2", version, err)
	}
	var name string
	if err := s.DB().QueryRow("SELECT name FROM sqlite_master WHERE type = 'table' AND name = 'conferences'").Scan(&name); err != nil {
		t.Fatalf("conferences table missing: %v", err)
	}
}

func TestVersionOneDatabaseUpgradesToVersionTwoKeepingReports(t *testing.T) {
	path := filepath.Join(t.TempDir(), "briefast.db")
	// Build a database exactly as the version-1 code would have left it.
	legacy, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	putReport(t, legacy, testReport("2026-08-07", "舊報告"))
	if _, err := legacy.DB().Exec("DROP TABLE conferences; DELETE FROM schema_migrations WHERE version = 2"); err != nil {
		t.Fatal(err)
	}
	if err := legacy.Close(); err != nil {
		t.Fatal(err)
	}

	s, err := Open(path)
	if err != nil {
		t.Fatalf("reopen at version 1: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	var version int
	if err := s.DB().QueryRow("SELECT MAX(version) FROM schema_migrations").Scan(&version); err != nil || version != 2 {
		t.Fatalf("migration version = %d, err = %v, want 2", version, err)
	}
	r, err := s.ReportByDate(context.Background(), "2026-08-07")
	if err != nil || r.Headline != "舊報告" {
		t.Fatalf("report after upgrade = %+v, err = %v", r, err)
	}
	putConference(t, s, testConference("2330", "2026-10-16", report.PredictionBull))
}

func TestConferenceUpsertKeepsSettlementAndListsSplitByToday(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	key, err := s.CreateAPIKey(ctx, "cowork", "tok")
	if err != nil {
		t.Fatal(err)
	}
	putConference(t, s, testConference("1314", "2026-09-11", report.PredictionBull))
	putConference(t, s, testConference("2330", "2026-09-12", report.PredictionBear))
	putConference(t, s, testConference("3443", "2026-09-14", report.PredictionNone))
	putConference(t, s, testConference("1101", "2026-09-16", report.PredictionBull))
	putConference(t, s, testConference("1101", "2026-09-14", report.PredictionBull))

	outcome, err := s.SettleConference(ctx, "2330", "2026-09-12", Settlement{PreCloseDate: "2026-09-11", PreClose: 500, PostCloseDate: "2026-09-14", PostClose: 480}, key)
	if err != nil || outcome != report.OutcomeHit {
		t.Fatalf("SettleConference() = %q, err = %v", outcome, err)
	}
	if _, err := s.SettleConference(ctx, "9999", "2026-09-12", Settlement{}, key); !errors.Is(err, ErrNotFound) {
		t.Fatalf("settle unknown error = %v, want ErrNotFound", err)
	}

	updated := testConference("2330", "2026-09-12", report.PredictionBear)
	updated.SummaryMD = "新依據"
	putConference(t, s, updated)
	got, err := s.ConferenceByKey(ctx, "2330", "2026-09-12")
	if err != nil || got.Brief.SummaryMD != "新依據" || got.Settlement == nil || got.Settlement.Outcome != report.OutcomeHit || got.Settlement.PostClose != 480 {
		t.Fatalf("ConferenceByKey() after re-upsert = %+v, err = %v", got, err)
	}

	upcoming, err := s.ListUpcomingConferences(ctx, "2026-09-14")
	if err != nil || len(upcoming) != 3 || upcoming[0].Symbol != "1101" || upcoming[0].HeldOn != "2026-09-14" || upcoming[1].Symbol != "3443" || upcoming[2].HeldOn != "2026-09-16" {
		t.Fatalf("ListUpcomingConferences() = %+v, err = %v", upcoming, err)
	}
	pending, err := s.ListPendingSettlement(ctx, "2026-09-14")
	if err != nil || len(pending) != 1 || pending[0].Symbol != "1314" {
		t.Fatalf("ListPendingSettlement() = %+v, err = %v", pending, err)
	}
	settled, err := s.ListSettledConferences(ctx, 5)
	if err != nil || len(settled) != 1 || settled[0].Symbol != "2330" || settled[0].Settlement.PreClose != 500 {
		t.Fatalf("ListSettledConferences() = %+v, err = %v", settled, err)
	}
	all, err := s.ListConferences(ctx)
	if err != nil || len(all) != 5 || all[0].HeldOn != "2026-09-16" {
		t.Fatalf("ListConferences() = %+v, err = %v", all, err)
	}
	if _, err := s.ConferenceByKey(ctx, "2330", "2026-01-01"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("ConferenceByKey unknown error = %v", err)
	}

	logs, err := s.ListUpdateLogs(ctx, 10)
	if err != nil || len(logs) != 1 || logs[0].Action != "settle_ok" || logs[0].ReportDate != "2026-09-12" {
		t.Fatalf("ListUpdateLogs() = %+v, err = %v", logs, err)
	}
	if err := s.IngestConference(ctx, testConference("2454", "2026-09-20", report.PredictionBull), []byte(`{}`), key); err != nil {
		t.Fatal(err)
	}
	logs, _ = s.ListUpdateLogs(ctx, 10)
	if len(logs) != 2 || logs[0].Action != "conference_ok" || logs[0].ReportDate != "2026-09-20" {
		t.Fatalf("conference_ok log missing: %+v", logs)
	}
}
