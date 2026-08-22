package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/HazelnutParadise/briefast/internal/admin"
	briefapi "github.com/HazelnutParadise/briefast/internal/api"
	"github.com/HazelnutParadise/briefast/internal/seo"
	"github.com/HazelnutParadise/briefast/internal/site"
	"github.com/HazelnutParadise/briefast/internal/store"
	sy "github.com/HazelnutParadise/syralit"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	databasePath := os.Getenv("BRIEFAST_DB_PATH")
	if databasePath == "" {
		databasePath = "data/briefast.db"
	}
	s, err := store.Open(databasePath)
	if err != nil {
		return err
	}
	defer s.Close()

	handler := newHandler(s)
	host, _ := sy.GetOption("server.host").(string)
	port, _ := sy.GetOption("server.port").(int)
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Printf("Briefast listening on http://%s", addr)
	return http.ListenAndServe(addr, handler)
}

func newHandler(s *store.Store) http.Handler {
	cfg := sy.Config{Title: "Briefast"}
	public := site.New(s)
	adminApp := admin.New(s)
	home := sy.Handler(cfg, public.Home)
	history := sy.Handler(cfg, public.History)
	adminHandler := sy.Handler(cfg, adminApp.Page)

	// 公開頁包一層 SEO 中介軟體改寫 head；後台與 API 不套。
	crawler := seo.Deps{Reports: s, Config: seo.Config{SiteURL: os.Getenv("BRIEFAST_SITE_URL")}}

	mux := http.NewServeMux()
	mux.Handle("/api/report", briefapi.NewReportHandler(s, public))
	mux.Handle("/api/report/{date}", briefapi.NewReadHandler(s))
	mux.Handle("GET /robots.txt", crawler.RobotsHandler())
	mux.Handle("GET /sitemap.xml", crawler.SitemapHandler())
	mux.Handle("/", crawler.Middleware(home))
	mux.Handle("/history/", crawler.Middleware(http.StripPrefix("/history", history)))
	mux.Handle("/admin/", http.StripPrefix("/admin", adminHandler))
	mux.HandleFunc("GET /history", trailingSlash("/history/"))
	mux.HandleFunc("GET /admin", trailingSlash("/admin/"))
	return mux
}

func trailingSlash(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path, http.StatusTemporaryRedirect)
	}
}
