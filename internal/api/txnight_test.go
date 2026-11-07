package api

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/text/encoding/traditionalchinese"
)

const txHeader = "交易日期,契約,到期月份(週別),開盤價,最高價,最低價,收盤價,漲跌價,漲跌%,成交量,結算價,未沖銷契約數,最後最佳買價,最後最佳賣價,歷史最高價,歷史最低價,是否因訊息面暫停交易,交易時段,價差對單式委託成交量"

func txCSV(t *testing.T, rows ...string) []byte {
	t.Helper()
	out, err := traditionalchinese.Big5.NewEncoder().String(txHeader + "\n" + strings.Join(rows, "\n") + "\n")
	if err != nil {
		t.Fatal(err)
	}
	return []byte(out)
}

func txUpload(t *testing.T, handler http.Handler, token string, fields map[string]string, files map[string][]byte) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	for k, v := range fields {
		if err := writer.WriteField(k, v); err != nil {
			t.Fatal(err)
		}
	}
	for k, v := range files {
		part, err := writer.CreateFormFile(k, k+".csv")
		if err != nil {
			t.Fatal(err)
		}
		_, _ = part.Write(v)
	}
	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/tx-night", &body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)
	return w
}

func txFiles(t *testing.T) map[string][]byte {
	return map[string][]byte{
		"report_csv":   txCSV(t, "2026/09/24,TX,202610  ,-,-,-,47909,-,-,28947,-,-,-,-,-,-,,盤後,"),
		"previous_csv": txCSV(t, "2026/09/23,TX,202610  ,-,-,-,48336,-,-,32929,-,-,-,-,-,-,,一般,"),
	}
}

var txDates = map[string]string{"report_date": "2026-09-24", "previous_date": "2026-09-23"}

func TestTXNightReturnsSignal(t *testing.T) {
	s, key, _, _ := setupHandler(t)
	w := txUpload(t, TXNightHandler(s), key.Token, txDates, txFiles(t))
	if w.Code != http.StatusOK {
		t.Fatalf("狀態碼 = %d, body=%s", w.Code, w.Body.String())
	}
	body := assertPointer(t, w)
	if body["ok"] != true || body["contract"] != "202610" || body["change_points"] != float64(-427) ||
		body["direction"] != "down" || body["night_volume"] != float64(28947) || body["change_pct"] != -0.8834 {
		t.Fatalf("回應 = %v", body)
	}
	if n := countLogsWithAction(t, s, "txnight_rejected_validation"); n != 0 {
		t.Fatalf("成功不應寫 log, got %d", n)
	}
}

func TestTXNightRejectsBadInputWithoutPartialResult(t *testing.T) {
	s, key, _, _ := setupHandler(t)
	stale := txFiles(t)
	stale["report_csv"] = txCSV(t, "2026/09/23,TX,202610,-,-,-,47909,-,-,28947,-,-,-,-,-,-,,盤後,")
	missing := txFiles(t)
	delete(missing, "previous_csv")
	cases := []struct {
		name   string
		fields map[string]string
		files  map[string][]byte
		want   string
	}{
		{"舊日夜盤", txDates, stale, "找不到已完成"},
		{"缺檔案", txDates, missing, "previous_csv"},
		{"缺日期", map[string]string{"report_date": "2026-09-24"}, txFiles(t), "previous_date"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			w := txUpload(t, TXNightHandler(s), key.Token, c.fields, c.files)
			if w.Code != http.StatusBadRequest || !strings.Contains(w.Body.String(), c.want) {
				t.Fatalf("狀態碼 = %d, body=%s", w.Code, w.Body.String())
			}
			body := assertPointer(t, w)
			if _, partial := body["contract"]; partial || body["ok"] != false {
				t.Fatalf("不應有部分結果: %v", body)
			}
		})
	}
	if n := countLogsWithAction(t, s, "txnight_rejected_validation"); n != len(cases) {
		t.Fatalf("txnight_rejected_validation = %d, want %d", n, len(cases))
	}
}

func TestTXNightRejectsMissingKey(t *testing.T) {
	s, _, _, _ := setupHandler(t)
	w := txUpload(t, TXNightHandler(s), "", txDates, txFiles(t))
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("狀態碼 = %d", w.Code)
	}
	if n := countLogsWithAction(t, s, "txnight_rejected_auth"); n != 1 {
		t.Fatalf("txnight_rejected_auth = %d", n)
	}
}

func TestTXNightRejectsNonMultipart(t *testing.T) {
	s, key, _, _ := setupHandler(t)
	req := httptest.NewRequest(http.MethodPost, "/api/tx-night", strings.NewReader("{}"))
	req.Header.Set("Authorization", "Bearer "+key.Token)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	TXNightHandler(s).ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("狀態碼 = %d, body=%s", w.Code, w.Body.String())
	}
	assertPointer(t, w)
}
