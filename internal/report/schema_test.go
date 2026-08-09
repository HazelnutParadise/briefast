package report

import (
	"strings"
	"testing"
)

func validReport() Report {
	return Report{
		Date:       "2026-08-07",
		Headline:   "美股收紅",
		OverviewMD: "盤前總覽",
		WatchMD:    "- 今日觀察",
		Calls: Calls{ShortBull: []CallEntry{{
			Symbol: "2330", Name: "台積電", Reason: "一句理由",
		}}},
		Industries: []Industry{{Name: "科技", SummaryMD: "產業摘要", WatchMD: "- 產業觀察"}},
		StockNews: []StockNews{{
			Symbol: "2330", Name: "台積電", Call: CallShortBull,
			SummaryMD: "新聞摘要", WatchMD: "- 個股觀察", Sources: []Source{{Title: "新聞", URL: "https://example.com/news"}},
		}},
		GeneratedAt: "2026-08-07T07:50:00+08:00",
	}
}

func TestReportValidateEntryWatchMD(t *testing.T) {
	tests := []struct {
		name     string
		industry string
		stock    string
		want     []string
	}{
		{name: "valid", industry: "- 產業觀察", stock: "- 個股觀察"},
		{name: "missing fields", want: []string{"industries[0].watch_md", "stock_news[0].watch_md"}},
		{name: "whitespace fields", industry: " \t\n", stock: " \r\n", want: []string{"industries[0].watch_md", "stock_news[0].watch_md"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validReport()
			r.Industries[0].WatchMD = tt.industry
			r.StockNews[0].WatchMD = tt.stock
			errs := strings.Join(r.Validate(), "\n")
			if len(tt.want) == 0 && errs != "" {
				t.Fatalf("Validate() errors = %q, want none", errs)
			}
			for _, want := range tt.want {
				if !strings.Contains(errs, want) {
					t.Errorf("Validate() errors = %q, want substring %q", errs, want)
				}
			}
		})
	}
}

func TestReportValidate(t *testing.T) {
	tests := []struct {
		name string
		edit func(*Report)
		want string
	}{
		{name: "valid", edit: func(*Report) {}},
		{name: "malformed date", edit: func(r *Report) { r.Date = "2026/08/07" }, want: "date"},
		{name: "impossible date", edit: func(r *Report) { r.Date = "2026-02-30" }, want: "date"},
		{name: "empty headline", edit: func(r *Report) { r.Headline = " " }, want: "headline"},
		{name: "empty overview", edit: func(r *Report) { r.OverviewMD = "" }, want: "overview_md"},
		{name: "empty watch", edit: func(r *Report) { r.WatchMD = "" }, want: "watch_md"},
		{name: "invalid call", edit: func(r *Report) { r.StockNews[0].Call = "maybe" }, want: "stock_news[0].call"},
		{name: "source URL missing", edit: func(r *Report) { r.StockNews[0].Sources[0].URL = "" }, want: "sources[0].url"},
		{name: "call detail missing", edit: func(r *Report) { r.StockNews = nil }, want: "symbol \"2330\""},
		{name: "neutral news without call membership", edit: func(r *Report) {
			r.Calls = Calls{}
			r.StockNews[0].Call = CallNone
		}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validReport()
			tt.edit(&r)
			errs := r.Validate()
			if tt.want == "" && len(errs) != 0 {
				t.Fatalf("Validate() errors = %v, want none", errs)
			}
			if tt.want != "" && !strings.Contains(strings.Join(errs, "\n"), tt.want) {
				t.Fatalf("Validate() errors = %v, want substring %q", errs, tt.want)
			}
		})
	}
}

func TestReportValidateReturnsEveryViolation(t *testing.T) {
	r := validReport()
	r.Date = "bad"
	r.Headline = ""
	errs := r.Validate()
	if len(errs) != 2 {
		t.Fatalf("Validate() returned %d errors, want 2: %v", len(errs), errs)
	}
}
