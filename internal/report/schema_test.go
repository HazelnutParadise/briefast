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
		Industries: []Industry{{
			Name: "科技",
			Events: []IndustryEvent{
				{Headline: "記憶體漲價改善上游產品組合", SummaryMD: "產業摘要"},
				{Headline: "AI 追單推升封測產能利用率", SummaryMD: "另一則產業摘要"},
			},
			WatchMD: "- 產業觀察",
		}},
		StockNews: []StockNews{{
			Symbol: "2330", Name: "台積電", Call: CallShortBull,
			Headline: "先進製程放量支撐短線動能", SummaryMD: "新聞摘要", WatchMD: "- 個股觀察", Sources: []Source{{Title: "新聞", URL: "https://example.com/news"}},
		}},
		GeneratedAt: "2026-08-07T07:50:00+08:00",
	}
}

func TestReportValidateIndustryEventsAndStockHeadline(t *testing.T) {
	tests := []struct {
		name string
		edit func(*Report)
		want string
	}{
		{name: "valid", edit: func(*Report) {}},
		{name: "empty events", edit: func(r *Report) { r.Industries[0].Events = nil }, want: "industries[0].events"},
		{name: "event headline missing", edit: func(r *Report) { r.Industries[0].Events[1].Headline = " \t\n" }, want: "industries[0].events[1].headline"},
		{name: "event summary missing", edit: func(r *Report) { r.Industries[0].Events[1].SummaryMD = " \r\n" }, want: "industries[0].events[1].summary_md"},
		{name: "stock headline missing", edit: func(r *Report) { r.StockNews[0].Headline = " " }, want: "stock_news[0].headline"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validReport()
			tt.edit(&r)
			errs := strings.Join(r.Validate(), "\n")
			if tt.want == "" && errs != "" {
				t.Fatalf("Validate() errors = %q, want none", errs)
			}
			if tt.want != "" && !strings.Contains(errs, tt.want) {
				t.Fatalf("Validate() errors = %q, want substring %q", errs, tt.want)
			}
		})
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
	r.Industries[0].Events[1].Headline = ""
	r.Industries[0].Events[1].SummaryMD = ""
	r.StockNews[0].Headline = ""
	errs := r.Validate()
	if len(errs) != 5 {
		t.Fatalf("Validate() returned %d errors, want 5: %v", len(errs), errs)
	}
}

func chipsSample() *Chips {
	margin, short := int64(-18270), int64(9056)
	return &Chips{
		Date:       "2026-08-20",
		ForeignNet: 54758664, TrustNet: -15000, DealerNet: 2063215, TotalNet: 56806879,
		MarginChange: &margin, ShortChange: &short,
	}
}

func TestReportValidateChips(t *testing.T) {
	tests := []struct {
		name string
		edit func(*Report)
		want string
	}{
		{name: "chips absent", edit: func(*Report) {}},
		{name: "chips present and valid", edit: func(r *Report) { r.StockNews[0].Chips = chipsSample() }},
		{name: "chips without margin fields", edit: func(r *Report) {
			c := chipsSample()
			c.MarginChange, c.ShortChange = nil, nil
			r.StockNews[0].Chips = c
		}},
		{name: "chips total not recomputed", edit: func(r *Report) {
			c := chipsSample()
			c.TotalNet = 0
			r.StockNews[0].Chips = c
		}},
		{name: "chips date malformed", edit: func(r *Report) {
			c := chipsSample()
			c.Date = "2026-8-20"
			r.StockNews[0].Chips = c
		}, want: "stock_news[0].chips.date"},
		{name: "chips date empty", edit: func(r *Report) {
			c := chipsSample()
			c.Date = ""
			r.StockNews[0].Chips = c
		}, want: "stock_news[0].chips.date"},
		{name: "chips date malformed on third entry", edit: func(r *Report) {
			first := r.StockNews[0]
			r.StockNews = append(r.StockNews, first, first)
			c := chipsSample()
			c.Date = "20260820"
			r.StockNews[2].Chips = c
		}, want: "stock_news[2].chips.date"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := validReport()
			tt.edit(&r)
			errs := strings.Join(r.Validate(), "\n")
			if tt.want == "" && errs != "" {
				t.Fatalf("Validate() errors = %q, want none", errs)
			}
			if tt.want != "" && !strings.Contains(errs, tt.want) {
				t.Fatalf("Validate() errors = %q, want substring %q", errs, tt.want)
			}
		})
	}
}
