package site

import (
	"context"
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/HazelnutParadise/briefast/internal/report"
	"github.com/HazelnutParadise/briefast/internal/store"
	sy "github.com/HazelnutParadise/syralit"
)

func setupSite(t *testing.T) (*store.Store, *Site) {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "briefast.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s, New(s)
}

func fullReport(date string) report.Report {
	return report.Report{
		Date: date, Headline: "盤前頭條", OverviewMD: "總覽內容", WatchMD: "- 觀察一\n- 觀察二",
		Calls: report.Calls{
			ShortBull: []report.CallEntry{{Symbol: "1001", Name: "短多股", Reason: "短多理由"}},
			ShortBear: []report.CallEntry{{Symbol: "1002", Name: "短空股", Reason: "短空理由"}},
			LongBull:  []report.CallEntry{{Symbol: "1003", Name: "長多股", Reason: "長多理由"}},
			LongBear:  []report.CallEntry{{Symbol: "1004", Name: "長空股", Reason: "長空理由"}},
		},
		Industries: []report.Industry{{
			Name: "科技",
			Events: []report.IndustryEvent{
				{Headline: "記憶體漲價改善上游產品組合", SummaryMD: "記憶體產業內容"},
				{Headline: "AI 追單推升封測產能利用率", SummaryMD: "封測產業內容"},
			},
			WatchMD: "- 產業觀察內容",
		}},
		StockNews: []report.StockNews{
			{Symbol: "1001", Name: "短多股", Call: report.CallShortBull, Headline: "擴產提前支撐短線交付", SummaryMD: "短多摘要", WatchMD: "- 短多觀察內容", Sources: []report.Source{{Title: "短多來源", URL: "https://example.com/1"}}},
			{Symbol: "1002", Name: "短空股", Call: report.CallShortBear, Headline: "成本上升壓縮本季毛利", SummaryMD: "短空摘要", WatchMD: "- 短空觀察內容", Sources: []report.Source{{Title: "短空來源", URL: "https://example.com/2"}}},
			{Symbol: "1003", Name: "長多股", Call: report.CallLongBull, Headline: "新訂單擴大長期成長能見度", SummaryMD: "長多摘要", WatchMD: "- 長多觀察內容", Sources: []report.Source{{Title: "長多來源", URL: "https://example.com/3"}}},
			{Symbol: "1004", Name: "長空股", Call: report.CallLongBear, Headline: "需求轉弱拖累長期利用率", SummaryMD: "長空摘要", WatchMD: "- 長空觀察內容", Sources: []report.Source{{Title: "長空來源", URL: "https://example.com/4"}}},
			{Symbol: "1005", Name: "中性股", Call: report.CallNone, Headline: "合作備忘錄尚待正式訂單確認", SummaryMD: "方向不明摘要", WatchMD: "- 中性觀察內容", Sources: []report.Source{{Title: "中性來源", URL: "https://example.com/5"}}},
		},
		GeneratedAt: date + "T07:50:00+08:00",
	}
}

func TestStockHeadlineRendersOnlyInIdentificationArea(t *testing.T) {
	_, app := setupSite(t)
	value := fullReport("2026-08-07")
	entry := newsEntry(app.renderStockNews(value.StockNews), "n-1001")
	bodyAt := strings.Index(entry, `<div class="news-body md">`)
	if bodyAt < 0 {
		t.Fatalf("news body missing: %s", entry)
	}
	identity, body := entry[:bodyAt], entry[bodyAt:]
	for _, want := range []string{"短多股", "1001", "擴產提前支撐短線交付", "短期看漲"} {
		if !strings.Contains(identity, want) {
			t.Errorf("identification area missing %q: %s", want, identity)
		}
	}
	if strings.Contains(body, "擴產提前支撐短線交付") {
		t.Errorf("stock headline leaked into body: %s", body)
	}
	if stock, headline, tag := strings.Index(identity, "短多股"), strings.Index(identity, "擴產提前支撐短線交付"), strings.Index(identity, "短期看漲"); stock < 0 || headline <= stock || tag <= headline {
		t.Errorf("stock name, headline, and tag are out of order: %s", identity)
	}
}

func TestHeadlineStylesCreateHierarchyWithoutDecorativeBoxes(t *testing.T) {
	for _, want := range []string{
		`.event-headline{font-family:var(--serif);font-size:17px;font-weight:700`,
		`.stock-headline{font-family:var(--serif);font-size:16px;font-weight:600`,
		`.industry-event{margin:0 0 22px}`,
	} {
		if !strings.Contains(styles, want) {
			t.Errorf("headline hierarchy style missing %q", want)
		}
	}
	for _, forbidden := range []string{".industry-event{border:", ".industry-event{box-shadow:", ".stock-headline{border:", ".stock-headline{box-shadow:"} {
		if strings.Contains(styles, forbidden) {
			t.Errorf("headline style adds forbidden decoration %q", forbidden)
		}
	}
}

func saveReport(t *testing.T, s *store.Store, value report.Report) {
	t.Helper()
	payload, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertReport(context.Background(), value, payload); err != nil {
		t.Fatal(err)
	}
}

func renderedHTML(t *testing.T, at *sy.AppTest) string {
	t.Helper()
	nodes := at.FindAll("html")
	if len(nodes) != 1 {
		t.Fatalf("html nodes = %d, want 1", len(nodes))
	}
	value, _ := nodes[0].Props["html"].(string)
	return value
}

func TestHomeEmptyAndLatestReportFixedSectionOrder(t *testing.T) {
	s, app := setupSite(t)
	at := sy.NewAppTest(app.Home)
	at.Run()
	if got := renderedHTML(t, at); !strings.Contains(got, "尚無報告") || !strings.Contains(got, Disclaimer) {
		t.Fatalf("empty page missing state or disclaimer")
	}

	saveReport(t, s, fullReport("2026-08-06"))
	latest := fullReport("2026-08-07")
	latest.Headline = "最新頭條"
	saveReport(t, s, latest)
	at.Run()
	got := renderedHTML(t, at)
	if !strings.Contains(got, "最新頭條") || !strings.Contains(got, "2026 年 8 月 7 日") || !strings.Contains(got, "07:50 更新") {
		t.Fatalf("latest report masthead missing: %s", got)
	}
	sections := []string{`data-section="overview"`, `data-section="watch"`, `data-section="calls"`, `data-section="industries"`, `data-section="stock-news"`}
	last := -1
	for _, section := range sections {
		index := strings.Index(got, section)
		if index <= last {
			t.Fatalf("section %s index = %d after %d", section, index, last)
		}
		last = index
	}
}

func TestSharedVersionRefreshesOpenHome(t *testing.T) {
	s, app := setupSite(t)
	at := sy.NewAppTest(app.Home)
	at.Run()
	if !strings.Contains(renderedHTML(t, at), "尚無報告") {
		t.Fatal("initial empty state missing")
	}
	saveReport(t, s, fullReport("2026-08-07"))
	app.Notify()
	at.Run()
	if !strings.Contains(renderedHTML(t, at), "盤前頭條") {
		t.Fatal("report did not refresh after shared version notification")
	}
}

func TestCallsColorsContentAndNeutralNewsWithoutTag(t *testing.T) {
	s, app := setupSite(t)
	saveReport(t, s, fullReport("2026-08-07"))
	at := sy.NewAppTest(app.Home)
	at.Run()
	got := renderedHTML(t, at)
	for _, want := range []string{"短期看漲", "短期看跌", "長期看好・可留意", "長期看壞・不建議", "短多股", "1001", "短多理由", "科技", "記憶體產業內容", "封測產業內容", "短多來源"} {
		if !strings.Contains(got, want) {
			t.Errorf("rendered report missing %q", want)
		}
	}
	if !strings.Contains(got, `.call-col.up .stock{color:var(--up)}`) || !strings.Contains(got, `.call-col.down .stock{color:var(--down)}`) {
		t.Error("red-up/green-down styles missing")
	}
	assertBullishIsRedAndBearishIsGreen(t, got)
	assertCallsLandInMatchingColumns(t, got)
	assertStockNewsTags(t, got)
	neutralStart := strings.Index(got, `id="n-1005"`)
	neutralEnd := strings.Index(got[neutralStart:], `</article>`)
	neutral := got[neutralStart : neutralStart+neutralEnd]
	if strings.Contains(neutral, `class="tag`) || !strings.Contains(neutral, "方向不明摘要") {
		t.Fatalf("neutral entry has a tag or lacks summary: %s", neutral)
	}
}

func TestIndustryEventsRenderAsSeparateUnits(t *testing.T) {
	_, app := setupSite(t)
	value := fullReport("2026-08-07")
	got := app.renderIndustries(value.Industries)
	if count := strings.Count(got, `class="industry-event"`); count != 2 {
		t.Fatalf("industry event units = %d, want 2: %s", count, got)
	}

	for _, want := range []struct {
		headline string
		body     string
	}{
		{headline: "記憶體漲價改善上游產品組合", body: "記憶體產業內容"},
		{headline: "AI 追單推升封測產能利用率", body: "封測產業內容"},
	} {
		unit := industryEvent(got, want.headline)
		if unit == "" || !strings.Contains(unit, want.body) || strings.Index(unit, want.headline) > strings.Index(unit, want.body) {
			t.Errorf("event %q does not render above its own body %q: %s", want.headline, want.body, unit)
		}
	}
}

func industryEvent(got, headline string) string {
	headlineAt := strings.Index(got, headline)
	if headlineAt < 0 {
		return ""
	}
	start := strings.LastIndex(got[:headlineAt], `<div class="industry-event">`)
	if start < 0 {
		return ""
	}
	end := strings.Index(got[headlineAt:], `</div></div>`)
	if end < 0 {
		return ""
	}
	return got[start : headlineAt+end+len(`</div></div>`)]
}

func TestEntryWatchBlocksRenderAndLegacyReportsOmitThem(t *testing.T) {
	_, app := setupSite(t)
	current := fullReport("2026-08-07")
	got := app.renderReport(&current, navToHistory)
	if count := strings.Count(got, `class="entry-watch"`); count != 6 {
		t.Fatalf("entry watch blocks = %d, want 6", count)
	}
	industryStart := strings.Index(got, `<article class="ind-col">`)
	if industryStart < 0 {
		t.Fatal("industry entry missing")
	}
	industryEnd := strings.Index(got[industryStart:], `</article>`)
	if industryEnd < 0 {
		t.Fatal("industry entry closing tag missing")
	}
	industry := got[industryStart : industryStart+industryEnd]
	if summary, label, watch := strings.Index(industry, "產業內容"), strings.Index(industry, "觀察重點"), strings.Index(industry, "產業觀察內容"); summary < 0 || label <= summary || watch <= label {
		t.Fatalf("industry watch is missing or out of order: %s", industry)
	}
	stock := newsEntry(got, "n-1001")
	if summary, label, watch := strings.Index(stock, "短多摘要"), strings.Index(stock, "觀察重點"), strings.Index(stock, "短多觀察內容"); summary < 0 || label <= summary || watch <= label {
		t.Fatalf("stock watch is missing or out of order: %s", stock)
	}

	legacy := fullReport("2026-08-06")
	legacy.Industries[0].WatchMD = " \t\n"
	for i := range legacy.StockNews {
		legacy.StockNews[i].WatchMD = ""
	}
	legacyHTML := app.renderReport(&legacy, navToHistory)
	if strings.Contains(legacyHTML, `class="entry-watch"`) || strings.Contains(legacyHTML, "觀察重點") {
		t.Fatalf("legacy report rendered an entry watch placeholder: %s", legacyHTML)
	}
}

func TestHistoryPagingOrderingAndHistoricalReport(t *testing.T) {
	s, app := setupSite(t)
	for day := 1; day <= 25; day++ {
		date := fmt.Sprintf("2026-07-%02d", day)
		r := fullReport(date)
		r.Headline = "頭條 " + date
		saveReport(t, s, r)
	}
	at := sy.NewAppTest(func() { app.HistoryPage(1, "") })
	at.Run()
	got := renderedHTML(t, at)
	if strings.Count(got, `class="hist-item"`) != 10 || strings.Index(got, "2026-07-25") > strings.Index(got, "2026-07-24") || !strings.Contains(got, `?page=3`) {
		t.Fatalf("history page 1 ordering or pagination incorrect")
	}

	page3 := sy.NewAppTest(func() { app.HistoryPage(3, "") })
	page3.Run()
	if got := renderedHTML(t, page3); strings.Count(got, `class="hist-item"`) != 5 || !strings.Contains(got, "2026-07-05") || !strings.Contains(got, "2026-07-01") {
		t.Fatalf("history page 3 incorrect")
	}

	historical := sy.NewAppTest(func() { app.HistoryPage(1, "2026-07-05") })
	historical.Run()
	got = renderedHTML(t, historical)
	if !strings.Contains(got, "頭條 2026-07-05") || !strings.Contains(got, `data-section="stock-news"`) || !strings.Contains(got, Disclaimer) {
		t.Fatalf("historical report did not use shared report layout")
	}
}

// Taiwan market convention: bullish is red, bearish is green, in both themes.
// Asserting on the hue keeps a token swap from silently inverting the meaning.
func assertBullishIsRedAndBearishIsGreen(t *testing.T, got string) {
	t.Helper()
	for _, token := range []struct {
		name    string
		wantHue string
	}{{"--up", "red"}, {"--down", "green"}} {
		for _, value := range cssTokenValues(got, token.name) {
			if hue := dominantHue(value); hue != token.wantHue {
				t.Errorf("%s = %s reads as %s, want %s", token.name, value, hue, token.wantHue)
			}
		}
	}
}

// cssTokenValues collects every declared value of a custom property, so both the
// light and the dark theme block are checked.
func cssTokenValues(css, name string) []string {
	var values []string
	for _, chunk := range strings.Split(css, name+":")[1:] {
		end := strings.IndexAny(chunk, ";}")
		if end > 0 {
			values = append(values, strings.TrimSpace(chunk[:end]))
		}
	}
	return values
}

func dominantHue(hex string) string {
	if len(hex) != 7 || hex[0] != '#' {
		return "unparsed:" + hex
	}
	var rgb [3]int64
	for i := range rgb {
		v, err := strconv.ParseInt(hex[1+i*2:3+i*2], 16, 32)
		if err != nil {
			return "unparsed:" + hex
		}
		rgb[i] = v
	}
	switch {
	case rgb[0] > rgb[1]:
		return "red"
	case rgb[1] > rgb[0]:
		return "green"
	default:
		return "neutral"
	}
}

// Every call must render inside the column matching its direction, so swapping
// two lists cannot pass unnoticed.
func assertCallsLandInMatchingColumns(t *testing.T, got string) {
	t.Helper()
	columns := map[string]string{}
	for _, chunk := range strings.Split(got, `<div class="call-col `)[1:] {
		class, rest, ok := strings.Cut(chunk, `"`)
		if !ok {
			continue
		}
		title, _, _ := strings.Cut(rest[strings.Index(rest, "</span>")+len("</span>"):], "<span")
		columns[title] = class + "|" + rest
	}
	for _, want := range []struct{ title, class, stock string }{
		{"短期看漲", "up", "短多股"},
		{"短期看跌", "down", "短空股"},
		{"長期看好・可留意", "up", "長多股"},
		{"長期看壞・不建議", "down", "長空股"},
	} {
		column, ok := columns[want.title]
		if !ok {
			t.Errorf("column %q missing", want.title)
			continue
		}
		class, body, _ := strings.Cut(column, "|")
		if class != want.class {
			t.Errorf("column %q has class %q, want %q", want.title, class, want.class)
		}
		if !strings.Contains(body, want.stock) {
			t.Errorf("column %q does not contain %q", want.title, want.stock)
		}
	}
}

// Each detail entry carries the tag for its own call, with the colour class that
// matches the direction.
func assertStockNewsTags(t *testing.T, got string) {
	t.Helper()
	for _, want := range []struct{ id, class, label string }{
		{"n-1001", "up", "短期看漲"},
		{"n-1002", "down", "短期看跌"},
		{"n-1003", "up", "長期看好"},
		{"n-1004", "down", "長期看壞"},
	} {
		entry := newsEntry(got, want.id)
		if entry == "" {
			t.Errorf("news entry %s missing", want.id)
			continue
		}
		wantTag := `<span class="tag ` + want.class + `">`
		if !strings.Contains(entry, wantTag) {
			t.Errorf("news entry %s missing tag %s", want.id, wantTag)
		}
		if !strings.Contains(entry, want.label) {
			t.Errorf("news entry %s missing label %q", want.id, want.label)
		}
	}
}

func newsEntry(got, id string) string {
	start := strings.Index(got, `id="`+id+`"`)
	if start < 0 {
		return ""
	}
	end := strings.Index(got[start:], `</article>`)
	if end < 0 {
		return ""
	}
	return got[start : start+end]
}
