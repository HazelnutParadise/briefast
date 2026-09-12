package report

import (
	"strings"
	"testing"
	"time"
)

func validConference() Conference {
	return Conference{
		Symbol: "2330", Name: "台積電", Market: MarketTWSE,
		HeldOn: "2026-10-16", HeldAt: "14:00", Venue: "線上法說會", AnnouncedOn: "2026-10-10",
		Prediction: PredictionBull, Headline: "營收連三月年增，法說可望上修展望",
		SummaryMD: "判斷依據", WatchMD: "- 觀察 2 奈米量產時程",
		Fundamentals: &Fundamentals{
			Revenue: &RevenueFigures{Month: "2026-08", Current: 100, PrevMonth: 90, LastYearMonth: 80, YTD: 700, LastYearYTD: 600},
			Quarter: &QuarterFigures{Label: "2026Q2 累計", Revenue: 500, OperatingIncome: 200, NetIncome: 150, EPS: "9.10"},
		},
		Chips:       &Chips{Date: "2026-10-09", ForeignNet: 1000, TrustNet: 0, DealerNet: -500, TotalNet: 500},
		Sources:     []Source{{Title: "來源", URL: "https://example.com/a"}},
		GeneratedAt: "2026-10-10T07:50:00+08:00",
	}
}

func TestConferenceValidateMatrix(t *testing.T) {
	tests := []struct {
		name string
		edit func(*Conference)
		want string
	}{
		{name: "valid", edit: func(*Conference) {}},
		{name: "no fundamentals", edit: func(c *Conference) { c.Fundamentals = nil }},
		{name: "no chips and empty time", edit: func(c *Conference) { c.Chips = nil; c.HeldAt = "" }},
		{name: "prediction up", edit: func(c *Conference) { c.Prediction = "up" }, want: "prediction"},
		{name: "market otc", edit: func(c *Conference) { c.Market = "otc" }, want: "market"},
		{name: "held_at 2pm", edit: func(c *Conference) { c.HeldAt = "2pm" }, want: "held_at"},
		{name: "held_on bad", edit: func(c *Conference) { c.HeldOn = "2026-13-01" }, want: "held_on"},
		{name: "announced_on bad", edit: func(c *Conference) { c.AnnouncedOn = "x" }, want: "announced_on"},
		{name: "headline blank", edit: func(c *Conference) { c.Headline = " \n" }, want: "headline"},
		{name: "summary blank", edit: func(c *Conference) { c.SummaryMD = "" }, want: "summary_md"},
		{name: "watch blank", edit: func(c *Conference) { c.WatchMD = "" }, want: "watch_md"},
		{name: "source url missing", edit: func(c *Conference) { c.Sources[0].URL = "" }, want: "sources[0].url"},
		{name: "chips date bad", edit: func(c *Conference) { c.Chips.Date = "2026-13-01" }, want: "chips.date"},
		{name: "revenue month bad", edit: func(c *Conference) { c.Fundamentals.Revenue.Month = "202608" }, want: "fundamentals.revenue.month"},
		{name: "generated_at bad", edit: func(c *Conference) { c.GeneratedAt = "2026-10-10 07:50" }, want: "generated_at"},
		{name: "symbol blank", edit: func(c *Conference) { c.Symbol = "" }, want: "symbol"},
		{name: "symbol too long", edit: func(c *Conference) { c.Symbol = "12345678901" }, want: "symbol"},
		{name: "symbol padded", edit: func(c *Conference) { c.Symbol = " 2330" }, want: "symbol"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := validConference()
			tt.edit(&c)
			errs := strings.Join(c.Validate(), "\n")
			if tt.want == "" && errs != "" {
				t.Fatalf("Validate() errors = %q, want none", errs)
			}
			if tt.want != "" && !strings.Contains(errs, tt.want) {
				t.Fatalf("Validate() errors = %q, want substring %q", errs, tt.want)
			}
		})
	}
}

func TestConferenceValidateListsEveryViolation(t *testing.T) {
	c := validConference()
	c.Prediction = "up"
	c.Headline = ""
	errs := c.Validate()
	if len(errs) != 2 {
		t.Fatalf("Validate() = %v, want two violations", errs)
	}
}

func TestOutcomeMatrix(t *testing.T) {
	tests := []struct {
		prediction string
		pre, post  float64
		want       string
	}{
		{PredictionBull, 1080, 1095, OutcomeHit},
		{PredictionBull, 1080, 1080, OutcomeMiss},
		{PredictionBull, 1080, 1060, OutcomeMiss},
		{PredictionBear, 500, 480, OutcomeHit},
		{PredictionBear, 500, 510, OutcomeMiss},
		{PredictionNone, 500, 520, ""},
	}
	for _, tt := range tests {
		if got := Outcome(tt.prediction, tt.pre, tt.post); got != tt.want {
			t.Errorf("Outcome(%s, %v, %v) = %q, want %q", tt.prediction, tt.pre, tt.post, got, tt.want)
		}
	}
}

func TestTaipeiDateCrossesMidnight(t *testing.T) {
	// 23:30 UTC on the 13th is already 07:30 on the 14th in Taipei.
	at := time.Date(2026, 9, 13, 23, 30, 0, 0, time.UTC)
	if got := TaipeiDate(at); got != "2026-09-14" {
		t.Fatalf("TaipeiDate() = %q, want 2026-09-14", got)
	}
}
