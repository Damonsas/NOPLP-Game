package handlers

import (
	"net/http"
	"path/filepath"
	"strings"
)

func AssetHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		switch ext := strings.ToLower(filepath.Ext(path)); ext {
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		}

		http.StripPrefix("/asset/", http.FileServer(http.Dir("asset"))).ServeHTTP(w, r)

	})
}
