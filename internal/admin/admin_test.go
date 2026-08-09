package admin_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/HazelnutParadise/briefast/internal/admin"
	briefapi "github.com/HazelnutParadise/briefast/internal/api"
	"github.com/HazelnutParadise/briefast/internal/store"
	sy "github.com/HazelnutParadise/syralit"
)

func setupAdmin(t *testing.T, password string) (*store.Store, *admin.App) {
	t.Helper()
	s, err := store.Open(filepath.Join(t.TempDir(), "briefast.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = s.Close() })
	a := admin.New(s,
		admin.WithPasswordSource(func() string { return password }),
		admin.WithTokenGenerator(func() (string, error) { return "brf_fixed_complete_token", nil }),
	)
	return s, a
}

func allHTML(at *sy.AppTest) string {
	var b strings.Builder
	for _, node := range at.FindAll("html") {
		value, _ := node.Props["html"].(string)
		b.WriteString(value)
	}
	return b.String()
}

func login(t *testing.T, at *sy.AppTest, password string) {
	t.Helper()
	at.Run()
	at.SetValue("admin-password", password).Click("admin-login").Run()
	at.Run()
}

func TestAdminLoginGateAndMissingConfiguration(t *testing.T) {
	_, missing := setupAdmin(t, "")
	at := sy.NewAppTest(missing.Page)
	at.Run()
	if got := allHTML(at); !strings.Contains(got, "尚未設定管理員密碼") || len(at.FindAll("text_input")) != 0 || len(at.FindAll("button")) != 0 {
		t.Fatalf("missing-password view rendered login or lacks error: %s", got)
	}

	_, app := setupAdmin(t, "correct-password")
	at = sy.NewAppTest(app.Page)
	at.Run()
	got := allHTML(at)
	if !strings.Contains(got, "後台登入") || strings.Contains(got, "API key 管理") || strings.Contains(got, "報告更新紀錄") {
		t.Fatalf("unauthenticated view leaked admin data: %s", got)
	}
	at.SetValue("admin-password", "wrong").Click("admin-login").Run()
	if len(at.FindAll("status")) == 0 {
		t.Fatal("wrong password did not render an error")
	}
	login(t, at, "correct-password")
	if got := allHTML(at); !strings.Contains(got, "API key 管理") || !strings.Contains(got, "報告更新紀錄") {
		t.Fatalf("authenticated dashboard missing: %s", got)
	}
}

func TestCreateKeyAndViewCompleteTokenOnLaterLogin(t *testing.T) {
	_, app := setupAdmin(t, "secret")
	at := sy.NewAppTest(app.Page)
	login(t, at, "secret")
	at.SetValue("key-name", "cowork-daily").Click("create-key").Run()
	got := allHTML(at)
	if !strings.Contains(got, "cowork-daily") || !strings.Contains(got, "brf_fixed_complete_token") {
		t.Fatalf("created key not fully visible: %s", got)
	}

	later := sy.NewAppTest(app.Page)
	login(t, later, "secret")
	got = allHTML(later)
	if !strings.Contains(got, "cowork-daily") || !strings.Contains(got, "brf_fixed_complete_token") {
		t.Fatalf("stored key not visible on later login: %s", got)
	}
}

func TestRevokeKeyChangesAdminStateAndAPIAuthentication(t *testing.T) {
	s, app := setupAdmin(t, "secret")
	key, err := s.CreateAPIKey(context.Background(), "revoke-me", "brf_revoke_me")
	if err != nil {
		t.Fatal(err)
	}
	at := sy.NewAppTest(app.Page)
	login(t, at, "secret")
	at.Click("revoke-key-" + strconv.FormatInt(key.ID, 10)).Run()
	at.Run()
	if got := allHTML(at); !strings.Contains(got, `class="key-row revoked"`) || !strings.Contains(got, "已撤銷") {
		t.Fatalf("revoked key not visually distinguished: %s", got)
	}

	handler := briefapi.NewReportHandler(s, nil)
	req := httptest.NewRequest(http.MethodPost, "/api/report", bytes.NewBufferString(`{}`))
	req.Header.Set("Authorization", "Bearer "+key.Token)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("revoked key status = %d, body = %s", w.Code, w.Body.String())
	}
}

func TestUpdateLogsRenderNewestFirstIncludingRejections(t *testing.T) {
	s, app := setupAdmin(t, "secret")
	older, _ := time.Parse(time.RFC3339, "2026-08-06T07:50:00+08:00")
	newer, _ := time.Parse(time.RFC3339, "2026-08-07T07:50:00+08:00")
	if err := s.AddUpdateLog(context.Background(), store.UpdateLog{At: older, KeyName: "first-key", ReportDate: "2026-08-06", Action: "ingest_ok"}); err != nil {
		t.Fatal(err)
	}
	if err := s.AddUpdateLog(context.Background(), store.UpdateLog{At: newer, KeyName: "second-key", ReportDate: "2026-08-07", Action: "ingest_rejected_schema", Detail: "headline 不得為空"}); err != nil {
		t.Fatal(err)
	}
	at := sy.NewAppTest(app.Page)
	login(t, at, "secret")
	got := allHTML(at)
	if strings.Index(got, "second-key") > strings.Index(got, "first-key") || !strings.Contains(got, "ingest_rejected_schema") || !strings.Contains(got, "headline 不得為空") {
		t.Fatalf("update log ordering or fields incorrect: %s", got)
	}
}
