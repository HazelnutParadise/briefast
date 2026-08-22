package seo

import (
	"encoding/xml"
	"net/http"
	"strings"
	"time"
)

// RobotsHandler 回傳 robots.txt。公開頁全部開放，後台與 API 擋掉，並用絕對網址
// 指向 sitemap。
func (d Deps) RobotsHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		base := d.Config.BaseURL(r)
		body := strings.Join([]string{
			"User-agent: *",
			"Allow: /",
			"Disallow: /admin/",
			"Disallow: /api/",
			"",
			"Sitemap: " + base + "/sitemap.xml",
			"",
		}, "\n")
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte(body))
	})
}

type sitemapURL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

type sitemapURLSet struct {
	XMLName xml.Name     `xml:"urlset"`
	Xmlns   string       `xml:"xmlns,attr"`
	URLs    []sitemapURL `xml:"url"`
}

// SitemapHandler 回傳涵蓋首頁、歷史列表與每一份報告的 sitemap。報告一天一份，
// 全量輸出仍遠低於 sitemap 的筆數上限，不需要分頁。
func (d Deps) SitemapHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		count, err := d.Reports.CountReports(ctx)
		if err != nil {
			http.Error(w, "無法產生 sitemap", http.StatusInternalServerError)
			return
		}
		var reports []summary
		if count > 0 {
			rows, err := d.Reports.ListReports(ctx, 1, count)
			if err != nil {
				http.Error(w, "無法產生 sitemap", http.StatusInternalServerError)
				return
			}
			for _, row := range rows {
				reports = append(reports, summary{date: row.Date, generatedAt: row.GeneratedAt})
			}
		}

		base := d.Config.BaseURL(r)
		newest := ""
		if len(reports) > 0 {
			newest = lastModOf(reports[0])
		}
		set := sitemapURLSet{
			Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
			URLs: []sitemapURL{
				{Loc: base + "/", LastMod: newest},
				{Loc: base + "/history/", LastMod: newest},
			},
		}
		for _, item := range reports {
			set.URLs = append(set.URLs, sitemapURL{
				Loc:     base + "/history/?date=" + item.date,
				LastMod: lastModOf(item),
			})
		}

		document, err := xml.MarshalIndent(set, "", "  ")
		if err != nil {
			http.Error(w, "無法產生 sitemap", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/xml; charset=utf-8")
		_, _ = w.Write([]byte(xml.Header))
		_, _ = w.Write(document)
		_, _ = w.Write([]byte("\n"))
	})
}

// summary 只留 sitemap 用得到的兩個欄位。
type summary struct {
	date        string
	generatedAt string
}

// lastModOf 用報告的產生時間當 lastmod，格式不合就退回報告日期。
func lastModOf(item summary) string {
	if t, err := time.Parse(time.RFC3339, item.generatedAt); err == nil {
		return t.UTC().Format(time.RFC3339)
	}
	return item.date
}
