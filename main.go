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
	addr := listenAddr(appConfig())
	log.Printf("Briefast listening on http://%s", addr)
	return http.ListenAndServe(addr, handler)
}

// appConfig 是全站唯一的設定來源。ResolveConfig 會把 syralit.toml 與內建預設
// 套進來，這是嵌入模式下取得設定的正確途徑——sy.GetOption 只在 page function
// 執行期間反映服務中的設定，在這裡呼叫會拿到未解析的零值。
func appConfig() sy.Config {
	return sy.ResolveConfig(sy.Config{Title: "Briefast"})
}

// listenAddr 抽成獨立函式，讓測試不必啟動伺服器就能驗證綁定位址。
func listenAddr(cfg sy.Config) string {
	return fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
}

func newHandler(s *store.Store) http.Handler {
	cfg := appConfig()
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
