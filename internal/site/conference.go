package site

import (
	"context"
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"

	"github.com/HazelnutParadise/briefast/internal/report"
	"github.com/HazelnutParadise/briefast/internal/store"
	sy "github.com/HazelnutParadise/syralit"
)

// accuracyWindow caps the tally so early misses do not haunt the number forever
// and one lucky week cannot inflate it either.
const accuracyWindow = 20

// settledListLimit is how many recently settled conferences sit under the cards.
const settledListLimit = 5

// renderConferences builds the homepage section. It is omitted entirely when
// nothing is upcoming: the section describes "now", so an empty frame would
// only advertise absence.
func (s *Site) renderConferences() string {
	ctx := context.Background()
	today := report.TaipeiDate(s.now())
	upcoming, err := s.store.ListUpcomingConferences(ctx, today)
	if err != nil {
		return `<section class="section" data-section="conferences"><div class="section-head"><h2>近期法說會</h2></div><p class="no-call">無法載入近期法說會</p></section>`
	}
	if len(upcoming) == 0 {
		return ""
	}
	settled, err := s.store.ListSettledConferences(ctx, accuracyWindow)
	if err != nil {
		settled = nil
	}

	var b strings.Builder
	b.WriteString(`<section class="section" data-section="conferences"><div class="section-head"><h2>近期法說會</h2><span class="note">AI 於法說會前依基本面與近期消息預測會後方向</span>`)
	if line := accuracyLine(settled); line != "" {
		b.WriteString(`<span class="accuracy">` + line + `</span>`)
	}
	b.WriteString(`</div><div class="confs">`)
	for _, item := range upcoming {
		label, class := predictionPresentation(item.Prediction)
		fmt.Fprintf(&b, `<article class="conf-card"><div class="conf-when">%s</div><div class="stock">%s<span class="code">%s</span></div>%s<p class="why">%s</p><a class="more" href="%s">會前報告 →</a></article>`,
			html.EscapeString(displayConferenceTime(item.HeldOn, item.HeldAt)), html.EscapeString(item.Name), html.EscapeString(item.Symbol),
			predictionTag(label, class), html.EscapeString(briefHeadline(ctx, s.store, item)), conferenceHref(item.Symbol, item.HeldOn))
	}
	b.WriteString(`</div>`)
	b.WriteString(renderSettledList(settled))
	b.WriteString(`</section>`)
	return b.String()
}

// briefHeadline decodes only the payloads of the few upcoming cards; the
// listing itself stays index-only.
func briefHeadline(ctx context.Context, s *store.Store, item store.ConferenceSummary) string {
	full, err := s.ConferenceByKey(ctx, item.Symbol, item.HeldOn)
	if err != nil {
		return ""
	}
	return full.Brief.Headline
}

// accuracyLine counts hits among the most recent settled directional briefs.
func accuracyLine(settled []store.ConferenceSummary) string {
	total, hits := 0, 0
	for _, item := range settled {
		if item.Settlement == nil || item.Prediction == report.PredictionNone {
			continue
		}
		if total == accuracyWindow {
			break
		}
		total++
		if item.Settlement.Outcome == report.OutcomeHit {
			hits++
		}
	}
	if total == 0 {
		return ""
	}
	return fmt.Sprintf("近 %d 場命中 %d 場", total, hits)
}

func renderSettledList(settled []store.ConferenceSummary) string {
	if len(settled) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<div class="conf-settled"><span class="entry-watch-label">最近結算</span>`)
	for i, item := range settled {
		if i == settledListLimit || item.Settlement == nil {
			break
		}
		label, class := predictionPresentation(item.Prediction)
		verdict := verdictLabel(item.Settlement.Outcome)
		if verdict != "" {
			verdict = `<span class="verdict">` + verdict + `</span>`
		}
		fmt.Fprintf(&b, `<div class="conf-settled-row"><a class="stock" href="%s">%s<span class="code">%s</span></a><span class="conf-date">%s</span>%s%s%s</div>`,
			conferenceHref(item.Symbol, item.HeldOn), html.EscapeString(item.Name), html.EscapeString(item.Symbol), html.EscapeString(item.HeldOn),
			predictionTag(label, class), changeSpan(item.Settlement.PreClose, item.Settlement.PostClose), verdict)
	}
	b.WriteString(`</div>`)
	return b.String()
}

func verdictLabel(outcome string) string {
	switch outcome {
	case report.OutcomeHit:
		return "命中"
	case report.OutcomeMiss:
		return "落空"
	default:
		return ""
	}
}

func predictionPresentation(prediction string) (label, class string) {
	switch prediction {
	case report.PredictionBull:
		return "會後看漲", "up"
	case report.PredictionBear:
		return "會後看跌", "down"
	default:
		return "無法判斷", "neutral"
	}
}

func predictionTag(label, class string) string {
	switch class {
	case "up":
		return `<span class="tag up"><span class="tri u"></span>` + label + `</span>`
	case "down":
		return `<span class="tag down"><span class="tri d"></span>` + label + `</span>`
	default:
		return `<span class="tag neutral">` + label + `</span>`
	}
}

// changeSpan renders the close-to-close move; a flat move takes the neutral
// colour because neither direction happened.
func changeSpan(pre, post float64) string {
	if pre <= 0 {
		return ""
	}
	pct := (post - pre) / pre * 100
	class := "flat"
	switch {
	case pct > 0:
		class = "up"
	case pct < 0:
		class = "down"
	}
	return fmt.Sprintf(`<span class="change %s">%s</span>`, class, formatPct(pct))
}

func formatPct(pct float64) string {
	return fmt.Sprintf("%+.2f%%", pct)
}

func conferenceHref(symbol, heldOn string) string {
	return "/conference/?symbol=" + url.QueryEscape(symbol) + "&date=" + url.QueryEscape(heldOn)
}

func displayConferenceTime(heldOn, heldAt string) string {
	out := displayDate(heldOn)
	if heldAt != "" {
		out += " " + heldAt
	}
	return out
}

// Conference is the page function mounted at /conference/.
func (s *Site) Conference() {
	s.ConferencePage(sy.QueryParam("symbol"), sy.QueryParam("date"))
}

func (s *Site) ConferencePage(symbol, date string) {
	s.pageConfig()
	_ = s.version.Get()
	symbol, date = strings.TrimSpace(symbol), strings.TrimSpace(date)
	if symbol == "" || len(symbol) > 10 || !validPageDate(date) {
		s.conferenceNotFound()
		return
	}
	item, err := s.store.ConferenceByKey(context.Background(), symbol, date)
	if err == store.ErrNotFound {
		s.conferenceNotFound()
		return
	}
	if err != nil {
		sy.Error("無法載入法說會報告")
		return
	}
	sy.HTML(styles + s.renderConference(item) + footer())
}

func (s *Site) conferenceNotFound() {
	sy.HTML(styles + masthead(dateLineConference, navToHome) + `<main class="briefast shell empty"><h2>找不到這場法說會的報告</h2><p><a href="/">回首頁</a></p></main>` + footer())
}

func validPageDate(value string) bool {
	if len(value) != len("2006-01-02") {
		return false
	}
	parsed, err := time.Parse("2006-01-02", value)
	return err == nil && parsed.Format("2006-01-02") == value
}

const dateLineConference = `<b>法說會前預測</b><span class="sep">｜</span><span>本報告由 AI 彙整公開資料產生</span>`

func (s *Site) renderConference(item *store.Conference) string {
	brief := item.Brief
	label, class := predictionPresentation(item.Prediction)
	var b strings.Builder
	b.WriteString(masthead(`<b>`+html.EscapeString(brief.Name)+` `+html.EscapeString(brief.Symbol)+`</b><span class="sep">｜</span><span>法說會前預測</span><span class="sep">｜</span><span>本報告由 AI 彙整公開資料產生</span>`, navToHome))
	b.WriteString(resultBand(item))
	b.WriteString(`<main class="briefast shell conf-page">`)
	fmt.Fprintf(&b, `<section class="conf-head" data-section="conference"><p class="kicker">法說會前預測</p><div class="stock">%s<span class="code">%s</span></div><h2 class="headline">%s</h2><div class="conf-meta"><span>法說會 %s</span>`,
		html.EscapeString(brief.Name), html.EscapeString(brief.Symbol), html.EscapeString(brief.Headline), html.EscapeString(displayConferenceTime(brief.HeldOn, brief.HeldAt)))
	if strings.TrimSpace(brief.Venue) != "" {
		b.WriteString(`<span class="sep">・</span><span>` + html.EscapeString(brief.Venue) + `</span>`)
	}
	b.WriteString(`</div>` + predictionTag(label, class) + `</section>`)

	b.WriteString(`<section class="section" data-section="basis"><div class="section-head"><h2>判斷依據</h2></div><div class="news-body md">` + s.markdown(brief.SummaryMD) + `</div></section>`)
	b.WriteString(renderFundamentals(brief.Fundamentals))
	b.WriteString(`<section class="section" data-section="detail">`)
	b.WriteString(s.renderEntryWatch(brief.WatchMD))
	b.WriteString(renderChips(brief.Chips))
	if len(brief.Sources) > 0 {
		b.WriteString(`<div class="srcs"><span class="lbl">來源</span>`)
		for _, source := range brief.Sources {
			if href, ok := safeHTTPURL(source.URL); ok {
				fmt.Fprintf(&b, `<a href="%s" rel="noopener noreferrer" target="_blank">%s</a>`, html.EscapeString(href), html.EscapeString(source.Title))
			} else {
				fmt.Fprintf(&b, `<span>%s</span>`, html.EscapeString(source.Title))
			}
		}
		b.WriteString(`</div>`)
	}
	fmt.Fprintf(&b, `<p class="conf-dates">報告產生 %s<span class="sep">・</span>公告日期 %s</p>`, html.EscapeString(displayGeneratedFull(brief.GeneratedAt)), html.EscapeString(brief.AnnouncedOn))
	b.WriteString(`</section></main>`)
	return b.String()
}

// resultBand mirrors the archive notice band: a dark strip that states the
// settled figures before the reader gets to the pre-event reasoning.
func resultBand(item *store.Conference) string {
	st := item.Settlement
	if st == nil {
		return ""
	}
	verdict := verdictLabel(st.Outcome)
	if verdict == "" {
		verdict = "未列入統計"
	}
	return `<div class="briefast archive-note result-band"><div class="shell archive-note-in"><span class="archive-flag">已結算</span><span class="archive-text">會前收盤（` + html.EscapeString(shortDate(st.PreCloseDate)) + `）` + formatPrice(st.PreClose) + ` 元 → 會後收盤（` + html.EscapeString(shortDate(st.PostCloseDate)) + `）` + formatPrice(st.PostClose) + ` 元</span>` + changeSpan(st.PreClose, st.PostClose) + `<span class="verdict">` + verdict + `</span></div></div>`
}

func formatPrice(value float64) string {
	text := fmt.Sprintf("%.2f", value)
	text = strings.TrimRight(strings.TrimRight(text, "0"), ".")
	whole, frac, hasFrac := strings.Cut(text, ".")
	out := groupThousandsString(whole)
	if hasFrac {
		out += "." + frac
	}
	return out
}

func groupThousandsString(digits string) string {
	for i := len(digits) - 3; i > 0; i -= 3 {
		digits = digits[:i] + "," + digits[i:]
	}
	return digits
}

func shortDate(value string) string {
	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return value
	}
	return fmt.Sprintf("%d/%d", t.Month(), t.Day())
}

func displayGeneratedFull(value string) string {
	t, ok := parseGeneratedAt(value)
	if !ok {
		return value
	}
	return t.Format("2006-01-02 15:04")
}

// renderFundamentals draws the revenue comparisons and the latest cumulative
// quarter. Percentages are computed here so every brief shows the same maths.
func renderFundamentals(f *report.Fundamentals) string {
	if f == nil || (f.Revenue == nil && f.Quarter == nil) {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<section class="section" data-section="fundamentals"><div class="section-head"><h2>基本面</h2><span class="note">單位：新台幣千元</span></div>`)
	if r := f.Revenue; r != nil {
		b.WriteString(`<table class="fund"><tbody>`)
		fundRow(&b, "當月營收（"+html.EscapeString(r.Month)+"）", r.Current, nil)
		fundRow(&b, "上月營收", r.PrevMonth, pctCell(r.Current, r.PrevMonth, "較上月"))
		fundRow(&b, "去年同月營收", r.LastYearMonth, pctCell(r.Current, r.LastYearMonth, "較去年同月"))
		fundRow(&b, "今年累計營收", r.YTD, nil)
		fundRow(&b, "去年累計營收", r.LastYearYTD, pctCell(r.YTD, r.LastYearYTD, "累計較去年"))
		b.WriteString(`</tbody></table>`)
	}
	if q := f.Quarter; q != nil {
		b.WriteString(`<table class="fund quarter"><tbody>`)
		fmt.Fprintf(&b, `<tr><th colspan="3">%s</th></tr>`, html.EscapeString(q.Label))
		fundRow(&b, "營業收入", q.Revenue, nil)
		fundRow(&b, "營業利益", q.OperatingIncome, nil)
		fundRow(&b, "本期淨利", q.NetIncome, nil)
		fmt.Fprintf(&b, `<tr><th>每股盈餘</th><td>%s 元</td><td></td></tr>`, html.EscapeString(q.EPS))
		b.WriteString(`</tbody></table>`)
	}
	b.WriteString(`</section>`)
	return b.String()
}

func fundRow(b *strings.Builder, label string, value int64, pct *string) {
	cell := ""
	if pct != nil {
		cell = *pct
	}
	fmt.Fprintf(b, `<tr><th>%s</th><td>%s</td><td>%s</td></tr>`, label, groupThousands(value), cell)
}

// pctCell renders "較上月 +2.70%" coloured by direction; a zero base yields
// nothing rather than a division by zero dressed up as a number.
func pctCell(current, base int64, label string) *string {
	if base == 0 {
		return nil
	}
	pct := (float64(current) - float64(base)) / float64(base) * 100
	class := "flat"
	switch {
	case pct > 0:
		class = "up"
	case pct < 0:
		class = "down"
	}
	out := fmt.Sprintf(`<span class="change %s">%s %s</span>`, class, label, formatPct(pct))
	return &out
}
