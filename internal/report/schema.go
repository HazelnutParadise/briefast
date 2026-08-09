package report

import (
	"fmt"
	"strings"
	"time"
)

const (
	CallShortBull = "short_bull"
	CallShortBear = "short_bear"
	CallLongBull  = "long_bull"
	CallLongBear  = "long_bear"
	CallNone      = "none"
)

type Report struct {
	Date        string      `json:"date"`
	Headline    string      `json:"headline"`
	OverviewMD  string      `json:"overview_md"`
	WatchMD     string      `json:"watch_md"`
	Calls       Calls       `json:"calls"`
	Industries  []Industry  `json:"industries"`
	StockNews   []StockNews `json:"stock_news"`
	GeneratedAt string      `json:"generated_at"`
}

type Calls struct {
	ShortBull []CallEntry `json:"short_bull"`
	ShortBear []CallEntry `json:"short_bear"`
	LongBull  []CallEntry `json:"long_bull"`
	LongBear  []CallEntry `json:"long_bear"`
}

type CallEntry struct {
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type Industry struct {
	Name      string `json:"name"`
	SummaryMD string `json:"summary_md"`
	WatchMD   string `json:"watch_md"`
}

type StockNews struct {
	Symbol    string   `json:"symbol"`
	Name      string   `json:"name"`
	Call      string   `json:"call"`
	SummaryMD string   `json:"summary_md"`
	WatchMD   string   `json:"watch_md"`
	Sources   []Source `json:"sources"`
}

type Source struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (r Report) Validate() []string {
	errs := make([]string, 0)
	if !validDate(r.Date) {
		errs = append(errs, "date 必須是有效的 YYYY-MM-DD 日期")
	}
	if strings.TrimSpace(r.Headline) == "" {
		errs = append(errs, "headline 不得為空")
	}
	if strings.TrimSpace(r.OverviewMD) == "" {
		errs = append(errs, "overview_md 不得為空")
	}
	if strings.TrimSpace(r.WatchMD) == "" {
		errs = append(errs, "watch_md 不得為空")
	}
	for i, industry := range r.Industries {
		if strings.TrimSpace(industry.WatchMD) == "" {
			errs = append(errs, fmt.Sprintf("industries[%d].watch_md 不得為空", i))
		}
	}

	newsSymbols := make(map[string]struct{}, len(r.StockNews))
	for i, item := range r.StockNews {
		newsSymbols[item.Symbol] = struct{}{}
		if strings.TrimSpace(item.WatchMD) == "" {
			errs = append(errs, fmt.Sprintf("stock_news[%d].watch_md 不得為空", i))
		}
		if !validCall(item.Call) {
			errs = append(errs, fmt.Sprintf("stock_news[%d].call 必須是 short_bull、short_bear、long_bull、long_bear 或 none", i))
		}
		for j, source := range item.Sources {
			if strings.TrimSpace(source.URL) == "" {
				errs = append(errs, fmt.Sprintf("stock_news[%d].sources[%d].url 不得為空", i, j))
			}
		}
	}

	for _, group := range []struct {
		name    string
		entries []CallEntry
	}{
		{CallShortBull, r.Calls.ShortBull},
		{CallShortBear, r.Calls.ShortBear},
		{CallLongBull, r.Calls.LongBull},
		{CallLongBear, r.Calls.LongBear},
	} {
		for _, entry := range group.entries {
			if _, ok := newsSymbols[entry.Symbol]; !ok {
				errs = append(errs, fmt.Sprintf("calls.%s 的 symbol %q 在 stock_news 中沒有對應條目", group.name, entry.Symbol))
			}
		}
	}
	return errs
}

func validCall(call string) bool {
	switch call {
	case CallShortBull, CallShortBear, CallLongBull, CallLongBear, CallNone:
		return true
	default:
		return false
	}
}

func validDate(value string) bool {
	if len(value) != len("2006-01-02") {
		return false
	}
	parsed, err := time.Parse("2006-01-02", value)
	return err == nil && parsed.Format("2006-01-02") == value
}
