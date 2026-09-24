package report

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// The skill's brief example is what the agent copies; it must pass the same
// validation the API applies, or the first real POST fails on a typo here.
func TestSkillConferenceExampleValidates(t *testing.T) {
	raw, err := os.ReadFile("../../skills/daily-brief/SKILL.md")
	if err != nil {
		t.Skip("skill file not available:", err)
	}
	section := string(raw)
	start := strings.Index(section, "### 組會前報告並 POST")
	if start < 0 {
		t.Fatal("brief section heading missing from SKILL.md")
	}
	section = section[start:]
	open := strings.Index(section, "```json\n")
	if open < 0 {
		t.Fatal("brief JSON example missing from SKILL.md")
	}
	section = section[open+len("```json\n"):]
	example := section[:strings.Index(section, "```")]

	var c Conference
	decoder := json.NewDecoder(strings.NewReader(example))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&c); err != nil {
		t.Fatalf("decode example: %v", err)
	}
	if errs := c.Validate(); len(errs) != 0 {
		t.Fatalf("example fails validation: %v", errs)
	}
}

func TestSkillDailyReportExampleValidatesMarketOutlook(t *testing.T) {
	raw, err := os.ReadFile("../../skills/daily-brief/SKILL.md")
	if err != nil {
		t.Skip("skill file not available:", err)
	}
	section := string(raw)
	start := strings.Index(section, "## 4. 組成並驗證報告 JSON")
	if start < 0 {
		t.Fatal("daily report section missing from SKILL.md")
	}
	section = section[start:]
	open := strings.Index(section, "```json\n")
	if open < 0 {
		t.Fatal("daily report JSON example missing from SKILL.md")
	}
	section = section[open+len("```json\n"):]
	close := strings.Index(section, "```")
	if close < 0 {
		t.Fatal("daily report JSON fence is not closed")
	}
	var r Report
	decoder := json.NewDecoder(strings.NewReader(section[:close]))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&r); err != nil {
		t.Fatalf("decode example: %v", err)
	}
	if r.MarketOutlook == nil || r.MarketOutlook.TrajectoryMD == nil || strings.TrimSpace(*r.MarketOutlook.TrajectoryMD) == "" {
		t.Fatal("daily report example has no market trajectory")
	}
	if errs := r.Validate(); len(errs) != 0 {
		t.Fatalf("example fails validation: %v", errs)
	}
}

func TestReadmeReportExampleValidates(t *testing.T) {
	raw, err := os.ReadFile("../../README.md")
	if err != nil {
		t.Skip("README not available:", err)
	}
	section := string(raw)
	start := strings.Index(section, "## 報告 JSON")
	if start < 0 {
		t.Fatal("report section missing from README")
	}
	section = section[start:]
	open := strings.Index(section, "```json\n")
	if open < 0 {
		t.Fatal("report JSON example missing from README")
	}
	section = section[open+len("```json\n"):]
	close := strings.Index(section, "```")
	if close < 0 {
		t.Fatal("report JSON fence is not closed")
	}
	var r Report
	decoder := json.NewDecoder(strings.NewReader(section[:close]))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&r); err != nil {
		t.Fatalf("decode example: %v", err)
	}
	if r.MarketOutlook == nil || r.MarketOutlook.TrajectoryMD == nil || strings.TrimSpace(*r.MarketOutlook.TrajectoryMD) == "" {
		t.Fatal("README example has no market trajectory")
	}
	if errs := r.Validate(); len(errs) != 0 {
		t.Fatalf("example fails validation: %v", errs)
	}
}
