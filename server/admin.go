package main

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
)

const maxContentBytes = 1 << 20 // 1 МБ на документ — с запасом

func (a *App) adminRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/admin/login", a.handleLogin)
	mux.HandleFunc("POST /api/admin/logout", a.handleLogout)
	mux.HandleFunc("GET /api/admin/session", a.handleSession)
	mux.HandleFunc("GET /api/admin/content", a.auth.requireAdmin(a.handleAdminGet))
	mux.HandleFunc("PUT /api/admin/content", a.auth.requireAdmin(a.handleAdminPut))
}

func (a *App) handleLogin(w http.ResponseWriter, r *http.Request) {
	if !isSameOrigin(r) {
		http.Error(w, "cross-origin request", http.StatusForbidden)
		return
	}

	var body struct {
		Password string `json:"password"`
	}
	if err := decodeJSON(w, r, &body); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	token, ok, retry := a.auth.Login(body.Password, clientIP(r))
	if !ok {
		if retry > 0 {
			w.Header().Set("Retry-After", "300")
			writeJSON(w, http.StatusTooManyRequests, map[string]string{
				"error": "слишком много попыток, подожди пару минут",
			})
			return
		}
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "неверный пароль"})
		return
	}

	a.auth.setCookie(w, r, token)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	if cookie, err := r.Cookie(sessionCookie); err == nil {
		a.auth.Logout(cookie.Value)
	}
	a.auth.clearCookie(w, r)
	writeJSON(w, http.StatusOK, map[string]bool{"ok": true})
}

func (a *App) handleSession(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]bool{
		"authenticated": a.auth.authenticated(r),
	})
}

func (a *App) handleAdminGet(w http.ResponseWriter, r *http.Request) {
	doc, savedAt, rev := a.store.Snapshot()
	w.Header().Set("Cache-Control", "no-store")
	writeJSON(w, http.StatusOK, map[string]any{
		"content": doc,
		"rev":     rev,
		"savedAt": savedAt,
	})
}

func (a *App) handleAdminPut(w http.ResponseWriter, r *http.Request) {
	if !isSameOrigin(r) {
		http.Error(w, "cross-origin request", http.StatusForbidden)
		return
	}
	// Кастомный заголовок — простая CSRF-защита: чужой сайт не сможет
	// отправить такой запрос без CORS-preflight, который мы не разрешаем.
	if r.Header.Get("X-Portfolio-Admin") != "1" {
		http.Error(w, "missing admin header", http.StatusForbidden)
		return
	}

	var doc Document
	if err := decodeJSON(w, r, &doc); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.store.Save(doc); err != nil {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	_, savedAt, rev := a.store.Snapshot()
	slog.Info("контент обновлён", "rev", rev, "projects", len(doc.RU.Projects))
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":      true,
		"rev":     rev,
		"savedAt": savedAt,
	})
}

// decodeJSON читает тело с ограничением размера и требует объект.
func decodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	if ct := r.Header.Get("Content-Type"); ct != "" &&
		!strings.HasPrefix(strings.ToLower(ct), "application/json") {
		return errors.New("ожидается Content-Type: application/json")
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxContentBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		var maxErr *http.MaxBytesError
		if errors.As(err, &maxErr) {
			return errors.New("слишком большой запрос")
		}
		return errors.New("некорректный JSON: " + err.Error())
	}
	return nil
}

// isSameOrigin отсекает кросс-доменные запросы браузера. Запросы без Origin
// (curl, тесты) пропускаются — за них отвечает заголовок X-Portfolio-Admin.
func isSameOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}
	origin = strings.TrimSuffix(origin, "/")
	if strings.EqualFold(origin, "http://"+r.Host) || strings.EqualFold(origin, "https://"+r.Host) {
		return true
	}
	return false
}
