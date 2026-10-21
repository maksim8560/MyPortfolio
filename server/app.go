package main

import (
	"encoding/json"
	"net/http"
	"time"
)

// App — состояние сервера: версия, время старта, хранилище контента и вход.
type App struct {
	version   string
	startedAt time.Time
	store     *Store
	auth      *Auth
}

func newApp(version string, store *Store, auth *Auth) *App {
	return &App{version: version, startedAt: time.Now(), store: store, auth: auth}
}

func (a *App) routes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/healthz", a.handleHealth)
	mux.HandleFunc("GET /api/content", a.handleContent)
	mux.HandleFunc("GET /api/projects", a.handleProjects)
	a.adminRoutes(mux)
	// Всё остальное — статика и SPA-роутинг.
	mux.Handle("/", spaHandler())
}

func (a *App) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"status":  "ok",
		"version": a.version,
		"uptime":  int64(time.Since(a.startedAt).Seconds()),
		"runtime": runtimeVersion(),
	})
}

// handleContent отдаёт весь контент страницы на выбранном языке.
// Язык выбирается так: ?lang= → Accept-Language → ru.
func (a *App) handleContent(w http.ResponseWriter, r *http.Request) {
	lang := pickLang(r.URL.Query().Get("lang"), r.Header.Get("Accept-Language"))
	c, ok := a.store.Content(lang)
	if !ok {
		http.Error(w, "unsupported lang", http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

// handleProjects оставлен отдельным эндпоинтом: карточки можно обновлять
// независимо от остальной страницы.
func (a *App) handleProjects(w http.ResponseWriter, r *http.Request) {
	lang := pickLang(r.URL.Query().Get("lang"), r.Header.Get("Accept-Language"))
	c, ok := a.store.Content(lang)
	if !ok {
		http.Error(w, "unsupported lang", http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"lang":     c.Lang,
		"projects": c.Projects,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// Заголовки уже отправлены — остаётся только записать в лог.
		slogError("не удалось записать JSON", err)
	}
}
