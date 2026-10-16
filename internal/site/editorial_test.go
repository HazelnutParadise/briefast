package site

import (
	"strings"
	"testing"

	"github.com/HazelnutParadise/briefast/internal/report"
	sy "github.com/HazelnutParadise/syralit"
)

func TestEditionGuideCountsSharesAndTargets(t *testing.T) {
	industries := []report.Industry{
		{Name: "科技", Events: make([]report.IndustryEvent, 2)},
		{Name: "金融", Events: make([]report.IndustryEvent, 1)},
		{Name: "傳產", Events: make([]report.IndustryEvent, 1)},
		{Name: "房地產"},
	}
	got := renderEditionGuide(industries)
	for _, want := range []string{`aria-label="本期產業新聞導覽"`, "本期新聞分布", `href="#industry-0"`, `href="#industry-1"`, `href="#industry-2"`, "2 則", "1 則", "50.00%", "25.00%", "sector-tech", "sector-finance", "sector-industry"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q: %s", want, got)
		}
	}
	if strings.Contains(got, "房地產") || strings.Contains(got, `href="#industry-3"`) {
		t.Error("zero-event sector rendered")
	}
	if strings.Count(got, "<a ") != 3 {
		t.Errorf("expected 3 sector links: %s", got)
	}
}

func TestEditionGuideEmptyUnknownAndOriginalIndexes(t *testing.T) {
	for _, industries := range [][]report.Industry{nil, {{Name: "科技"}}} {
		if got := renderEditionGuide(industries); got != "" {
			t.Errorf("empty guide = %q", got)
		}
	}
	label := `<script>alert("x")</script>`
	got := renderEditionGuide([]report.Industry{{Name: "科技"}, {Name: label, Events: make([]report.IndustryEvent, 1)}})
	for _, want := range []string{`href="#industry-1"`, "&lt;script&gt;", "100.00%", "1 則"} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q: %s", want, got)
		}
	}
	if strings.Contains(got, "<script>") || strings.Contains(got, "sector-art") {
		t.Errorf("unescaped name or mismatched art: %s", got)
	}
}

func TestEditionGuideIllustrationMapping(t *testing.T) {
	for name, class := range map[string]string{"科技": "sector-tech", "金融": "sector-finance", "傳產": "sector-industry", "房地產": "sector-property"} {
		got := renderIndustryIllustration(name)
		if !strings.Contains(got, class) || !strings.Contains(got, `aria-hidden="true"`) {
			t.Errorf("%s: %s", name, got)
		}
	}
	if got := renderIndustryIllustration(`unknown" onerror="alert(1)`); got != "" {
		t.Fatalf("unknown illustration: %s", got)
	}
}

func TestEditorialDisclosuresOnHomeAndHistory(t *testing.T) {
	s, app := setupSite(t)
	value := fullReport("2026-08-07")
	saveReport(t, s, value)
	for name, page := range map[string]func(){"home": app.Home, "history": func() { app.HistoryPage(1, value.Date) }} {
		t.Run(name, func(t *testing.T) {
			at := sy.NewAppTest(page)
			at.Run()
			got := renderedHTML(t, at)
			for _, want := range []string{`class="overview-details"`, "閱讀盤前總覽", "總覽內容", `id="industry-0"`, `href="#industry-0"`, "sector-tech", `class="event-details"`, "記憶體產業內容", "封測產業內容", Disclaimer} {
				if !strings.Contains(got, want) {
					t.Errorf("missing %q", want)
				}
			}
			if strings.Contains(got, `class="overview-details" open`) || strings.Contains(got, `class="event-details" open`) {
				t.Error("long body starts expanded")
			}
			for _, section := range []string{"overview", "watch", "calls", "industries", "stock-news"} {
				if !strings.Contains(got, `data-section="`+section+`"`) {
					t.Errorf("missing section %s", section)
				}
			}
		})
	}
}
