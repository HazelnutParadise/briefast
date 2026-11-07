// Package txnight reads two TAIFEX TX daily CSV downloads and derives the
// pre-market overnight signal: the report date's night-session close against
// the preceding trading day's regular close for the same contract.
package txnight

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/encoding/traditionalchinese"
)

var requiredColumns = []string{"交易日期", "契約", "到期月份(週別)", "收盤價", "成交量", "交易時段"}

var monthContract = regexp.MustCompile(`^\d{6}$`)

// Result is the verified signal returned by POST /api/tx-night.
type Result struct {
	ReportDate           string  `json:"report_date"`
	PreviousRegularDate  string  `json:"previous_regular_date"`
	Contract             string  `json:"contract"`
	NightClose           float64 `json:"night_close"`
	PreviousRegularClose float64 `json:"previous_regular_close"`
	ChangePoints         float64 `json:"change_points"`
	ChangePct            float64 `json:"change_pct"`
	Direction            string  `json:"direction"`
	NightVolume          int64   `json:"night_volume"`
}

type csvRow map[string]string

type nightRow struct {
	volume   int64
	contract string
	close    *big.Rat
}

// Signal validates both files and returns the signal, or an error describing
// the first problem. It never returns a partial result.
func Signal(reportDate, previousDate string, reportCSV, previousCSV []byte) (Result, error) {
	reportDay, err := parseDate(reportDate)
	if err != nil {
		return Result{}, err
	}
	previousDay, err := parseDate(previousDate)
	if err != nil {
		return Result{}, err
	}
	if !previousDay.Before(reportDay) {
		return Result{}, errors.New("前一交易日必須早於報告日")
	}
	reportRows, err := readCSV("report_csv", reportCSV)
	if err != nil {
		return Result{}, err
	}
	previousRows, err := readCSV("previous_csv", previousCSV)
	if err != nil {
		return Result{}, err
	}

	var nights []nightRow
	for _, row := range reportRows {
		if row["交易日期"] != reportDay.Format("2006/01/02") || row["契約"] != "TX" || row["交易時段"] != "盤後" {
			continue
		}
		contract := row["到期月份(週別)"]
		if !monthContract.MatchString(contract) {
			continue
		}
		traded, err := volume(row["成交量"])
		if err != nil {
			return Result{}, err
		}
		if traded == 0 {
			continue
		}
		price, err := positivePrice(row["收盤價"], "夜盤")
		if err != nil {
			return Result{}, err
		}
		nights = append(nights, nightRow{volume: traded, contract: contract, close: price})
	}
	if len(nights) == 0 {
		return Result{}, fmt.Errorf("%s 找不到已完成且有成交的 TX 盤後單月契約", reportDate)
	}
	sort.SliceStable(nights, func(i, j int) bool { return nights[i].volume > nights[j].volume })
	if len(nights) > 1 && nights[0].volume == nights[1].volume {
		return Result{}, fmt.Errorf("%s TX 夜盤最高成交量有多筆，無法選定契約", reportDate)
	}
	night := nights[0]

	var regular []csvRow
	for _, row := range previousRows {
		if row["交易日期"] == previousDay.Format("2006/01/02") && row["契約"] == "TX" &&
			row["交易時段"] == "一般" && row["到期月份(週別)"] == night.contract {
			regular = append(regular, row)
		}
	}
	if len(regular) != 1 {
		return Result{}, fmt.Errorf("%s 找不到唯一的 TX %s 一般時段收盤", previousDate, night.contract)
	}
	regularClose, err := positivePrice(regular[0]["收盤價"], "前一交易日日盤")
	if err != nil {
		return Result{}, err
	}

	change := new(big.Rat).Sub(night.close, regularClose)
	pct := new(big.Rat).Mul(new(big.Rat).Quo(change, regularClose), big.NewRat(100, 1))
	direction := "flat"
	switch change.Sign() {
	case 1:
		direction = "up"
	case -1:
		direction = "down"
	}
	return Result{
		ReportDate:           reportDate,
		PreviousRegularDate:  previousDate,
		Contract:             night.contract,
		NightClose:           toFloat(night.close),
		PreviousRegularClose: toFloat(regularClose),
		ChangePoints:         toFloat(change),
		ChangePct:            roundTo(pct, 4),
		Direction:            direction,
		NightVolume:          night.volume,
	}, nil
}

func parseDate(value string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil || parsed.Format("2006-01-02") != value {
		return time.Time{}, fmt.Errorf("日期不是有效的 YYYY-MM-DD：%s", value)
	}
	return parsed, nil
}

func readCSV(name string, raw []byte) ([]csvRow, error) {
	raw = bytes.TrimPrefix(raw, []byte("\xef\xbb\xbf"))
	decoded, err := traditionalchinese.Big5.NewDecoder().Bytes(raw)
	if err != nil {
		return nil, fmt.Errorf("無法以 CP950 讀取期交所 CSV：%s（%v）", name, err)
	}
	reader := csv.NewReader(bytes.NewReader(decoded))
	reader.FieldsPerRecord = -1
	header, err := reader.Read()
	if errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("期交所 CSV 為空：%s", name)
	}
	if err != nil {
		return nil, fmt.Errorf("期交所 CSV 格式錯誤：%s（%v）", name, err)
	}
	for i := range header {
		header[i] = strings.TrimPrefix(strings.TrimSpace(header[i]), "\ufeff")
	}
	indexes := map[string]int{}
	maxIndex := 0
	for _, column := range requiredColumns {
		count, at := 0, -1
		for i, field := range header {
			if field == column {
				count++
				at = i
			}
		}
		if count != 1 {
			return nil, fmt.Errorf("期交所 CSV 欄位錯誤：%s 缺少或重複 %s", name, column)
		}
		indexes[column] = at
		maxIndex = max(maxIndex, at)
	}
	var rows []csvRow
	for {
		fields, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("期交所 CSV 格式錯誤：%s（%v）", name, err)
		}
		if allBlank(fields) {
			continue
		}
		if len(fields) <= maxIndex {
			line, _ := reader.FieldPos(0)
			return nil, fmt.Errorf("期交所 CSV 資料欄位不足：%s 第 %d 行", name, line)
		}
		row := csvRow{}
		for column, i := range indexes {
			row[column] = strings.TrimSpace(fields[i])
		}
		rows = append(rows, row)
	}
	return rows, nil
}

func allBlank(fields []string) bool {
	for _, field := range fields {
		if strings.TrimSpace(field) != "" {
			return false
		}
	}
	return true
}

func positivePrice(value, context string) (*big.Rat, error) {
	price, ok := new(big.Rat).SetString(strings.ReplaceAll(value, ",", ""))
	if !ok || price.Sign() <= 0 {
		return nil, fmt.Errorf("%s收盤價無效：%s", context, value)
	}
	return price, nil
}

func volume(value string) (int64, error) {
	traded, err := strconv.ParseInt(strings.ReplaceAll(value, ",", ""), 10, 64)
	if err != nil || traded < 0 {
		return 0, fmt.Errorf("夜盤成交量無效：%s", value)
	}
	return traded, nil
}

func toFloat(value *big.Rat) float64 {
	f, _ := value.Float64()
	return f
}

func roundTo(value *big.Rat, places int) float64 {
	f, _ := strconv.ParseFloat(value.FloatString(places), 64)
	return f
}
