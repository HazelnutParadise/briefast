package site

import (
	"bytes"
	"crypto/sha256"
	_ "embed"
	"fmt"
	"net/http"
	"time"
)

//go:embed assets/editorial-sectors-v1.webp
var editorialIllustrations []byte

// EditorialIllustrationHandler serves the newspaper artwork bundled into the binary.
func EditorialIllustrationHandler() http.Handler {
	etag := fmt.Sprintf(`"%x"`, sha256.Sum256(editorialIllustrations))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/webp")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		w.Header().Set("ETag", etag)
		http.ServeContent(w, r, "editorial-sectors-v1.webp", time.Time{}, bytes.NewReader(editorialIllustrations))
	})
}
