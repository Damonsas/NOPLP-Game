package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
)

func AssetHandler() http.Handler {
	fs := http.StripPrefix("/asset/", http.FileServer(http.Dir("asset")))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ext := strings.ToLower(filepath.Ext(r.URL.Path))

		switch ext {
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		}

		fs.ServeHTTP(w, r)
	})
}
