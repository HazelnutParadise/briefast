#!/usr/bin/env python3
"""Offline checks for the TAIFEX night-session parser."""

import contextlib
import io
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

import tx_night


HEADER = ('交易日期,契約,到期月份(週別),開盤價,最高價,最低價,收盤價,'
          '漲跌價,漲跌%,成交量,結算價,未沖銷契約數,最後最佳買價,最後最佳賣價,'
          '歷史最高價,歷史最低價,是否因訊息面暫停交易,交易時段,價差對單式委託成交量')


def row(day, contract, close, traded, session, product='TX'):
    return f'{day},{product},{contract},-,-,-,{close},-,-,{traded},-,-,-,-,-,-,,{session},'


class TXNightTest(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.report = Path(self.temp.name) / 'report.csv'
        self.previous = Path(self.temp.name) / 'previous.csv'
        self.save(
            [row('2026/09/24', '202610  ', '47909', '28947', '盤後'),
             row('2026/09/24', '202609  ', '47940', '20', '盤後'),
             row('2026/09/24', '202610  ', '48123', '37196', '一般')],
            [row('2026/09/23', '202610  ', '48336', '32929', '一般'),
             row('2026/09/23', '202609  ', '48000', '5', '一般')],
        )

    def save(self, report_rows, previous_rows, header=HEADER):
        self.report.write_bytes(('\n'.join([header] + report_rows) + '\n').encode('cp950'))
        self.previous.write_bytes(('\n'.join([HEADER] + previous_rows) + '\n').encode('cp950'))

    def signal(self, report_date='2026-09-24', previous_date='2026-09-23'):
        return tx_night.build_signal(report_date, previous_date, self.report, self.previous)

    def test_report_date_night_uses_previous_same_contract_regular_close(self):
        result = self.signal()
        self.assertEqual(result['contract'], '202610')
        self.assertEqual(result['night_close'], 47909)
        self.assertEqual(result['previous_regular_close'], 48336)
        self.assertEqual(result['change_points'], -427)
        self.assertEqual(result['change_pct'], -0.8834)
        self.assertEqual(result['direction'], 'down')
        self.assertEqual(result['night_volume'], 28947)

    def test_weekend_and_rollover_use_supplied_trading_date_and_same_contract(self):
        self.save(
            [row('2026/09/21', '202610  ', '47405', '21514', '盤後'),
             row('2026/09/21', '202609  ', '47500', '100', '盤後')],
            [row('2026/09/18', '202610  ', '47418', '44655', '一般'),
             row('2026/09/18', '202609  ', '47000', '30', '一般')],
        )
        result = self.signal('2026-09-21', '2026-09-18')
        self.assertEqual((result['contract'], result['change_points'], result['direction']),
                         ('202610', -13, 'down'))

    def test_flat_move_has_no_forced_direction(self):
        self.save(
            [row('2026/09/24', '202610  ', '48336', '28947', '盤後')],
            [row('2026/09/23', '202610  ', '48336', '32929', '一般')],
        )
        result = self.signal()
        self.assertEqual((result['change_points'], result['change_pct'], result['direction']),
                         (0, 0, 'flat'))

    def test_missing_night_and_stale_date_fail(self):
        for report_rows in (
            [row('2026/09/24', '202610  ', '48123', '37196', '一般')],
            [row('2026/09/23', '202610  ', '47909', '28947', '盤後')],
        ):
            with self.subTest(report_rows=report_rows):
                self.save(report_rows, [row('2026/09/23', '202610  ', '48336', '32929', '一般')])
                with self.assertRaisesRegex(ValueError, '找不到已完成'):
                    self.signal()

    def test_tied_maximum_volume_fails(self):
        self.save(
            [row('2026/09/24', '202610  ', '47909', '28947', '盤後'),
             row('2026/09/24', '202611  ', '47900', '28947', '盤後')],
            [row('2026/09/23', '202610  ', '48336', '32929', '一般')],
        )
        with self.assertRaisesRegex(ValueError, '最高成交量有多筆'):
            self.signal()

    def test_missing_same_contract_regular_close_fails(self):
        self.save(
            [row('2026/09/24', '202610  ', '47909', '28947', '盤後')],
            [row('2026/09/23', '202609  ', '48000', '32929', '一般')],
        )
        with self.assertRaisesRegex(ValueError, '找不到唯一'):
            self.signal()

    def test_bad_dates_prices_volumes_and_header_fail(self):
        with self.assertRaisesRegex(ValueError, '前一交易日必須早於'):
            self.signal(previous_date='2026-09-24')
        with self.assertRaisesRegex(ValueError, '日期不是有效'):
            self.signal(report_date='2026-09-31')
        self.save([row('2026/09/24', '202610  ', 'NaN', '28947', '盤後')],
                  [row('2026/09/23', '202610  ', '48336', '32929', '一般')])
        with self.assertRaisesRegex(ValueError, '收盤價無效'):
            self.signal()
        self.save([row('2026/09/24', '202610  ', '47909', 'bad', '盤後')],
                  [row('2026/09/23', '202610  ', '48336', '32929', '一般')])
        with self.assertRaisesRegex(ValueError, '成交量無效'):
            self.signal()
        self.save([row('2026/09/24', '202610  ', '47909', '28947', '盤後')],
                  [row('2026/09/23', '202610  ', '48336', '32929', '一般')],
                  header=HEADER.replace('交易時段', '錯誤欄位'))
        with self.assertRaisesRegex(ValueError, '欄位錯誤'):
            self.signal()
        self.report.write_bytes((HEADER + '\n"unterminated\n').encode('cp950'))
        with self.assertRaisesRegex(ValueError, 'CSV 格式錯誤'):
            self.signal()

    def test_cli_failure_has_no_partial_stdout(self):
        self.save([row('2026/09/23', '202610  ', '47909', '28947', '盤後')],
                  [row('2026/09/23', '202610  ', '48336', '32929', '一般')])
        args = ['tx_night.py', '--report-date', '2026-09-24', '--previous-date', '2026-09-23',
                '--report-csv', str(self.report), '--previous-csv', str(self.previous)]
        stdout, stderr = io.StringIO(), io.StringIO()
        with patch('sys.argv', args), contextlib.redirect_stdout(stdout), contextlib.redirect_stderr(stderr):
            code = tx_night.main()
        self.assertEqual(code, 1)
        self.assertEqual(stdout.getvalue(), '')
        self.assertIn('找不到已完成', stderr.getvalue())


if __name__ == '__main__':
    unittest.main()
