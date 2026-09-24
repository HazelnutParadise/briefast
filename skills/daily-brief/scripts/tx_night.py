#!/usr/bin/env python3
"""Read saved TAIFEX TX CSV files for one pre-market overnight signal."""

import argparse
import csv
import io
import json
import math
import re
import sys
from datetime import date
from decimal import Decimal, InvalidOperation
from pathlib import Path


REQUIRED_COLUMNS = ('交易日期', '契約', '到期月份(週別)', '收盤價', '成交量', '交易時段')


def parse_date(value):
    try:
        parsed = date.fromisoformat(value)
    except ValueError as exc:
        raise ValueError(f'日期不是有效的 YYYY-MM-DD：{value}') from exc
    if parsed.isoformat() != value:
        raise ValueError(f'日期不是有效的 YYYY-MM-DD：{value}')
    return parsed


def read_csv(path):
    try:
        content = Path(path).read_bytes().decode('cp950')
    except (OSError, UnicodeError) as exc:
        raise ValueError(f'無法讀取期交所 CSV：{path}（{exc}）') from exc
    reader = csv.reader(io.StringIO(content, newline=''), strict=True)
    try:
        header = [field.strip().lstrip('\ufeff') for field in next(reader)]
    except StopIteration as exc:
        raise ValueError(f'期交所 CSV 為空：{path}') from exc
    except csv.Error as exc:
        raise ValueError(f'期交所 CSV 格式錯誤：{path}（{exc}）') from exc
    for name in REQUIRED_COLUMNS:
        if header.count(name) != 1:
            raise ValueError(f'期交所 CSV 欄位錯誤：{path} 缺少或重複 {name}')
    indexes = {name: header.index(name) for name in REQUIRED_COLUMNS}
    rows = []
    try:
        for line, fields in enumerate(reader, start=2):
            if not fields or all(not field.strip() for field in fields):
                continue
            if len(fields) <= max(indexes.values()):
                raise ValueError(f'期交所 CSV 資料欄位不足：{path}:{line}')
            rows.append({name: fields[indexes[name]].strip() for name in REQUIRED_COLUMNS})
    except csv.Error as exc:
        raise ValueError(f'期交所 CSV 格式錯誤：{path}（{exc}）') from exc
    return rows


def positive_price(value, context):
    try:
        result = Decimal(value.replace(',', ''))
    except InvalidOperation as exc:
        raise ValueError(f'{context}收盤價無效：{value}') from exc
    if not result.is_finite() or result <= 0:
        raise ValueError(f'{context}收盤價無效：{value}')
    return result


def volume(value):
    try:
        result = int(value.replace(',', ''))
    except ValueError as exc:
        raise ValueError(f'夜盤成交量無效：{value}') from exc
    if result < 0:
        raise ValueError(f'夜盤成交量無效：{value}')
    return result


def json_number(value):
    converted = float(value)
    if not math.isfinite(converted):
        raise ValueError('期交所數值超出有效範圍')
    if value == value.to_integral_value() and abs(value) <= 9007199254740991:
        return int(value)
    return converted


def build_signal(report_date, previous_date, report_csv, previous_csv):
    report_day = parse_date(report_date)
    previous_day = parse_date(previous_date)
    if previous_day >= report_day:
        raise ValueError('前一交易日必須早於報告日')
    report_rows = read_csv(report_csv)
    previous_rows = read_csv(previous_csv)
    night_rows = []
    for row in report_rows:
        if (row['交易日期'] != report_day.strftime('%Y/%m/%d') or row['契約'] != 'TX'
                or row['交易時段'] != '盤後'):
            continue
        contract = row['到期月份(週別)']
        if not re.fullmatch(r'\d{6}', contract):
            continue
        traded = volume(row['成交量'])
        if traded == 0:
            continue
        night_rows.append((traded, contract, positive_price(row['收盤價'], '夜盤')))
    if not night_rows:
        raise ValueError(f'{report_date} 找不到已完成且有成交的 TX 盤後單月契約')
    night_rows.sort(reverse=True)
    if len(night_rows) > 1 and night_rows[0][0] == night_rows[1][0]:
        raise ValueError(f'{report_date} TX 夜盤最高成交量有多筆，無法選定契約')
    night_volume, contract, night_close = night_rows[0]
    regular_rows = [row for row in previous_rows
                    if row['交易日期'] == previous_day.strftime('%Y/%m/%d')
                    and row['契約'] == 'TX'
                    and row['交易時段'] == '一般'
                    and row['到期月份(週別)'] == contract]
    if len(regular_rows) != 1:
        raise ValueError(f'{previous_date} 找不到唯一的 TX {contract} 一般時段收盤')
    regular_close = positive_price(regular_rows[0]['收盤價'], '前一交易日日盤')
    change = night_close - regular_close
    return {
        'report_date': report_date,
        'previous_regular_date': previous_date,
        'contract': contract,
        'night_close': json_number(night_close),
        'previous_regular_close': json_number(regular_close),
        'change_points': json_number(change),
        'change_pct': json_number(round(change / regular_close * 100, 4)),
        'direction': 'up' if change > 0 else 'down' if change < 0 else 'flat',
        'night_volume': night_volume,
    }


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument('--report-date', required=True)
    parser.add_argument('--previous-date', required=True)
    parser.add_argument('--report-csv', required=True)
    parser.add_argument('--previous-csv', required=True)
    args = parser.parse_args()
    try:
        result = build_signal(args.report_date, args.previous_date,
                              args.report_csv, args.previous_csv)
    except ValueError as exc:
        print(f'tx_night: {exc}', file=sys.stderr)
        return 1
    print(json.dumps(result, ensure_ascii=False, separators=(',', ':')))
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
