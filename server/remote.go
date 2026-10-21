package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"
)

// remoteStore — контент в Upstash Redis поверх REST API.
//
// Нужен на хостингах без постоянного диска (Render free): файл внутри
// контейнера стирается при каждом рестарте и засыпании инстанса.
// Документ хранится одной строкой JSON под ключом UPSTASH_REDIS_KEY.
type remoteStore struct {
	url   string
	token string
	key   string
	http  *http.Client
}

// newRemoteStore возвращает nil, если переменные окружения не заданы:
// тогда сервер работает как раньше, с одним файлом на диске.
func newRemoteStore() *remoteStore {
	url := strings.TrimRight(strings.TrimSpace(os.Getenv("UPSTASH_REDIS_REST_URL")), "/")
	token := strings.TrimSpace(os.Getenv("UPSTASH_REDIS_REST_TOKEN"))
	if url == "" || token == "" {
		return nil
	}

	key := strings.TrimSpace(os.Getenv("UPSTASH_REDIS_KEY"))
	if key == "" {
		key = "portfolio:content"
	}

	return &remoteStore{
		url:   url,
		token: token,
		key:   key,
		http:  &http.Client{Timeout: 15 * time.Second},
	}
}

func (r *remoteStore) Name() string { return "upstash" }

// command отправляет одну команду в корень REST-API: тело — сама команда,
// поэтому 20 КБ JSON едет в body, а не в URL.
func (r *remoteStore) command(ctx context.Context, args ...any) (json.RawMessage, error) {
	payload, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("upstash: сборка команды: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, r.url, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("upstash: запрос: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("upstash: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 16<<20))
	if err != nil {
		return nil, fmt.Errorf("upstash: чтение ответа: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("upstash: %s: %s", resp.Status, strings.TrimSpace(string(body)))
	}

	var reply struct {
		Result json.RawMessage `json:"result"`
		Error  string          `json:"error"`
	}
	if err := json.Unmarshal(body, &reply); err != nil {
		return nil, fmt.Errorf("upstash: разбор ответа: %w", err)
	}
	if reply.Error != "" {
		return nil, fmt.Errorf("upstash: %s", reply.Error)
	}
	return reply.Result, nil
}

// Load возвращает сохранённый документ. found=false — ключа ещё нет.
func (r *remoteStore) Load(ctx context.Context) (data []byte, found bool, err error) {
	raw, err := r.command(ctx, "GET", r.key)
	if err != nil {
		return nil, false, err
	}
	if len(raw) == 0 || string(raw) == "null" {
		return nil, false, nil
	}

	var value string
	if err := json.Unmarshal(raw, &value); err != nil {
		return nil, false, fmt.Errorf("upstash: значение не строка: %w", err)
	}
	if strings.TrimSpace(value) == "" {
		return nil, false, nil
	}
	return []byte(value), true, nil
}

func (r *remoteStore) Save(ctx context.Context, data []byte) error {
	if _, err := r.command(ctx, "SET", r.key, string(data)); err != nil {
		return err
	}
	slog.Debug("контент записан в облако", "key", r.key, "bytes", len(data))
	return nil
}
