package txnight

import (
	"strings"
	"testing"

	"golang.org/x/text/encoding/traditionalchinese"
)

const header = "交易日期,契約,到期月份(週別),開盤價,最高價,最低價,收盤價," +
	"漲跌價,漲跌%,成交量,結算價,未沖銷契約數,最後最佳買價,最後最佳賣價," +
	"歷史最高價,歷史最低價,是否因訊息面暫停交易,交易時段,價差對單式委託成交量"

func row(day, contract, close, traded, session string) string {
	return day + ",TX," + contract + ",-,-,-," + close + ",-,-," + traded + ",-,-,-,-,-,-,," + session + ","
}

func cp950(t *testing.T, lines ...string) []byte {
	t.Helper()
	encoded, err := traditionalchinese.Big5.NewEncoder().String(strings.Join(lines, "\n") + "\n")
	if err != nil {
		t.Fatal(err)
	}
	return []byte(encoded)
}

func defaultFiles(t *testing.T) ([]byte, []byte) {
	return cp950(t, header,
			row("2026/09/24", "202610  ", "47909", "28947", "盤後"),
			row("2026/09/24", "202609  ", "47940", "20", "盤後"),
			row("2026/09/24", "202610  ", "48123", "37196", "一般")),
		cp950(t, header,
			row("2026/09/23", "202610  ", "48336", "32929", "一般"),
			row("2026/09/23", "202609  ", "48000", "5", "一般"))
}

func TestNightUsesPreviousSameContractRegularClose(t *testing.T) {
	report, previous := defaultFiles(t)
	got, err := Signal("2026-09-24", "2026-09-23", report, previous)
	if err != nil {
		t.Fatal(err)
	}
	want := Result{ReportDate: "2026-09-24", PreviousRegularDate: "2026-09-23", Contract: "202610",
		NightClose: 47909, PreviousRegularClose: 48336, ChangePoints: -427, ChangePct: -0.8834,
		Direction: "down", NightVolume: 28947}
	if got != want {
		t.Fatalf("Signal = %+v\nwant %+v", got, want)
	}
}

func TestWeekendAndRolloverUseSuppliedDateAndSameContract(t *testing.T) {
	report := cp950(t, header,
		row("2026/09/21", "202610  ", "47405", "21514", "盤後"),
		row("2026/09/21", "202609  ", "47500", "100", "盤後"))
	previous := cp950(t, header,
		row("2026/09/18", "202610  ", "47418", "44655", "一般"),
		row("2026/09/18", "202609  ", "47000", "30", "一般"))
	got, err := Signal("2026-09-21", "2026-09-18", report, previous)
	if err != nil {
		t.Fatal(err)
	}
	if got.Contract != "202610" || got.ChangePoints != -13 || got.Direction != "down" {
		t.Fatalf("Signal = %+v", got)
	}
}

func TestFlatMoveHasNoForcedDirection(t *testing.T) {
	got, err := Signal("2026-09-24", "2026-09-23",
		cp950(t, header, row("2026/09/24", "202610", "48336", "28947", "盤後")),
		cp950(t, header, row("2026/09/23", "202610", "48336", "32929", "一般")))
	if err != nil {
		t.Fatal(err)
	}
	if got.ChangePoints != 0 || got.ChangePct != 0 || got.Direction != "flat" {
		t.Fatalf("Signal = %+v", got)
	}
}

func TestFractionalPricesKeepDecimals(t *testing.T) {
	got, err := Signal("2026-09-24", "2026-09-23",
		cp950(t, header, row("2026/09/24", "202610", "48336.5", "10", "盤後")),
		cp950(t, header, row("2026/09/23", "202610", "48336", "32929", "一般")))
	if err != nil {
		t.Fatal(err)
	}
	if got.NightClose != 48336.5 || got.ChangePoints != 0.5 || got.Direction != "up" {
		t.Fatalf("Signal = %+v", got)
	}
}

func TestUnusableInputsFail(t *testing.T) {
	regular := cp950(t, header, row("2026/09/23", "202610", "48336", "32929", "一般"))
	cases := []struct {
		name, report, previous string
		reportCSV, previousCSV []byte
		want                   string
	}{
		{"只有日盤", "2026-09-24", "2026-09-23", cp950(t, header, row("2026/09/24", "202610", "48123", "37196", "一般")), regular, "找不到已完成"},
		{"舊日夜盤", "2026-09-24", "2026-09-23", cp950(t, header, row("2026/09/23", "202610", "47909", "28947", "盤後")), regular, "找不到已完成"},
		{"成交量並列", "2026-09-24", "2026-09-23", cp950(t, header,
			row("2026/09/24", "202610", "47909", "28947", "盤後"),
			row("2026/09/24", "202611", "47900", "28947", "盤後")), regular, "最高成交量有多筆"},
		{"缺同契約日盤", "2026-09-24", "2026-09-23", cp950(t, header, row("2026/09/24", "202610", "47909", "28947", "盤後")),
			cp950(t, header, row("2026/09/23", "202609", "48000", "32929", "一般")), "找不到唯一"},
		{"前一日不早於報告日", "2026-09-24", "2026-09-24", nil, nil, "前一交易日必須早於"},
		{"無效日期", "2026-09-31", "2026-09-23", nil, nil, "日期不是有效"},
		{"收盤價無效", "2026-09-24", "2026-09-23", cp950(t, header, row("2026/09/24", "202610", "NaN", "28947", "盤後")), regular, "收盤價無效"},
		{"收盤價非正數", "2026-09-24", "2026-09-23", cp950(t, header, row("2026/09/24", "202610", "0", "28947", "盤後")), regular, "收盤價無效"},
		{"成交量無效", "2026-09-24", "2026-09-23", cp950(t, header, row("2026/09/24", "202610", "47909", "bad", "盤後")), regular, "成交量無效"},
		{"欄位錯誤", "2026-09-24", "2026-09-23", cp950(t, strings.Replace(header, "交易時段", "錯誤欄位", 1),
			row("2026/09/24", "202610", "47909", "28947", "盤後")), regular, "欄位錯誤"},
		{"引號未關", "2026-09-24", "2026-09-23", cp950(t, header, `"unterminated`), regular, "CSV 格式錯誤"},
		{"空檔", "2026-09-24", "2026-09-23", []byte{}, regular, "CSV 為空"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, err := Signal(c.report, c.previous, c.reportCSV, c.previousCSV)
			if err == nil || !strings.Contains(err.Error(), c.want) {
				t.Fatalf("err = %v, want containing %q", err, c.want)
			}
		})
	}
}
