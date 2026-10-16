package site

import (
	"fmt"
	"html"
	"strings"

	"github.com/HazelnutParadise/briefast/internal/report"
)

// renderEditionGuide links actual report coverage to its industry sections.
// Empty input renders nothing; counts and shares describe news, never sentiment.
func renderEditionGuide(industries []report.Industry) string {
	total := 0
	for _, ind := range industries {
		total += len(ind.Events)
	}
	if total == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<nav class="edition-guide" aria-label="本期產業新聞導覽"><div class="edition-guide-head"><strong>本期新聞分布</strong><span>依報導則數</span></div><div class="edition-sectors">`)
	for i, ind := range industries {
		count := len(ind.Events)
		if count == 0 {
			continue
		}
		fmt.Fprintf(&b, `<a class="edition-sector" href="#industry-%d">%s<strong>%s</strong><span class="edition-count">%d 則</span><span class="edition-track" aria-hidden="true"><span class="edition-share" style="width:%.2f%%"></span></span></a>`,
			i, renderIndustryIllustration(ind.Name), html.EscapeString(ind.Name), count, float64(count)/float64(total)*100)
	}
	b.WriteString(`</div></nav>`)
	return b.String()
}

// renderIndustryIllustration returns the matching newspaper vignette.
// Unknown labels render no art; arbitrary input never enters CSS or attributes.
func renderIndustryIllustration(name string) string {
	class := map[string]string{
		"科技":   "sector-tech",
		"金融":   "sector-finance",
		"傳產":   "sector-industry",
		"房地產":  "sector-property",
	}[name]
	if class == "" {
		return ""
	}
	return fmt.Sprintf(`<div class="sector-art %s" aria-hidden="true"></div>`, class)
}
