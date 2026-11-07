package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/HazelnutParadise/briefast/internal/store"
	"github.com/HazelnutParadise/briefast/internal/txnight"
)

const maxTXNightBytes = 4 << 20

// TXNightHandler serves POST /api/tx-night: the agent uploads the two TAIFEX
// CSV downloads and the server returns the verified overnight signal.
func TXNightHandler(s *store.Store) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !requireMethod(w, r, http.MethodPost) {
			return
		}
		key, ok := authorizeRequest(s, w, r, "txnight_rejected_auth", "")
		if !ok {
			return
		}
		reject := func(errs ...string) {
			rejectRequest(s, w, r, key, "txnight_rejected_validation", r.FormValue("report_date"), errs)
		}
		r.Body = http.MaxBytesReader(w, r.Body, maxTXNightBytes)
		if err := r.ParseMultipartForm(maxTXNightBytes); err != nil {
			reject(fmt.Sprintf("request 必須是 4 MiB 以內的 multipart/form-data：%v", err))
			return
		}
		var errs []string
		files := map[string][]byte{}
		for _, field := range []string{"report_csv", "previous_csv"} {
			content, err := formFile(r, field)
			if err != nil {
				errs = append(errs, err.Error())
				continue
			}
			files[field] = content
		}
		for _, field := range []string{"report_date", "previous_date"} {
			if r.FormValue(field) == "" {
				errs = append(errs, field+" 不得為空")
			}
		}
		if len(errs) != 0 {
			reject(errs...)
			return
		}
		result, err := txnight.Signal(r.FormValue("report_date"), r.FormValue("previous_date"), files["report_csv"], files["previous_csv"])
		if err != nil {
			reject(err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"ok":                     true,
			"report_date":            result.ReportDate,
			"previous_regular_date":  result.PreviousRegularDate,
			"contract":               result.Contract,
			"night_close":            result.NightClose,
			"previous_regular_close": result.PreviousRegularClose,
			"change_points":          result.ChangePoints,
			"change_pct":             result.ChangePct,
			"direction":              result.Direction,
			"night_volume":           result.NightVolume,
		})
	})
}

func formFile(r *http.Request, field string) ([]byte, error) {
	file, _, err := r.FormFile(field)
	if errors.Is(err, http.ErrMissingFile) {
		return nil, fmt.Errorf("缺少檔案欄位 %s", field)
	}
	if err != nil {
		return nil, fmt.Errorf("無法讀取 %s：%v", field, err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("無法讀取 %s：%v", field, err)
	}
	return content, nil
}
