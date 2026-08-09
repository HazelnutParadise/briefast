package admin

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"html"
	"strings"

	"github.com/HazelnutParadise/briefast/internal/store"
	sy "github.com/HazelnutParadise/syralit"
)

type Option func(*App)

func WithPasswordSource(source func() string) Option {
	return func(a *App) { a.password = source }
}

func WithTokenGenerator(generator func() (string, error)) Option {
	return func(a *App) { a.token = generator }
}

type App struct {
	store    *store.Store
	password func() string
	token    func() (string, error)
}

func New(s *store.Store, opts ...Option) *App {
	a := &App{
		store: s,
		password: func() string {
			return sy.Secrets("BRIEFAST_ADMIN_PASSWORD")
		},
		token: generateToken,
	}
	for _, opt := range opts {
		opt(a)
	}
	return a
}

func (a *App) Page() {
	sy.SetPageConfig(
		sy.PageTitle("Briefast 後台"),
		sy.PageLayout("wide"),
		sy.InitialSidebarState("collapsed"),
	)
	configured := a.password()
	if configured == "" {
		sy.HTML(adminStyles + `<main class="admin-shell"><h1>Briefast 後台</h1><p class="config-error">尚未設定管理員密碼。請設定 BRIEFAST_ADMIN_PASSWORD 後重新啟動。</p></main>`)
		return
	}

	authenticated := sy.State("briefast-admin-authenticated", false)
	if !authenticated.Get() {
		sy.HTML(adminStyles + `<main class="admin-shell login"><p class="eyebrow">ADMINISTRATION</p><h1>Briefast 後台登入</h1></main>`)
		password := sy.PasswordInput("管理員密碼", sy.Key("admin-password"))
		if sy.Button("登入", sy.Key("admin-login")) {
			if secureEqual(password, configured) {
				authenticated.Set(true)
				sy.Rerun()
			} else {
				sy.Error("密碼錯誤")
			}
		}
		return
	}

	a.renderDashboard()
}

func (a *App) renderDashboard() {
	sy.HTML(adminStyles + `<header class="admin-head"><div class="admin-shell"><p class="eyebrow">ADMINISTRATION</p><h1>Briefast 後台</h1><a href="/">返回晨報</a></div></header>`)

	name := sy.TextInput("Key 名稱", sy.Key("key-name"), sy.Placeholder("例如 cowork-daily"))
	if sy.Button("建立 API key", sy.Key("create-key")) {
		if strings.TrimSpace(name) == "" {
			sy.Error("請輸入 Key 名稱")
		} else {
			token, err := a.token()
			if err != nil {
				sy.Error("無法產生 API key")
			} else if _, err := a.store.CreateAPIKey(context.Background(), name, token); err != nil {
				sy.Error("無法建立 API key")
			} else {
				sy.ResetWidget("key-name")
				sy.Success("API key 已建立")
			}
		}
	}

	keys, err := a.store.ListAPIKeys(context.Background())
	if err != nil {
		sy.Error("無法載入 API keys")
		return
	}
	sy.HTML(renderKeys(keys))
	for _, key := range keys {
		if key.Active() && sy.Button("撤銷 "+key.Name, sy.Key(fmt.Sprintf("revoke-key-%d", key.ID)), sy.ButtonType("secondary")) {
			if err := a.store.RevokeAPIKey(context.Background(), key.ID); err != nil {
				sy.Error("無法撤銷 API key")
			} else {
				sy.Rerun()
			}
		}
	}

	logs, err := a.store.ListUpdateLogs(context.Background(), 200)
	if err != nil {
		sy.Error("無法載入更新紀錄")
		return
	}
	sy.HTML(renderLogs(logs))
}

func generateToken() (string, error) {
	random := make([]byte, 32)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}
	return "brf_" + base64.RawURLEncoding.EncodeToString(random), nil
}

func secureEqual(got, want string) bool {
	if len(got) != len(want) {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(want)) == 1
}

func renderKeys(keys []store.APIKey) string {
	var b strings.Builder
	b.WriteString(`<main class="admin-shell"><section><div class="section-head"><h2>API key 管理</h2><span>完整 token 會持續顯示，請保護資料庫檔案。</span></div><div class="key-list">`)
	if len(keys) == 0 {
		b.WriteString(`<p>尚未建立 API key。</p>`)
	}
	for _, key := range keys {
		status, class := "有效", "active"
		if !key.Active() {
			status, class = "已撤銷", "revoked"
		}
		fmt.Fprintf(&b, `<article class="key-row %s" data-key-id="%d"><div><strong>%s</strong><span class="status">%s</span><time>%s</time></div><code>%s</code></article>`,
			class, key.ID, html.EscapeString(key.Name), status, key.CreatedAt.Format("2006-01-02 15:04:05"), html.EscapeString(key.Token))
	}
	b.WriteString(`</div></section></main>`)
	return b.String()
}

func renderLogs(logs []store.UpdateLog) string {
	var b strings.Builder
	b.WriteString(`<main class="admin-shell"><section><div class="section-head"><h2>報告更新紀錄</h2></div><div class="log-list">`)
	if len(logs) == 0 {
		b.WriteString(`<p>尚無更新紀錄。</p>`)
	}
	for _, item := range logs {
		fmt.Fprintf(&b, `<article class="log-row" data-log-id="%d"><time>%s</time><strong>%s</strong><span>%s</span><code>%s</code>`,
			item.ID, item.At.Format("2006-01-02 15:04:05"), html.EscapeString(item.KeyName), html.EscapeString(item.ReportDate), html.EscapeString(item.Action))
		if item.Detail != "" {
			fmt.Fprintf(&b, `<p>%s</p>`, html.EscapeString(item.Detail))
		}
		b.WriteString(`</article>`)
	}
	b.WriteString(`</div></section></main>`)
	return b.String()
}

const adminStyles = `<style>
:root{--paper:#F9F1E7;--ink:#33302E;--soft:#66605A;--rule:#E0D8CA;--accent:#990F3D;--sans:"Noto Sans TC","PingFang TC",sans-serif;--serif:"Noto Serif TC","PMingLiU",serif}@media(prefers-color-scheme:dark){:root{--paper:#201D1A;--ink:#EDE6DC;--soft:#A89F94;--rule:#38342D;--accent:#E5697F}}*{border-radius:0!important}.sy-main,.sy-app{max-width:none!important;background:var(--paper)!important}.admin-shell{max-width:1120px;margin:0 auto;padding:28px;color:var(--ink);font-family:var(--sans)}.admin-shell h1,.admin-shell h2{font-family:var(--serif)}.admin-head{border-top:3px solid var(--ink);border-bottom:1px solid var(--rule)}.admin-head .admin-shell{display:flex;align-items:baseline;gap:20px}.admin-head h1{margin:0}.admin-head a{margin-left:auto}.eyebrow{color:var(--accent);font-size:12px;font-weight:700;letter-spacing:.16em}.config-error{border-top:3px solid var(--accent);padding-top:12px}.section-head{border-top:3px solid var(--ink);display:flex;align-items:baseline;gap:14px;margin-top:28px}.section-head h2{font-size:24px}.section-head span{color:var(--soft);font-size:13px}.key-row,.log-row{border-bottom:1px solid var(--rule);padding:14px 0}.key-row>div{display:flex;gap:14px;align-items:baseline}.key-row time,.status{color:var(--soft);font-size:13px}.key-row code{display:block;overflow-wrap:anywhere;margin-top:6px}.key-row.revoked{color:var(--soft);text-decoration-color:var(--soft)}.key-row.revoked code{text-decoration:line-through}.log-row{display:grid;grid-template-columns:180px 180px 140px 1fr;gap:16px}.log-row p{grid-column:1/-1;margin:0;color:var(--soft)}@media(max-width:767px){.log-row{grid-template-columns:1fr;gap:3px}.admin-head .admin-shell{display:block}}
</style>`
