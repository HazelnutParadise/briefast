package site

import (
	"bytes"
	"context"
	"fmt"
	"html"
	"math"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/HazelnutParadise/briefast/internal/report"
	"github.com/HazelnutParadise/briefast/internal/store"
	sy "github.com/HazelnutParadise/syralit"
	"github.com/yuin/goldmark"
)

const Disclaimer = "本報告由 AI 自動彙整公開新聞產生，多空判斷為模型觀點，僅供參考，不構成任何投資建議。投資人應自行判斷並承擔風險。"

type Site struct {
	store   *store.Store
	version *sy.SharedVar[int]
	md      goldmark.Markdown
}

func New(s *store.Store) *Site {
	return &Site{
		store:   s,
		version: sy.Shared("briefast-report-version", 0),
		md:      goldmark.New(),
	}
}

func (s *Site) Notify() {
	s.version.Update(func(value int) int { return value + 1 })
}

func (s *Site) Home() {
	s.pageConfig()
	_ = s.version.Get()
	r, err := s.store.LatestReport(context.Background())
	if err == store.ErrNotFound {
		sy.HTML(styles + masthead(dateLineWaiting, navToHistory) + `<main class="briefast shell empty"><h2>尚無報告</h2><p>每日報告發布後會顯示於此。</p></main>` + footer())
		return
	}
	if err != nil {
		sy.Error("無法載入報告")
		return
	}
	sy.HTML(styles + s.renderReport(r, navToHistory) + footer())
}

func (s *Site) History() {
	page, _ := strconv.Atoi(sy.QueryParam("page"))
	if page < 1 {
		page = 1
	}
	s.HistoryPage(page, sy.QueryParam("date"))
}

func (s *Site) HistoryPage(page int, date string) {
	s.pageConfig()
	_ = s.version.Get()
	if date != "" {
		r, err := s.store.ReportByDate(context.Background(), date)
		if err == store.ErrNotFound {
			sy.HTML(styles + masthead(dateLineArchive, navToHistory) + `<main class="briefast shell empty"><h2>找不到這份報告</h2><a href="/history/">返回歷史報告</a></main>` + footer())
			return
		}
		if err != nil {
			sy.Error("無法載入歷史報告")
			return
		}
		sy.HTML(styles + s.renderReport(r, navToHistory) + footer())
		return
	}

	rows, err := s.store.ListReports(context.Background(), page, 10)
	if err != nil {
		sy.Error("無法載入歷史報告")
		return
	}
	count, err := s.store.CountReports(context.Background())
	if err != nil {
		sy.Error("無法載入歷史報告")
		return
	}
	sy.HTML(styles + masthead(dateLineArchive, navToHome) + renderHistory(rows, page, count) + footer())
}

func (s *Site) pageConfig() {
	sy.SetPageConfig(
		sy.PageTitle("Briefast｜每日產業與股市報告"),
		sy.PageLayout("wide"),
		sy.InitialSidebarState("collapsed"),
	)
}

func (s *Site) renderReport(r *report.Report, nav string) string {
	var b strings.Builder
	b.WriteString(masthead(reportDateLine(r), nav))
	b.WriteString(`<main class="briefast shell">`)
	b.WriteString(`<section class="lead">`)
	b.WriteString(`<article class="lead-main" data-section="overview"><p class="kicker">盤前總覽</p><h2 class="headline">`)
	b.WriteString(html.EscapeString(r.Headline))
	b.WriteString(`</h2><div class="standfirst md">`)
	b.WriteString(s.markdown(r.OverviewMD))
	b.WriteString(`</div></article>`)
	b.WriteString(`<aside class="watch" data-section="watch"><h2>今日觀察</h2><div class="md">`)
	b.WriteString(s.markdown(r.WatchMD))
	b.WriteString(`</div></aside></section>`)
	b.WriteString(renderCalls(r.Calls))
	b.WriteString(s.renderIndustries(r.Industries))
	b.WriteString(s.renderStockNews(r.StockNews))
	b.WriteString(`</main>`)
	return b.String()
}

func renderCalls(calls report.Calls) string {
	groups := []struct {
		class   string
		title   string
		horizon string
		items   []report.CallEntry
	}{
		{"up", "短期看漲", "短期", calls.ShortBull},
		{"down", "短期看跌", "短期", calls.ShortBear},
		{"up", "長期看好・可留意", "長期", calls.LongBull},
		{"down", "長期看壞・不建議", "長期", calls.LongBear},
	}
	var b strings.Builder
	b.WriteString(`<section class="section" data-section="calls"><div class="section-head"><h2>個股多空判斷</h2><span class="note">AI 從新聞判讀，每檔附一句理由與新聞來源</span></div><div class="calls">`)
	for _, group := range groups {
		triangle := "u"
		if group.class == "down" {
			triangle = "d"
		}
		fmt.Fprintf(&b, `<div class="call-col %s"><h3><span class="tri %s"></span>%s<span class="horizon">%s</span></h3>`, group.class, triangle, html.EscapeString(group.title), group.horizon)
		if len(group.items) == 0 {
			b.WriteString(`<p class="no-call">今日無明確判斷</p>`)
		}
		for _, item := range group.items {
			fmt.Fprintf(&b, `<article class="call-item"><div class="stock">%s<span class="code">%s</span></div><p class="why">%s</p><a class="more" href="#n-%s">新聞詳情 ↓</a></article>`,
				html.EscapeString(item.Name), html.EscapeString(item.Symbol), html.EscapeString(item.Reason), url.PathEscape(item.Symbol))
		}
		b.WriteString(`</div>`)
	}
	b.WriteString(`</div></section>`)
	return b.String()
}

func (s *Site) renderIndustries(industries []report.Industry) string {
	var b strings.Builder
	b.WriteString(`<section class="section" data-section="industries"><div class="section-head"><h2>產業新聞摘要</h2></div><div class="industries">`)
	if len(industries) == 0 {
		b.WriteString(`<p>今日無產業摘要。</p>`)
	}
	for _, industry := range industries {
		fmt.Fprintf(&b, `<article class="ind-col"><h3>%s</h3><div class="md">%s</div></article>`, html.EscapeString(industry.Name), s.markdown(industry.SummaryMD))
	}
	b.WriteString(`</div></section>`)
	return b.String()
}

func (s *Site) renderStockNews(news []report.StockNews) string {
	var b strings.Builder
	b.WriteString(`<section class="section" data-section="stock-news"><div class="section-head"><h2>個股新聞詳情</h2></div><div class="news-list">`)
	if len(news) == 0 {
		b.WriteString(`<p>今日無個股新聞詳情。</p>`)
	}
	for _, item := range news {
		fmt.Fprintf(&b, `<article class="news-entry" id="n-%s"><div class="news-id"><div class="stock">%s<span class="code">%s</span></div>`, url.PathEscape(item.Symbol), html.EscapeString(item.Name), html.EscapeString(item.Symbol))
		if label, class, triangle := callPresentation(item.Call); label != "" {
			fmt.Fprintf(&b, `<span class="tag %s"><span class="tri %s"></span>%s</span>`, class, triangle, label)
		}
		b.WriteString(`</div><div class="news-body md">`)
		b.WriteString(s.markdown(item.SummaryMD))
		if len(item.Sources) > 0 {
			b.WriteString(`<div class="srcs"><span class="lbl">來源</span>`)
			for _, source := range item.Sources {
				if href, ok := safeHTTPURL(source.URL); ok {
					fmt.Fprintf(&b, `<a href="%s" rel="noopener noreferrer" target="_blank">%s</a>`, html.EscapeString(href), html.EscapeString(source.Title))
				} else {
					fmt.Fprintf(&b, `<span>%s</span>`, html.EscapeString(source.Title))
				}
			}
			b.WriteString(`</div>`)
		}
		b.WriteString(`</div></article>`)
	}
	b.WriteString(`</div></section>`)
	return b.String()
}

func callPresentation(call string) (label, class, triangle string) {
	switch call {
	case report.CallShortBull:
		return "短期看漲", "up", "u"
	case report.CallShortBear:
		return "短期看跌", "down", "d"
	case report.CallLongBull:
		return "長期看好", "up", "u"
	case report.CallLongBear:
		return "長期看壞", "down", "d"
	default:
		return "", "", ""
	}
}

func (s *Site) markdown(value string) string {
	var b bytes.Buffer
	if err := s.md.Convert([]byte(value), &b); err != nil {
		return `<p>` + html.EscapeString(value) + `</p>`
	}
	return b.String()
}

func safeHTTPURL(value string) (string, bool) {
	u, err := url.Parse(value)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return "", false
	}
	return u.String(), true
}

const (
	navToHistory = `<a href="/history/">歷史報告 →</a>`
	navToHome    = `<a href="/">回首頁 →</a>`

	dateLineWaiting = `<b>等待今日晨報</b><span class="sep">｜</span><span>本報告由 AI 彙整公開新聞產生</span>`
	dateLineArchive = `<b>歷史報告</b><span class="sep">｜</span><span>本報告由 AI 彙整公開新聞產生</span>`
)

func reportDateLine(r *report.Report) string {
	return `<b>` + html.EscapeString(displayDate(r.Date)) + `</b><span class="sep">｜</span><span>` + html.EscapeString(displayGeneratedAt(r.GeneratedAt)) + ` 更新</span><span class="sep">｜</span><span>本報告由 AI 彙整公開新聞產生</span>`
}

func masthead(dateLine, nav string) string {
	return `<div class="briefast"><header><div class="topbar"><div class="shell topbar-in"><div class="edition"><span>台北・開盤前晨報</span><span>每個交易日開盤前出刊</span></div>` + nav + `</div></div><div class="shell masthead"><h1>Briefast</h1><div class="tagline">每日產業與股市報告</div></div><div class="masthead-rule"></div><div class="shell dateline">` + dateLine + `</div></header></div>`
}

func renderHistory(rows []store.ReportSummary, page, count int) string {
	if page < 1 {
		page = 1
	}
	pages := int(math.Ceil(float64(count) / 10))
	var b strings.Builder
	b.WriteString(`<main class="briefast shell"><section class="section history"><div class="section-head"><h2>歷史報告</h2><span class="note">點日期進入該日完整報告</span></div><div class="hist-list">`)
	if len(rows) == 0 {
		b.WriteString(`<p>尚無歷史報告。</p>`)
	}
	for _, row := range rows {
		fmt.Fprintf(&b, `<div class="hist-item"><a class="d" href="/history/?date=%s">%s</a><span class="t">%s</span></div>`, url.QueryEscape(row.Date), html.EscapeString(row.Date), html.EscapeString(row.Headline))
	}
	b.WriteString(`</div>`)
	if pages > 1 {
		b.WriteString(`<nav class="pager" aria-label="歷史報告分頁">`)
		for p := 1; p <= pages; p++ {
			if p == page {
				fmt.Fprintf(&b, `<span class="cur" aria-current="page">%d</span>`, p)
			} else {
				fmt.Fprintf(&b, `<a href="/history/?page=%d">%d</a>`, p, p)
			}
		}
		b.WriteString(`</nav>`)
	}
	b.WriteString(`</section></main>`)
	return b.String()
}

func footer() string {
	return `<footer class="briefast site-foot"><div class="shell"><p class="foot-brand">Briefast</p><p>` + html.EscapeString(Disclaimer) + `</p></div></footer>`
}

func displayDate(value string) string {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return value
	}
	weekdays := []string{"日", "一", "二", "三", "四", "五", "六"}
	return fmt.Sprintf("%d 年 %d 月 %d 日（%s）", t.Year(), t.Month(), t.Day(), weekdays[t.Weekday()])
}

func displayGeneratedAt(value string) string {
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return value
	}
	return "開盤前 " + t.Format("15:04")
}
