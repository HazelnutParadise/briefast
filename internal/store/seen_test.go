package store

import (
	"context"
	"errors"
	"path/filepath"
	"testing"
	"time"
)

func TestSeenArticleRecordMatchAndDuplicate(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	a := SeenArticle{SourceID: "cnyes.com:6001", URLHash: "u1", TitleHash: "t1", Decision: "reported", Sketch: []string{"0001"}}
	added, err := s.RecordSeenArticle(ctx, a)
	if err != nil || !added {
		t.Fatalf("first record added=%v err=%v", added, err)
	}
	added, err = s.RecordSeenArticle(ctx, SeenArticle{URLHash: "u1", Decision: "skipped"})
	if err != nil || added {
		t.Fatalf("duplicate URL added=%v err=%v", added, err)
	}
	for _, c := range []struct{ sid, url, title, want string }{
		{"cnyes.com:6001", "", "", "source_id"},
		{"", "u1", "", "url"},
		{"", "", "t1", "title"},
		{"", "", "", ""},
		{"other", "u2", "t2", ""},
	} {
		got, err := s.SeenMatch(ctx, c.sid, c.url, c.title)
		if err != nil || got != c.want {
			t.Errorf("SeenMatch(%q,%q,%q) = %q, %v; want %q", c.sid, c.url, c.title, got, err, c.want)
		}
	}
	sketches, err := s.SeenSketches(ctx)
	if err != nil || len(sketches) != 1 || sketches[0].Sketch[0] != "0001" {
		t.Fatalf("sketches = %+v, %v", sketches, err)
	}
}

func TestSeenEventsAddOnceAndExpire(t *testing.T) {
	s := openTestStore(t)
	ctx := context.Background()
	start := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	s.now = func() time.Time { return start }
	if added, err := s.AddSeenEvent(ctx, "2330|月營收|2026-07", "首次"); err != nil || !added {
		t.Fatalf("add = %v, %v", added, err)
	}
	if added, err := s.AddSeenEvent(ctx, "2330|月營收|2026-07", "覆寫"); err != nil || added {
		t.Fatalf("second add = %v, %v", added, err)
	}
	event, err := s.SeenEventByKey(ctx, "2330|月營收|2026-07")
	if err != nil || event.Note != "首次" || !event.FirstSeenAt.Equal(start) {
		t.Fatalf("event = %+v, %v", event, err)
	}
	if _, err := s.RecordSeenArticle(ctx, SeenArticle{URLHash: "old", Decision: "skipped"}); err != nil {
		t.Fatal(err)
	}

	s.now = func() time.Time { return start.Add(31 * 24 * time.Hour) }
	if err := s.PruneSeen(ctx); err != nil {
		t.Fatal(err)
	}
	if _, err := s.SeenEventByKey(ctx, "2330|月營收|2026-07"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("31 天後事件應過期, err = %v", err)
	}
	if got, _ := s.SeenMatch(ctx, "", "old", ""); got != "url" {
		t.Fatal("31 天後文章不應過期")
	}

	s.now = func() time.Time { return start.Add(91 * 24 * time.Hour) }
	if err := s.PruneSeen(ctx); err != nil {
		t.Fatal(err)
	}
	if got, _ := s.SeenMatch(ctx, "", "old", ""); got != "" {
		t.Fatal("91 天後文章應過期")
	}
}

func TestVersionTwoDatabaseUpgradesWithSeenTables(t *testing.T) {
	path := filepath.Join(t.TempDir(), "briefast.db")
	legacy, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	putReport(t, legacy, testReport("2026-08-07", "舊報告"))
	if _, err := legacy.DB().Exec("DROP TABLE seen_articles; DROP TABLE seen_events; DELETE FROM schema_migrations WHERE version = 3"); err != nil {
		t.Fatal(err)
	}
	if err := legacy.Close(); err != nil {
		t.Fatal(err)
	}
	s, err := Open(path)
	if err != nil {
		t.Fatalf("reopen at version 2: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	if r, err := s.ReportByDate(context.Background(), "2026-08-07"); err != nil || r.Headline != "舊報告" {
		t.Fatalf("report after upgrade = %+v, err = %v", r, err)
	}
	if added, err := s.AddSeenEvent(context.Background(), "2330|法說會|2026-10", ""); err != nil || !added {
		t.Fatalf("seen tables unusable after upgrade: %v %v", added, err)
	}
}
