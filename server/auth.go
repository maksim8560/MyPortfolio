package main

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	sessionCookie = "portfolio_admin"
	sessionTTL    = 12 * time.Hour
	maxAttempts   = 8
	attemptWindow = 5 * time.Minute
)

// Auth — вход в админку по паролю. Сессии живут в памяти: перезапуск сервера
// разлогинивает всех, что для одного администратора скорее плюс.
type Auth struct {
	mu       sync.Mutex
	password string
	sessions map[string]time.Time
	attempts map[string][]time.Time
}

func NewAuth(dataDir string) (*Auth, error) {
	password, source, err := resolvePassword(dataDir)
	if err != nil {
		return nil, err
	}
	slog.Info("вход в админку готов", "источник пароля", source)
	return &Auth{
		password: password,
		sessions: make(map[string]time.Time),
		attempts: make(map[string][]time.Time),
	}, nil
}

// resolvePassword берёт пароль из ADMIN_PASSWORD, иначе из файла рядом с
// данными, иначе генерирует новый и сохраняет его.
func resolvePassword(dataDir string) (password, source string, err error) {
	if p := strings.TrimSpace(os.Getenv("ADMIN_PASSWORD")); p != "" {
		return p, "ADMIN_PASSWORD", nil
	}

	path := filepath.Join(dataDir, "admin-password.txt")
	raw, readErr := os.ReadFile(path)
	if readErr == nil {
		if p := strings.TrimSpace(string(raw)); p != "" {
			return p, path, nil
		}
	}

	generated, err := randomPassword(18)
	if err != nil {
		return "", "", err
	}
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		return "", "", err
	}
	if err := os.WriteFile(path, []byte(generated+"\n"), 0o600); err != nil {
		return "", "", err
	}
	slog.Warn("пароль администратора создан — сохрани его", "password", generated, "path", path)
	return generated, path, nil
}

func randomPassword(n int) (string, error) {
	const alphabet = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i := range buf {
		buf[i] = alphabet[int(buf[i])%len(alphabet)]
	}
	return string(buf), nil
}

// Login проверяет пароль и выдаёт токен сессии.
func (a *Auth) Login(password, ip string) (string, bool, time.Duration) {
	if !a.allow(ip) {
		return "", false, attemptWindow
	}

	a.mu.Lock()
	defer a.mu.Unlock()

	if subtle.ConstantTimeCompare([]byte(password), []byte(a.password)) != 1 {
		a.attempts[ip] = append(a.attempts[ip], time.Now())
		return "", false, 0
	}

	token, err := randomToken()
	if err != nil {
		slog.Error("не удалось выпустить токен", "err", err)
		return "", false, 0
	}
	a.sessions[token] = time.Now().Add(sessionTTL)
	delete(a.attempts, ip)
	return token, true, 0
}

// Verify проверяет токен сессии.
func (a *Auth) Verify(token string) bool {
	if token == "" {
		return false
	}
	a.mu.Lock()
	defer a.mu.Unlock()

	expires, ok := a.sessions[token]
	if !ok {
		return false
	}
	if time.Now().After(expires) {
		delete(a.sessions, token)
		return false
	}
	return true
}

func (a *Auth) Logout(token string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.sessions, token)
}

// allow ограничивает число попыток входа с одного адреса.
func (a *Auth) allow(ip string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()

	cutoff := time.Now().Add(-attemptWindow)
	recent := a.attempts[ip][:0]
	for _, t := range a.attempts[ip] {
		if t.After(cutoff) {
			recent = append(recent, t)
		}
	}
	if len(recent) == 0 {
		delete(a.attempts, ip)
	} else {
		a.attempts[ip] = recent
	}
	return len(recent) < maxAttempts
}

func randomToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

func (a *Auth) setCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   requestIsHTTPS(r),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(sessionTTL.Seconds()),
	})
}

func (a *Auth) clearCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   requestIsHTTPS(r),
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})
}

func requestIsHTTPS(r *http.Request) bool {
	return r.TLS != nil || strings.EqualFold(r.Header.Get("X-Forwarded-Proto"), "https")
}

func clientIP(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		if i := strings.IndexByte(fwd, ','); i > 0 {
			return strings.TrimSpace(fwd[:i])
		}
		return strings.TrimSpace(fwd)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// requireAdmin пускает дальше только с валидной сессией.
func (a *Auth) requireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(sessionCookie)
		if err != nil || !a.Verify(cookie.Value) {
			w.Header().Set("Cache-Control", "no-store")
			writeJSON(w, http.StatusUnauthorized, map[string]string{
				"error": "нужен вход в админку",
			})
			return
		}
		next(w, r)
	}
}

func (a *Auth) authenticated(r *http.Request) bool {
	cookie, err := r.Cookie(sessionCookie)
	return err == nil && a.Verify(cookie.Value)
}
