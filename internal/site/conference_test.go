package site

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/HazelnutParadise/briefast/internal/report"
	"github.com/HazelnutParadise/briefast/internal/store"
	sy "github.com/HazelnutParadise/syralit"
)

func pinToday(app *Site, date string) {
	t, _ := time.Parse("2006-01-02", date)
	app.now = func() time.Time {
		return time.Date(t.Year(), t.Month(), t.Day(), 7, 50, 0, 0, time.FixedZone("Asia/Taipei", 8*3600))
	}
}

func brief(symbol, name, heldOn, prediction string) report.Conference {
	return report.Conference{
		Symbol: symbol, Name: name, Market: report.MarketTWSE,
		HeldOn: heldOn, HeldAt: "14:00", Venue: "線上法說會", AnnouncedOn: "2026-10-10",
		Prediction: prediction, Headline: name + "標題", SummaryMD: "判斷依據內容", WatchMD: "- 觀察 2 奈米量產時程",
		Sources: []report.Source{{Title: "來源甲", URL: "https://example.com/a"}}, GeneratedAt: "2026-10-10T07:50:00+08:00",
	}
}

func saveBrief(t *testing.T, s *store.Store, c report.Conference) {
	t.Helper()
	payload, err := json.Marshal(c)
	if err != nil {
		t.Fatal(err)
	}
	if err := s.UpsertConference(context.Background(), c, payload); err != nil {
		t.Fatal(err)
	}
}

func settle(t *testing.T, s *store.Store, symbol, heldOn string, pre, post float64) {
	t.Helper()
	key, err := s.CreateAPIKey(context.Background(), "k-"+symbol+heldOn, "t-"+symbol+heldOn)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.SettleConference(context.Background(), symbol, heldOn, store.Settlement{PreCloseDate: "2026-09-01", PreClose: pre, PostCloseDate: "2026-09-03", PostClose: post}, key); err != nil {
		t.Fatal(err)
	}
}

func TestConferencesSectionPlacementAndOmission(t *testing.T) {
	s, app := setupSite(t)
	pinToday(app, "2026-10-14")
	saveReport(t, s, fullReport("2026-10-14"))
	at := sy.NewAppTest(app.Home)
	at.Run()
	got := renderedHTML(t, at)
	if strings.Contains(got, `data-section="conferences"`) || strings.Contains(got, "近期法說會") {
		t.Fatalf("empty conferences rendered a section: %s", got)
	}

	saveBrief(t, s, brief("2330", "台積電", "2026-10-16", report.PredictionBull))
	app.Notify()
	at.Run()
	got = renderedHTML(t, at)
	calls, confs, slot, industries := strings.Index(got, `data-section="calls"`), strings.Index(got, `data-section="conferences"`), strings.Index(got, `data-section="ad"`), strings.Index(got, `data-section="industries"`)
	if !(calls < confs && confs < slot && slot < industries) {
		t.Fatalf("conferences misplaced: calls=%d confs=%d ad=%d industries=%d", calls, confs, slot, industries)
	}
	last := -1
	for _, section := range []string{`data-section="overview"`, `data-section="watch"`, `data-section="calls"`, `data-section="industries"`, `data-section="stock-news"`} {
		if index := strings.Index(got, section); index <= last {
			t.Fatalf("fixed section %s out of order", section)
		} else {
			last = index
		}
	}

	// A past-dated conference alone must not keep the section alive.
	pinToday(app, "2026-10-17")
	at.Run()
	if strings.Contains(renderedHTML(t, at), `data-section="conferences"`) {
		t.Fatal("past conference kept the section on the homepage")
	}

	historical := sy.NewAppTest(func() { app.HistoryPage(1, "2026-10-14") })
	historical.Run()
	if strings.Contains(renderedHTML(t, historical), `data-section="conferences"`) {
		t.Fatal("history view rendered the conferences section")
	}
}

func TestConferenceCardsOrderTagsAndLinks(t *testing.T) {
	s, app := setupSite(t)
	pinToday(app, "2026-10-14")
	saveReport(t, s, fullReport("2026-10-14"))
	saveBrief(t, s, brief("2330", "台積電", "2026-10-16", report.PredictionBull))
	saveBrief(t, s, brief("1314", "中石化", "2026-10-14", report.PredictionNone))
	saveBrief(t, s, brief("2412", "中華電", "2026-10-16", report.PredictionBear))
	at := sy.NewAppTest(app.Home)
	at.Run()
	got := renderedHTML(t, at)
	section := got[strings.Index(got, `data-section="conferences"`):strings.Index(got, `data-section="ad"`)]

	cards := strings.Split(section, `<article class="conf-card">`)[1:]
	if len(cards) != 3 {
		t.Fatalf("cards = %d, want 3: %s", len(cards), section)
	}
	for i, want := range []struct{ name, tag, label, href string }{
		{"中石化", `<span class="tag neutral">`, "無法判斷", "/conference/?symbol=1314&amp;date=2026-10-14"},
		{"台積電", `<span class="tag up">`, "會後看漲", "/conference/?symbol=2330&amp;date=2026-10-16"},
		{"中華電", `<span class="tag down">`, "會後看跌", "/conference/?symbol=2412&amp;date=2026-10-16"},
	} {
		card := cards[i]
		for _, fragment := range []string{want.name, want.tag, want.label, want.name + "標題", "會前報告 →", "10 月 16 日", "14:00"} {
			if i == 0 && fragment == "10 月 16 日" {
				fragment = "10 月 14 日"
			}
			if !strings.Contains(card, fragment) {
				t.Errorf("card %d missing %q: %s", i, fragment, card)
			}
		}
		if !strings.Contains(card, `href="`+want.href+`"`) && !strings.Contains(card, `href="`+strings.ReplaceAll(want.href, "&amp;", "&")+`"`) {
			t.Errorf("card %d missing link %s: %s", i, want.href, card)
		}
	}
	if !strings.Contains(section, "AI 於法說會前") || strings.Contains(section, "近 ") || strings.Contains(section, "最近結算") {
		t.Fatalf("unsettled state rendered a tally or settled list: %s", section)
	}
	if !strings.Contains(styles, `.tag.neutral{color:var(--ink-soft);border:1px solid var(--rule)}`) {
		t.Fatal("neutral tag style missing")
	}
	assertBullishIsRedAndBearishIsGreen(t, got)
}

func TestConferenceAccuracyAndSettledList(t *testing.T) {
	s, app := setupSite(t)
	pinToday(app, "2026-10-14")
	saveReport(t, s, fullReport("2026-10-14"))
	saveBrief(t, s, brief("2330", "台積電", "2026-10-16", report.PredictionBull))
	// 7 settled: 5 directional (3 hit), 2 none.
	fixtures := []struct {
		symbol, heldOn, prediction string
		pre, post                  float64
	}{
		{"1", "2026-09-01", report.PredictionBull, 100, 110},
		{"2", "2026-09-02", report.PredictionBear, 100, 110},
		{"3", "2026-09-03", report.PredictionNone, 100, 120},
		{"4", "2026-09-04", report.PredictionBull, 100, 90},
		{"5", "2026-09-05", report.PredictionBear, 100, 90},
		{"6", "2026-09-06", report.PredictionNone, 100, 100},
		{"7", "2026-09-07", report.PredictionBull, 1080, 1095},
	}
	for _, f := range fixtures {
		saveBrief(t, s, brief(f.symbol, "公司"+f.symbol, f.heldOn, f.prediction))
		settle(t, s, f.symbol, f.heldOn, f.pre, f.post)
	}
	at := sy.NewAppTest(app.Home)
	at.Run()
	got := renderedHTML(t, at)
	section := got[strings.Index(got, `data-section="conferences"`):strings.Index(got, `data-section="ad"`)]
	if !strings.Contains(section, `<span class="accuracy">近 5 場命中 3 場</span>`) {
		t.Fatalf("accuracy line wrong: %s", section)
	}
	rows := strings.Split(section, `<div class="conf-settled-row">`)[1:]
	if len(rows) != 5 {
		t.Fatalf("settled rows = %d, want 5", len(rows))
	}
	if !strings.Contains(rows[0], "公司7") || !strings.Contains(rows[0], `<span class="change up">+1.39%</span>`) || !strings.Contains(rows[0], `<span class="verdict">命中</span>`) || !strings.Contains(rows[0], "2026-09-07") {
		t.Fatalf("newest settled row wrong: %s", rows[0])
	}
	if !strings.Contains(rows[1], "公司6") || strings.Contains(rows[1], "verdict") || !strings.Contains(rows[1], `<span class="change flat">+0.00%</span>`) {
		t.Fatalf("none row should carry no verdict: %s", rows[1])
	}
	if !strings.Contains(rows[2], `<span class="change down">-10.00%</span>`) || !strings.Contains(rows[2], "命中") {
		t.Fatalf("bear hit row wrong: %s", rows[2])
	}
	if strings.Contains(section, "公司1") {
		t.Fatal("settled list exceeded five rows")
	}
}

func TestConferenceAccuracyWindowStopsAtTwenty(t *testing.T) {
	s, app := setupSite(t)
	pinToday(app, "2026-10-14")
	saveReport(t, s, fullReport("2026-10-14"))
	saveBrief(t, s, brief("2330", "台積電", "2026-10-16", report.PredictionBull))
	for i := 0; i < 25; i++ {
		heldOn := time.Date(2026, 8, 1+i, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		post := 90.0
		if i >= 11 { // the newest 14 hit, the oldest 11 miss
			post = 110
		}
		symbol := "s" + heldOn
		saveBrief(t, s, brief(symbol, "公司", heldOn, report.PredictionBull))
		settle(t, s, symbol, heldOn, 100, post)
	}
	at := sy.NewAppTest(app.Home)
	at.Run()
	if got := renderedHTML(t, at); !strings.Contains(got, "近 20 場命中 14 場") {
		t.Fatalf("window tally wrong: %s", got[strings.Index(got, "近 "):strings.Index(got, "近 ")+30])
	}
}

func TestConferencePageRendersAllBlocks(t *testing.T) {
	s, app := setupSite(t)
	value := brief("2330", "台積電", "2026-10-16", report.PredictionBull)
	value.Fundamentals = &report.Fundamentals{
		Revenue: &report.RevenueFigures{Month: "2026-08", Current: 13744103, PrevMonth: 13382706, LastYearMonth: 13535929, YTD: 85211435, LastYearYTD: 90000000},
		Quarter: &report.QuarterFigures{Label: "2026Q2 累計", Revenue: 71289957, OperatingIncome: 5170177, NetIncome: 4569799, EPS: "0.38"},
	}
	value.Chips = chipsFixture()
	value.Sources = append(value.Sources, report.Source{Title: "來源乙", URL: "https://example.com/b"}, report.Source{Title: "來源丙", URL: "https://example.com/c"})
	saveBrief(t, s, value)
	at := sy.NewAppTest(func() { app.ConferencePage("2330", "2026-10-16") })
	at.Run()
	got := renderedHTML(t, at)
	for _, want := range []string{
		`<span class="tag up">`, "會後看漲", "台積電標題", "判斷依據", "判斷依據內容",
		"法說會 2026 年 10 月 16 日（五） 14:00", "線上法說會",
		"基本面", "當月營收（2026-08）", "13,744,103", "13,382,706",
		`<span class="change up">較上月 +2.70%</span>`, `<span class="change up">較去年同月 +1.54%</span>`, `<span class="change down">累計較去年 -5.32%</span>`,
		"2026Q2 累計", "71,289,957", "0.38 元",
		"觀察重點", "觀察 2 奈米量產時程", "籌碼面", "+54,759 張",
		`>來源甲</a>`, `>來源乙</a>`, `>來源丙</a>`,
		"報告產生 2026-10-10 07:50", "公告日期 2026-10-10", Disclaimer, `href="/"`,
	} {
		if !strings.Contains(got, want) {
			t.Errorf("conference page missing %q", want)
		}
	}
	if strings.Contains(got, "已結算") || strings.Contains(got, "會前收盤") {
		t.Fatal("unsettled brief rendered a result band")
	}
	if strings.Contains(got, archiveStyles) {
		t.Fatal("conference page used the archive paper")
	}
}

func TestConferencePageResultBandAndOmittedFundamentals(t *testing.T) {
	s, app := setupSite(t)
	saveBrief(t, s, brief("2330", "台積電", "2026-10-16", report.PredictionBull))
	key, _ := s.CreateAPIKey(context.Background(), "k", "t")
	if _, err := s.SettleConference(context.Background(), "2330", "2026-10-16", store.Settlement{PreCloseDate: "2026-10-15", PreClose: 1080, PostCloseDate: "2026-10-17", PostClose: 1095}, key); err != nil {
		t.Fatal(err)
	}
	at := sy.NewAppTest(func() { app.ConferencePage("2330", "2026-10-16") })
	at.Run()
	got := renderedHTML(t, at)
	for _, want := range []string{"已結算", "會前收盤（10/15）1,080 元", "會後收盤（10/17）1,095 元", `<span class="change up">+1.39%</span>`, `<span class="verdict">命中</span>`} {
		if !strings.Contains(got, want) {
			t.Errorf("result band missing %q: %s", want, got)
		}
	}
	if strings.Contains(got, `data-section="fundamentals"`) || strings.Contains(got, "基本面") {
		t.Fatal("brief without fundamentals rendered the block")
	}

	saveBrief(t, s, brief("9", "無方向", "2026-10-16", report.PredictionNone))
	if _, err := s.SettleConference(context.Background(), "9", "2026-10-16", store.Settlement{PreCloseDate: "2026-10-15", PreClose: 10, PostCloseDate: "2026-10-17", PostClose: 11}, key); err != nil {
		t.Fatal(err)
	}
	none := sy.NewAppTest(func() { app.ConferencePage("9", "2026-10-16") })
	none.Run()
	if !strings.Contains(renderedHTML(t, none), "未列入統計") {
		t.Fatal("prediction none settled without 未列入統計")
	}
}

func TestConferencePageNotFound(t *testing.T) {
	s, app := setupSite(t)
	saveBrief(t, s, brief("2330", "台積電", "2026-10-16", report.PredictionBull))
	for _, tc := range [][2]string{{"9999", "2026-01-01"}, {"", "2026-10-16"}, {"2330", "2026-13-01"}, {"2330", ""}} {
		at := sy.NewAppTest(func() { app.ConferencePage(tc[0], tc[1]) })
		at.Run()
		got := renderedHTML(t, at)
		if !strings.Contains(got, "找不到這場法說會的報告") || !strings.Contains(got, `href="/"`) || !strings.Contains(got, Disclaimer) {
			t.Errorf("not-found state wrong for %v: %s", tc, got)
		}
	}
}
