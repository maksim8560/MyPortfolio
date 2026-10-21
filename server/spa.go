package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"path"
	"strings"
)

// dist — результат `npm run build` из каталога web.
//
//go:embed all:dist
var distFS embed.FS

// spaHandler раздаёт собранный фронтенд и отдаёт index.html на любой
// не-API маршрут, чтобы работал клиентский роутинг.
func spaHandler() http.Handler {
	sub, err := fs.Sub(distFS, "dist")
	if err != nil {
		slog.Error("не удалось открыть dist", "err", err)
		return http.NotFoundHandler()
	}
	files := http.FileServerFS(sub)
	index, err := fs.ReadFile(sub, "index.html")
	if err != nil {
		slog.Error("в dist нет index.html — фронтенд не собран", "err", err)
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			w.Header().Set("Allow", "GET, HEAD")
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		name := strings.TrimPrefix(path.Clean(r.URL.Path), "/")
		if name == "" || name == "." {
			serveIndex(w, r, index)
			return
		}

		if _, statErr := fs.Stat(sub, name); statErr != nil {
			serveIndex(w, r, index)
			return
		}

		// Хешированные ассеты Vite неизменяемы, остальное — перепроверяем.
		if strings.HasPrefix(name, "assets/") {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "public, max-age=3600")
		}
		files.ServeHTTP(w, r)
	})
}

func serveIndex(w http.ResponseWriter, r *http.Request, index []byte) {
	if index == nil {
		http.Error(w, "frontend not built", http.StatusServiceUnavailable)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	if r.Method == http.MethodHead {
		return
	}
	_, _ = w.Write(index)
}
