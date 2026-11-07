package dedup

import "testing"

const articleBody = `台積電今日公布七月營收，單月合併營收達新台幣三千二百億元，較上月成長百分之十二，較去年同期成長百分之四十五，創下歷史新高。公司表示，先進製程需求持續強勁，三奈米與五奈米產能利用率維持高檔，人工智慧加速器訂單帶動高效能運算平台營收占比進一步提升。法人預期第三季營收可望優於財測上緣，毛利率也有機會維持在五成以上。`

const unrelatedBody = `中央氣象署表示，颱風預計明天清晨在花蓮一帶登陸，北部及東半部地區將有豪雨，山區累積雨量可能超過五百毫米。各縣市政府已陸續宣布停班停課，民眾應避免前往海邊與山區，並留意土石流警戒區的撤離通知。交通部也提醒，部分國道客運與台鐵班次將調整或停駛。`

func TestCanonicalURLDropsTrackingAndTrailingSlash(t *testing.T) {
	cases := map[string]string{
		"https://News.Cnyes.com/news/id/6001/?utm_source=line&fbclid=x#top": "https://news.cnyes.com/news/id/6001",
		"https://news.cnyes.com/news/id/6001?page=2&SPM=a":                  "https://news.cnyes.com/news/id/6001?page=2",
		"https://example.com": "https://example.com/",
	}
	for in, want := range cases {
		if got := CanonicalURL(in); got != want {
			t.Errorf("CanonicalURL(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestSourceIDFallsBackToNumericSegment(t *testing.T) {
	cases := []struct{ explicit, url, want string }{
		{"cnyes:6001", "https://news.cnyes.com/news/id/9999", "cnyes:6001"},
		{"", "https://www.cnyes.com/news/id/6001/slug", "cnyes.com:6001"},
		{"", "https://example.com/story/abc", "example.com:abc"},
		{"", "", ""},
	}
	for _, c := range cases {
		if got := SourceID(c.explicit, c.url); got != c.want {
			t.Errorf("SourceID(%q, %q) = %q, want %q", c.explicit, c.url, got, c.want)
		}
	}
}

func TestNormalizeTextRemovesTagsBylinesAndWidth(t *testing.T) {
	got := NormalizeText("<p>記者王小明／台北報導</p>　台積電ＡＩ　 營收\n\n成長 責任編輯：李四")
	if got != "台積電AI 營收 成長" {
		t.Fatalf("NormalizeText = %q", got)
	}
}

func TestSimilaritySeparatesReprintsFromUnrelated(t *testing.T) {
	original := Sketch(articleBody)
	reprint := Sketch("記者王小明／台北報導 " + articleBody + " 版權所有，轉載請註明")
	if score := Similarity(original, reprint); score < 0.9 {
		t.Fatalf("照登轉載相似度 = %.3f, want >= 0.9", score)
	}
	if score := Similarity(original, Sketch(unrelatedBody)); score >= 0.25 {
		t.Fatalf("不相關文章相似度 = %.3f, want < 0.25", score)
	}
	if len(original) != SketchSize {
		t.Fatalf("sketch 長度 = %d, want %d", len(original), SketchSize)
	}
	if Similarity(nil, original) != 0 {
		t.Fatal("空 sketch 應回 0")
	}
}
