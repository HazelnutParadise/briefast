// Package dedup holds the text and URL fingerprints the /api/seen endpoints
// use to recognise articles the daily brief has already processed.
package dedup

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"golang.org/x/text/unicode/norm"
)

// SketchSize is the number of smallest bigram hashes kept per article. With 64
// the Jaccard estimate stays within about 0.03 of the exact value.
const SketchSize = 64

var (
	trackingPrefixes = []string{"utm_", "pk_", "mtm_"}
	trackingKeys     = map[string]bool{"fbclid": true, "gclid": true, "igshid": true, "ref": true, "share": true, "from": true, "src": true, "spm": true}

	htmlTag    = regexp.MustCompile(`<[^>]+>`)
	whitespace = regexp.MustCompile(`[\s\p{Z}]+`)
	// Bylines and copyright tails differ between reprints of the same story.
	boilerplate = regexp.MustCompile(`記者[\x{4e00}-\x{9fff}]{2,4}[／/][^\n]{0,12}報導` +
		`|[\x{4e00}-\x{9fff}]{2,6}／(?:綜合|台北|編譯)報導` +
		`|責任編輯[：:][^\n]{0,12}` +
		`|本文(?:不構成|僅供)[^\n]{0,30}` +
		`|版權所有[^\n]{0,30}`)
	punctuation = regexp.MustCompile(`[\s\p{Z}，。、：；！？（）「」【】《》%,.\-—–~〜·]`)
)

// Digest is the short hash stored for URLs, titles and content.
func Digest(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:8])
}

// CanonicalURL drops tracking query keys, the fragment and a trailing slash so
// different links to one article collapse to the same string. The host is
// kept as is: m.cnyes.com and news.cnyes.com are left to the source ID.
func CanonicalURL(raw string) string {
	u, err := url.Parse(strings.TrimSpace(raw))
	if err != nil {
		return strings.TrimSpace(raw)
	}
	query := url.Values{}
	for key, values := range u.Query() {
		lower := strings.ToLower(key)
		if trackingKeys[lower] || hasTrackingPrefix(lower) {
			continue
		}
		query[key] = values
	}
	path := strings.TrimRight(u.Path, "/")
	if path == "" {
		path = "/"
	}
	clean := url.URL{Scheme: strings.ToLower(u.Scheme), Host: strings.ToLower(u.Host), Path: path, RawQuery: query.Encode()}
	return clean.String()
}

func hasTrackingPrefix(key string) bool {
	for _, prefix := range trackingPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

// NormalizeText strips tags and bylines, applies NFKC and collapses
// whitespace. Raw HTML hashes are useless because ad slots change every load.
func NormalizeText(value string) string {
	text := htmlTag.ReplaceAllString(value, " ")
	text = norm.NFKC.String(text)
	text = boilerplate.ReplaceAllString(text, " ")
	return strings.TrimSpace(whitespace.ReplaceAllString(text, " "))
}

// SourceID prefers the source's own article ID over the URL, because sites
// change URL layouts but keep article numbers.
func SourceID(explicit, rawURL string) string {
	if id := strings.TrimSpace(explicit); id != "" {
		return id
	}
	if rawURL == "" {
		return ""
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	host := strings.ReplaceAll(strings.ToLower(u.Host), "www.", "")
	var tail []string
	for _, segment := range strings.Split(u.Path, "/") {
		if segment != "" {
			tail = append(tail, segment)
		}
	}
	marker := ""
	for _, segment := range tail {
		if isDigits(segment) {
			marker = segment
		}
	}
	if marker == "" && len(tail) > 0 {
		marker = tail[len(tail)-1]
	}
	if marker == "" {
		return ""
	}
	return fmt.Sprintf("%s:%s", host, marker)
}

func isDigits(value string) bool {
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return value != ""
}

// Sketch returns the bottom-k MinHash signature of the text's character
// bigrams. Bigrams need no word segmentation for Chinese and separate
// reprints (about 0.9), heavy rewrites (about 0.27) and unrelated articles
// (about 0.04) clearly.
func Sketch(text string) []string {
	compact := []rune(punctuation.ReplaceAllString(NormalizeText(text), ""))
	seen := map[uint64]bool{}
	for i := 0; i+1 < len(compact); i++ {
		sum := sha256.Sum256([]byte(string(compact[i : i+2])))
		seen[binary.BigEndian.Uint64(sum[:8])] = true
	}
	hashes := make([]uint64, 0, len(seen))
	for value := range seen {
		hashes = append(hashes, value)
	}
	sort.Slice(hashes, func(i, j int) bool { return hashes[i] < hashes[j] })
	if len(hashes) > SketchSize {
		hashes = hashes[:SketchSize]
	}
	out := make([]string, len(hashes))
	for i, value := range hashes {
		out[i] = fmt.Sprintf("%016x", value)
	}
	return out
}

// Similarity estimates the Jaccard similarity of two sketches.
func Similarity(left, right []string) float64 {
	if len(left) == 0 || len(right) == 0 {
		return 0
	}
	a, b := toSet(left), toSet(right)
	union := map[string]bool{}
	for value := range a {
		union[value] = true
	}
	for value := range b {
		union[value] = true
	}
	merged := make([]string, 0, len(union))
	for value := range union {
		merged = append(merged, value)
	}
	sort.Strings(merged)
	if len(merged) > SketchSize {
		merged = merged[:SketchSize]
	}
	shared := 0
	for _, value := range merged {
		if a[value] && b[value] {
			shared++
		}
	}
	return float64(shared) / float64(len(merged))
}

func toSet(values []string) map[string]bool {
	out := make(map[string]bool, len(values))
	for _, value := range values {
		out[value] = true
	}
	return out
}
