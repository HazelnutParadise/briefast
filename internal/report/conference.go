package report

import (
	"fmt"
	"strings"
	"time"
)

// Conference predictions carry a single horizon — the first trading day after
// the event — so the four-way short/long call vocabulary does not apply here.
const (
	PredictionBull = "bull"
	PredictionBear = "bear"
	PredictionNone = "none"
)

const (
	MarketTWSE = "twse"
	MarketTPEx = "tpex"
)

// maxSymbolLength matches the page router's cut-off so every stored brief is
// reachable at its URL; Taiwan tickers are at most six characters.
const maxSymbolLength = 10

const (
	OutcomeHit  = "hit"
	OutcomeMiss = "miss"
)

// Conference is one company-held investor conference with its pre-event brief.
type Conference struct {
	Symbol       string        `json:"symbol"`
	Name         string        `json:"name"`
	Market       string        `json:"market"`
	HeldOn       string        `json:"held_on"`
	HeldAt       string        `json:"held_at"`
	Venue        string        `json:"venue"`
	AnnouncedOn  string        `json:"announced_on"`
	Prediction   string        `json:"prediction"`
	Headline     string        `json:"headline"`
	SummaryMD    string        `json:"summary_md"`
	WatchMD      string        `json:"watch_md"`
	Fundamentals *Fundamentals `json:"fundamentals,omitempty"`
	Chips        *Chips        `json:"chips,omitempty"`
	Sources      []Source      `json:"sources"`
	GeneratedAt  string        `json:"generated_at"`
}

// Fundamentals holds the structured figures fetched for the brief. Both parts
// are optional because a symbol can be missing from either dataset.
type Fundamentals struct {
	Revenue *RevenueFigures `json:"revenue,omitempty"`
	Quarter *QuarterFigures `json:"quarter,omitempty"`
}

// RevenueFigures are monthly revenue values in thousands of NTD, as published
// by the exchanges' monthly revenue datasets.
type RevenueFigures struct {
	Month         string `json:"month"`
	Current       int64  `json:"current"`
	PrevMonth     int64  `json:"prev_month"`
	LastYearMonth int64  `json:"last_year_month"`
	YTD           int64  `json:"ytd"`
	LastYearYTD   int64  `json:"last_year_ytd"`
}

// QuarterFigures are the latest cumulative quarter figures in thousands of
// NTD; EPS stays a string so "0.38" is not reformatted on the way through.
type QuarterFigures struct {
	Label           string `json:"label"`
	Revenue         int64  `json:"revenue"`
	OperatingIncome int64  `json:"operating_income"`
	NetIncome       int64  `json:"net_income"`
	EPS             string `json:"eps"`
}

func (c Conference) Validate() []string {
	errs := make([]string, 0)
	if trimmed := strings.TrimSpace(c.Symbol); trimmed == "" || len(trimmed) > maxSymbolLength || trimmed != c.Symbol {
		errs = append(errs, "symbol 必須是 1 到 10 個字元、不含前後空白的代號")
	}
	if strings.TrimSpace(c.Name) == "" {
		errs = append(errs, "name 不得為空")
	}
	if c.Market != MarketTWSE && c.Market != MarketTPEx {
		errs = append(errs, "market 必須是 twse 或 tpex")
	}
	if !validDate(c.HeldOn) {
		errs = append(errs, "held_on 必須是有效的 YYYY-MM-DD 日期")
	}
	if !validClock(c.HeldAt) {
		errs = append(errs, "held_at 必須是 HH:MM 或空字串")
	}
	if !validDate(c.AnnouncedOn) {
		errs = append(errs, "announced_on 必須是有效的 YYYY-MM-DD 日期")
	}
	if !validPrediction(c.Prediction) {
		errs = append(errs, "prediction 必須是 bull、bear 或 none")
	}
	if strings.TrimSpace(c.Headline) == "" {
		errs = append(errs, "headline 不得為空")
	}
	if strings.TrimSpace(c.SummaryMD) == "" {
		errs = append(errs, "summary_md 不得為空")
	}
	if strings.TrimSpace(c.WatchMD) == "" {
		errs = append(errs, "watch_md 不得為空")
	}
	if c.Fundamentals != nil && c.Fundamentals.Revenue != nil && !validYearMonth(c.Fundamentals.Revenue.Month) {
		errs = append(errs, "fundamentals.revenue.month 必須是 YYYY-MM")
	}
	if c.Chips != nil && !validDate(c.Chips.Date) {
		errs = append(errs, "chips.date 必須是有效的 YYYY-MM-DD 日期")
	}
	for i, source := range c.Sources {
		if strings.TrimSpace(source.URL) == "" {
			errs = append(errs, fmt.Sprintf("sources[%d].url 不得為空", i))
		}
	}
	if _, err := time.Parse(time.RFC3339, c.GeneratedAt); err != nil {
		errs = append(errs, "generated_at 必須是含時區的 RFC 3339 時間")
	}
	return errs
}

func validPrediction(value string) bool {
	switch value {
	case PredictionBull, PredictionBear, PredictionNone:
		return true
	default:
		return false
	}
}

func validClock(value string) bool {
	if value == "" {
		return true
	}
	if len(value) != len("15:04") {
		return false
	}
	parsed, err := time.Parse("15:04", value)
	return err == nil && parsed.Format("15:04") == value
}

func validYearMonth(value string) bool {
	if len(value) != len("2006-01") {
		return false
	}
	parsed, err := time.Parse("2006-01", value)
	return err == nil && parsed.Format("2006-01") == value
}

// Outcome applies the settlement rule: a flat close never counts as a hit,
// and a brief without a direction produces no outcome at all.
func Outcome(prediction string, preClose, postClose float64) string {
	switch prediction {
	case PredictionBull:
		if postClose > preClose {
			return OutcomeHit
		}
		return OutcomeMiss
	case PredictionBear:
		if postClose < preClose {
			return OutcomeHit
		}
		return OutcomeMiss
	default:
		return ""
	}
}

var taipei = loadTaipei()

func loadTaipei() *time.Location {
	if loc, err := time.LoadLocation("Asia/Taipei"); err == nil {
		return loc
	}
	return time.FixedZone("Asia/Taipei", 8*60*60)
}

// TaipeiDate returns the calendar date in Taipei for the given instant.
func TaipeiDate(t time.Time) string {
	return t.In(taipei).Format("2006-01-02")
}
