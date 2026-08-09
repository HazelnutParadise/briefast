#!/usr/bin/env python3
"""Briefast 去重資料庫：每日晨報流程唯一的已讀紀錄讀寫入口。

資料存成 JSONL，路徑由 BRIEFAST_SEEN_PATH 決定，預設是執行目錄下的
seen.jsonl（也就是 Cowork 的工作資料夾，不進 repo）。

每個指令執行前都會先清掉過期紀錄，檔案不會無限膨脹。

指令
  peek     從 stdin 讀候選新聞，只印出沒看過的，抓全文前先過濾
  record   記錄處理過的新聞，含最後決策
  event    查詢或登記事件層級的鍵（同一件事跨天不重複報）
  stats    看目前存了什麼
  prune    只做清理

典型用法
  curl -s "$API" | ./seen.py peek | ./fetch-and-judge
  ./seen.py record < processed.jsonl
  ./seen.py event check '2330|月營收|2026-07'
"""

from __future__ import annotations

import argparse
import hashlib
import json
import os
import re
import sys
import tempfile
import unicodedata
from datetime import datetime, timedelta, timezone
from pathlib import Path
from urllib.parse import urlparse, parse_qsl, urlencode, urlunparse

TAIPEI = timezone(timedelta(hours=8))

# 文章層紀錄留三個月，足夠涵蓋一季內的重複報導。
# 事件層只留一個月：它的用途是把同一件事的多篇報導合成一條個股詳情，以及讓
# 「延續昨日營收利多」這種敘述接得上，不需要更長的記憶。
ARTICLE_TTL_DAYS = 90
EVENT_TTL_DAYS = 30

# 轉址與分享用的追蹤參數，去掉才能讓同一篇文章的不同來路收斂成同一個 URL。
TRACKING_PREFIXES = ("utm_", "pk_", "mtm_")
TRACKING_KEYS = {"fbclid", "gclid", "igshid", "ref", "share", "from", "src", "spm"}

HTML_TAG = re.compile(r"<[^>]+>")
WHITESPACE = re.compile(r"\s+")
# 常見的署名與版權尾巴，留著會讓同一篇文章在不同版位算出不同雜湊。
BOILERPLATE = re.compile(
    r"(記者[一-鿿]{2,4}[／/][^\n]{0,12}報導"
    r"|[一-鿿]{2,6}／(?:綜合|台北|編譯)報導"
    r"|責任編輯[：:][^\n]{0,12}"
    r"|本文(?:不構成|僅供)[^\n]{0,30}"
    r"|版權所有[^\n]{0,30})"
)


def data_path() -> Path:
    return Path(os.environ.get("BRIEFAST_SEEN_PATH", "seen.jsonl")).expanduser()


def now() -> datetime:
    return datetime.now(TAIPEI)


def digest(value: str) -> str:
    return hashlib.sha256(value.encode("utf-8")).hexdigest()[:16]


def canonical_url(value: str) -> str:
    """去掉追蹤參數與片段，讓同一篇文章的不同連結收斂成同一個字串。

    不動主機名：m.cnyes.com 與 news.cnyes.com 是不同來源路徑，交給 source_id
    與內容雜湊去判斷是不是同一篇，比在這裡猜測改寫規則安全。
    """
    parsed = urlparse(value.strip())
    keep = [
        (key, val)
        for key, val in parse_qsl(parsed.query, keep_blank_values=True)
        if key.lower() not in TRACKING_KEYS
        and not key.lower().startswith(TRACKING_PREFIXES)
    ]
    path = parsed.path.rstrip("/") or "/"
    return urlunparse(
        (
            parsed.scheme.lower(),
            parsed.netloc.lower(),
            path,
            "",
            urlencode(keep),
            "",
        )
    )


def normalize_text(value: str) -> str:
    """把文字壓成可比較的形式：去標籤、去署名、全半形統一、摺疊空白。

    直接雜湊原始 HTML 沒有用——廣告版位與瀏覽數每次都不一樣。
    """
    text = HTML_TAG.sub(" ", value or "")
    text = unicodedata.normalize("NFKC", text)
    text = BOILERPLATE.sub(" ", text)
    return WHITESPACE.sub(" ", text).strip()


SKETCH_SIZE = 64
PUNCTUATION = re.compile(r"[\s，。、：；！？（）「」【】《》%,.\-—–~〜·]")


def bigrams(text: str) -> set[str]:
    """中文沒有詞界，字元二連文比斷詞穩定，也不需要外部套件。"""
    compact = PUNCTUATION.sub("", normalize_text(text))
    return {compact[i : i + 2] for i in range(len(compact) - 1)}


def sketch(text: str) -> list[str]:
    """bottom-k MinHash 簽章，用來估計兩篇文章的 Jaccard 相似度。

    實測過 SimHash：中文改寫稿與不相關文章的漢明距離只差 6 bit（25 對 31），
    分不出門檻。改用二連文 Jaccard 後分離度很清楚——照登轉載 0.91、大幅改寫
    0.27、不相關 0.04——所以存簽章而不是存整組二連文，避免 JSONL 膨脹。
    """
    hashes = sorted(
        int(hashlib.blake2b(gram.encode("utf-8"), digest_size=8).hexdigest(), 16)
        for gram in bigrams(text)
    )
    return [f"{value:016x}" for value in hashes[:SKETCH_SIZE]]


def similarity(left: list[str], right: list[str]) -> float:
    """由兩份 bottom-k 簽章估計 Jaccard，估計誤差實測在 0.03 以內。"""
    if not left or not right:
        return 0.0
    a, b = set(left), set(right)
    merged = sorted(a | b)[:SKETCH_SIZE]
    if not merged:
        return 0.0
    return len(set(merged) & a & b) / len(merged)


def source_id_of(record: dict) -> str:
    """來源原生 ID 優先於 URL：網站改版時 URL 會變，文章編號不會。"""
    explicit = (record.get("source_id") or "").strip()
    if explicit:
        return explicit
    url = record.get("url") or ""
    if not url:
        return ""
    host = urlparse(url).netloc.lower().replace("www.", "")
    tail = [seg for seg in urlparse(url).path.split("/") if seg]
    numeric = [seg for seg in tail if seg.isdigit()]
    marker = numeric[-1] if numeric else (tail[-1] if tail else "")
    return f"{host}:{marker}" if marker else ""


def load(path: Path) -> list[dict]:
    if not path.exists():
        return []
    records = []
    for line in path.read_text(encoding="utf-8").splitlines():
        line = line.strip()
        if not line:
            continue
        try:
            records.append(json.loads(line))
        except json.JSONDecodeError:
            # 壞掉的一行不該讓整天的流程停擺，跳過並在 stderr 留痕。
            print(f"seen.py: 略過無法解析的一行: {line[:80]}", file=sys.stderr)
    return records


def save(path: Path, records: list[dict]) -> None:
    """先寫暫存檔再 rename，中途中斷不會留下半個檔案。"""
    path.parent.mkdir(parents=True, exist_ok=True)
    handle = tempfile.NamedTemporaryFile(
        "w", encoding="utf-8", dir=path.parent, delete=False, suffix=".tmp"
    )
    try:
        for record in records:
            handle.write(json.dumps(record, ensure_ascii=False) + "\n")
        handle.flush()
        os.fsync(handle.fileno())
    finally:
        handle.close()
    os.replace(handle.name, path)


def prune(records: list[dict], reference: datetime | None = None) -> tuple[list[dict], int]:
    reference = reference or now()
    kept = []
    for record in records:
        ttl = EVENT_TTL_DAYS if record.get("kind") == "event" else ARTICLE_TTL_DAYS
        stamp = record.get("first_seen_at", "")
        try:
            seen_at = datetime.fromisoformat(stamp)
        except ValueError:
            # 沒有可解析的時間就當作剛看到，下一輪再判斷，不要默默刪掉資料。
            record["first_seen_at"] = reference.isoformat()
            kept.append(record)
            continue
        if reference - seen_at <= timedelta(days=ttl):
            kept.append(record)
    return kept, len(records) - len(kept)


def open_db(quiet: bool = False) -> tuple[Path, list[dict]]:
    """任何指令的共同入口：載入、清理、寫回，再把資料交出去。"""
    path = data_path()
    records = load(path)
    records, removed = prune(records)
    if removed:
        save(path, records)
        if not quiet:
            print(f"seen.py: 清掉 {removed} 筆過期紀錄", file=sys.stderr)
    return path, records


def read_input(raw: str) -> list[dict]:
    """吃兩種輸入：每行一個 JSON 物件，或每行一個網址。"""
    items = []
    for line in raw.splitlines():
        line = line.strip()
        if not line:
            continue
        if line.startswith("{"):
            try:
                items.append(json.loads(line))
            except json.JSONDecodeError:
                print(f"seen.py: 略過無法解析的輸入: {line[:80]}", file=sys.stderr)
        elif line.startswith("http"):
            items.append({"url": line})
        else:
            print(f"seen.py: 略過無法辨識的輸入: {line[:80]}", file=sys.stderr)
    return items


def index_of(records: list[dict]) -> tuple[set[str], set[str], set[str]]:
    ids, urls, titles = set(), set(), set()
    for record in records:
        if record.get("kind") == "event":
            continue
        if record.get("source_id"):
            ids.add(record["source_id"])
        if record.get("url_hash"):
            urls.add(record["url_hash"])
        if record.get("title_hash"):
            titles.add(record["title_hash"])
    return ids, urls, titles


def cmd_peek(args: argparse.Namespace) -> int:
    """抓全文之前先問：這些候選有哪些沒看過。

    只靠清單頁就有的欄位（網址、標題、原生 ID）判斷，不需要先下載內文。
    """
    _, records = open_db(quiet=args.quiet)
    seen_ids, seen_urls, seen_titles = index_of(records)

    raw = sys.stdin.read() if not args.url else "\n".join(args.url)
    items = read_input(raw)

    fresh, known = [], []
    for item in items:
        sid = source_id_of(item)
        url_hash = digest(canonical_url(item["url"])) if item.get("url") else ""
        title_hash = digest(normalize_text(item["title"])) if item.get("title") else ""

        reason = ""
        if sid and sid in seen_ids:
            reason = "source_id"
        elif url_hash and url_hash in seen_urls:
            reason = "url"
        elif title_hash and title_hash in seen_titles:
            reason = "title"

        if reason:
            known.append({**item, "seen_by": reason})
        else:
            fresh.append(item)

    if args.json:
        print(json.dumps({"new": fresh, "seen": known}, ensure_ascii=False, indent=2))
    else:
        for item in fresh:
            print(json.dumps(item, ensure_ascii=False))

    if not args.quiet:
        print(
            f"seen.py: 候選 {len(items)} 筆，未讀 {len(fresh)} 筆，已讀 {len(known)} 筆",
            file=sys.stderr,
        )
    return 0


def cmd_record(args: argparse.Namespace) -> int:
    """記下處理過的新聞，包含判定不報的，否則明天會重新評估同一批雜訊。"""
    path, records = open_db(quiet=args.quiet)
    items = read_input(sys.stdin.read())
    stamp = now().isoformat()
    seen_ids, seen_urls, _ = index_of(records)

    added = 0
    for item in items:
        sid = source_id_of(item)
        url_hash = digest(canonical_url(item["url"])) if item.get("url") else ""
        if (sid and sid in seen_ids) or (url_hash and url_hash in seen_urls):
            continue
        content = item.get("content") or item.get("summary") or ""
        normalized = normalize_text(content)
        entry = {
            "kind": "article",
            "source_id": sid,
            "url": item.get("url", ""),
            "url_hash": url_hash,
            "title": item.get("title", ""),
            "title_hash": digest(normalize_text(item.get("title", ""))),
            "content_hash": digest(normalized) if normalized else "",
            "lead_hash": digest(normalized[:100]) if normalized else "",
            "sketch": sketch(content) if normalized else [],
            "source_domain": urlparse(item.get("url", "")).netloc.lower(),
            "symbols": item.get("symbols") or item.get("stock") or [],
            "event_key": item.get("event_key", ""),
            "published_at": item.get("published_at", ""),
            "first_seen_at": stamp,
            "decision": item.get("decision", "skipped"),
            "call": item.get("call", ""),
        }
        records.append(entry)
        if sid:
            seen_ids.add(sid)
        if url_hash:
            seen_urls.add(url_hash)
        added += 1

    save(path, records)
    if not args.quiet:
        print(f"seen.py: 新增 {added} 筆，略過重複 {len(items) - added} 筆", file=sys.stderr)
    return 0


def cmd_similar(args: argparse.Namespace) -> int:
    """找內容相近的舊文，抓 exact 雜湊漏掉的改寫轉載稿。

    只標出候選並附相似度，最後是不是同一則由判讀的人或 agent 決定——
    大幅改寫的稿子分數會落在 0.25 附近，機器不該替它下定論。
    """
    _, records = open_db(quiet=args.quiet)
    target = sketch(sys.stdin.read())
    matches = []
    for record in records:
        if record.get("kind") == "event" or not record.get("sketch"):
            continue
        score = similarity(target, record["sketch"])
        if score >= args.threshold:
            matches.append(
                {
                    "similarity": round(score, 3),
                    "title": record.get("title", ""),
                    "url": record.get("url", ""),
                    "event_key": record.get("event_key", ""),
                }
            )
    matches.sort(key=lambda m: -m["similarity"])
    print(json.dumps({"threshold": args.threshold, "matches": matches}, ensure_ascii=False, indent=2))
    return 0


def cmd_event(args: argparse.Namespace) -> int:
    """事件層登記：把同一件事的多篇報導收斂成一條個股詳情。

    鍵的形狀是「代號|事件類型|期間」，例如 2330|月營收|2026-07。

    這裡查到已登記，**不代表該檔不能再出現在多空清單**——報告呈現的是當天的
    判斷狀態，同一檔連續多天看漲是正常的。事件登記只影響個股詳情怎麼寫：同一
    件事不必用「剛剛公布」重述一次，可以接成「延續昨日的營收利多」。
    """
    path, records = open_db(quiet=args.quiet)
    existing = {r["event_key"]: r for r in records if r.get("kind") == "event"}

    if args.action == "check":
        hit = existing.get(args.key)
        print(json.dumps({"key": args.key, "reported": bool(hit), "record": hit}, ensure_ascii=False))
        return 1 if hit else 0

    if args.key in existing:
        if not args.quiet:
            print(f"seen.py: 事件 {args.key} 已存在，未重複登記", file=sys.stderr)
        return 0
    records.append(
        {
            "kind": "event",
            "event_key": args.key,
            "note": args.note or "",
            "first_seen_at": now().isoformat(),
        }
    )
    save(path, records)
    if not args.quiet:
        print(f"seen.py: 登記事件 {args.key}", file=sys.stderr)
    return 0


def cmd_stats(args: argparse.Namespace) -> int:
    _, records = open_db(quiet=True)
    articles = [r for r in records if r.get("kind") != "event"]
    events = [r for r in records if r.get("kind") == "event"]
    decisions: dict[str, int] = {}
    domains: dict[str, int] = {}
    for record in articles:
        decisions[record.get("decision", "unknown")] = decisions.get(record.get("decision", "unknown"), 0) + 1
        domains[record.get("source_domain", "")] = domains.get(record.get("source_domain", ""), 0) + 1
    print(
        json.dumps(
            {
                "path": str(data_path()),
                "articles": len(articles),
                "events": len(events),
                "decisions": decisions,
                "by_domain": dict(sorted(domains.items(), key=lambda kv: -kv[1])),
                "ttl_days": {"article": ARTICLE_TTL_DAYS, "event": EVENT_TTL_DAYS},
            },
            ensure_ascii=False,
            indent=2,
        )
    )
    return 0


def cmd_prune(args: argparse.Namespace) -> int:
    path, records = open_db(quiet=True)
    save(path, records)
    print(f"seen.py: 清理完成，現存 {len(records)} 筆", file=sys.stderr)
    return 0


def main() -> int:
    # -q 放在子指令前後都要能用，否則實際打指令時很容易踩到順序限制。
    common = argparse.ArgumentParser(add_help=False)
    common.add_argument("-q", "--quiet", action="store_true", help="不輸出統計訊息")

    parser = argparse.ArgumentParser(
        description="Briefast 已讀新聞去重資料庫", parents=[common]
    )
    sub = parser.add_subparsers(dest="command", required=True)

    peek = sub.add_parser("peek", help="抓全文前先問哪些沒看過", parents=[common])
    peek.add_argument("--url", action="append", help="直接指定網址，可重複；省略時從 stdin 讀")
    peek.add_argument("--json", action="store_true", help="同時輸出已讀清單與判定依據")
    peek.set_defaults(func=cmd_peek)

    record = sub.add_parser("record", parents=[common], help="記錄處理過的新聞（含判定不報的）")
    record.set_defaults(func=cmd_record)

    similar = sub.add_parser("similar", parents=[common], help="找內容相近的舊文（轉載稿）")
    similar.add_argument("--threshold", type=float, default=0.25, help="Jaccard 相似度門檻，預設 0.25（不相關文章實測在 0.05 以下）")
    similar.set_defaults(func=cmd_similar)

    event = sub.add_parser("event", parents=[common], help="事件層去重")
    event.add_argument("action", choices=["check", "add"])
    event.add_argument("key", help="格式：代號|事件類型|期間")
    event.add_argument("--note", help="附註")
    event.set_defaults(func=cmd_event)

    sub.add_parser("stats", parents=[common], help="看目前存了什麼").set_defaults(func=cmd_stats)
    sub.add_parser("prune", parents=[common], help="只做清理").set_defaults(func=cmd_prune)

    args = parser.parse_args()
    return args.func(args)


if __name__ == "__main__":
    sys.exit(main())
