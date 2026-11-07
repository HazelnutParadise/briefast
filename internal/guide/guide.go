// Package guide holds the daily-brief workflow guide served at GET /api/guide.
// It ships inside the binary so the guide and the validation rules it
// describes are always deployed together.
package guide

import _ "embed"

//go:embed daily-brief.md
var DailyBrief string
