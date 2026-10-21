package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Document — контент сайта на всех поддерживаемых языках.
type Document struct {
	RU Content `json:"ru"`
	EN Content `json:"en"`
}

// Store хранит документ в памяти и сохраняет его в JSON на диск.
// Запись атомарная (tmp + rename) с резервной копией предыдущей версии.
type Store struct {
	mu      sync.RWMutex
	path    string
	doc     Document
	savedAt time.Time
	rev     int
}

var slugRe = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,40}$`)

func NewStore(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("каталог данных: %w", err)
	}
	s := &Store{
		path: filepath.Join(dir, "content.json"),
		doc:  Document{RU: defaultRU(), EN: defaultEN()},
	}

	raw, err := os.ReadFile(s.path)
	switch {
	case errors.Is(err, os.ErrNotExist):
		if err := s.persist(); err != nil {
			return nil, err
		}
		slog.Info("создан файл контента", "path", s.path)
	case err != nil:
		return nil, fmt.Errorf("чтение %s: %w", s.path, err)
	default:
		var doc Document
		if err := json.Unmarshal(raw, &doc); err != nil {
			broken := fmt.Sprintf("%s.broken-%s", s.path, time.Now().Format("20060102-150405"))
			_ = os.WriteFile(broken, raw, 0o600)
			slog.Error("content.json повреждён, использую значения по умолчанию",
				"err", err, "backup", broken)
			if err := s.persist(); err != nil {
				return nil, err
			}
		} else {
			doc.RU.Lang = "ru"
			doc.EN.Lang = "en"
			s.doc = doc
			s.rev = 1
			// Время последнего сохранения берём у файла, иначе админка
			// покажет нулевую дату.
			if info, statErr := os.Stat(s.path); statErr == nil {
				s.savedAt = info.ModTime()
			}
			slog.Info("контент загружен", "path", s.path, "projects", len(doc.RU.Projects))
		}
	}
	return s, nil
}

// Content возвращает копию контента для языка.
func (s *Store) Content(lang string) (Content, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	switch lang {
	case "ru":
		return s.doc.RU, true
	case "en":
		return s.doc.EN, true
	}
	return Content{}, false
}

// Snapshot отдаёт весь документ (для админки).
func (s *Store) Snapshot() (Document, time.Time, int) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.doc, s.savedAt, s.rev
}

// Save валидирует и сохраняет новый документ.
func (s *Store) Save(doc Document) error {
	if err := normalizeDocument(&doc); err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	prev := s.doc
	s.doc = doc
	if err := s.persist(); err != nil {
		s.doc = prev // откат: на диске и в памяти должно быть одно и то же
		return err
	}
	s.rev++
	return nil
}

// persist пишет файл. Вызывается под mu (или при инициализации).
func (s *Store) persist() error {
	data, err := json.MarshalIndent(s.doc, "", "  ")
	if err != nil {
		return fmt.Errorf("сериализация: %w", err)
	}
	data = append(data, '\n')

	tmp := s.path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return fmt.Errorf("запись tmp: %w", err)
	}
	// Резервная копия предыдущей версии — до подмены файла.
	if _, statErr := os.Stat(s.path); statErr == nil {
		if err := copyFile(s.path, s.path+".bak"); err != nil {
			slog.Warn("не удалось сделать резервную копию", "err", err)
		}
	}
	if err := os.Rename(tmp, s.path); err != nil {
		return fmt.Errorf("переименование: %w", err)
	}
	s.savedAt = time.Now()
	return nil
}

func copyFile(src, dst string) error {
	raw, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, raw, 0o600)
}

// normalizeDocument приводит документ в согласованный вид и проверяет
// то, без чего сайт нельзя отрендерить.
func normalizeDocument(doc *Document) error {
	docs := map[string]*Content{"ru": &doc.RU, "en": &doc.EN}

	for lang, c := range docs {
		c.Lang = lang
		trimContent(c)

		if c.Meta.Title == "" {
			return fmt.Errorf("[%s] заголовок страницы обязателен", lang)
		}
		if len(c.Meta.Title) > 200 {
			return fmt.Errorf("[%s] заголовок страницы длиннее 200 символов", lang)
		}
		if len(c.Projects) == 0 {
			return fmt.Errorf("[%s] нужен хотя бы один проект", lang)
		}

		seen := make(map[string]bool, len(c.Projects))
		for i := range c.Projects {
			p := &c.Projects[i]
			if !slugRe.MatchString(p.ID) {
				return fmt.Errorf("[%s] id проекта «%s»: только латиница, цифры и дефис", lang, p.ID)
			}
			if seen[p.ID] {
				return fmt.Errorf("[%s] id проекта «%s» повторяется", lang, p.ID)
			}
			seen[p.ID] = true
			if p.Name == "" {
				return fmt.Errorf("[%s] у проекта «%s» нет названия", lang, p.ID)
			}
		}

		if err := validateLinks(c, lang); err != nil {
			return err
		}
	}
	return nil
}

func validateLinks(c *Content, lang string) error {
	links := []struct {
		where string
		url   string
	}{
		{"hero.primary", c.Hero.Primary.Href},
		{"hero.secondary", c.Hero.Secondary.Href},
		{"contact", ""},
	}
	for _, l := range links {
		if l.url != "" && !isSafeHref(l.url) {
			return fmt.Errorf("[%s] %s: ссылка должна начинаться с https://, http://, / или #", lang, l.where)
		}
	}
	for i := range c.Projects {
		p := &c.Projects[i]
		if p.Demo != "" && !isSafeHref(p.Demo) {
			return fmt.Errorf("[%s] проект «%s»: ссылка должна начинаться с https://, http://, / или #", lang, p.Name)
		}
	}
	for i := range c.Contact.Items {
		it := &c.Contact.Items[i]
		if it.Href != "" && !isSafeHref(it.Href) {
			return fmt.Errorf("[%s] контакт «%s»: недопустимая ссылка", lang, it.Label)
		}
	}
	for i := range c.Nav {
		n := &c.Nav[i]
		if n.Href != "" && !isSafeHref(n.Href) {
			return fmt.Errorf("[%s] навигация «%s»: недопустимая ссылка", lang, n.Label)
		}
	}
	return nil
}

// isSafeHref не даёт сохранить javascript:-ссылку, которая станет XSS.
func isSafeHref(href string) bool {
	h := strings.ToLower(strings.TrimSpace(href))
	switch {
	case strings.HasPrefix(h, "https://"),
		strings.HasPrefix(h, "http://"),
		strings.HasPrefix(h, "/"),
		strings.HasPrefix(h, "#"),
		strings.HasPrefix(h, "mailto:"),
		strings.HasPrefix(h, "tel:"):
		return true
	}
	return false
}

// trimContent обрезает пробелы и ограничивает длину полей, чтобы в файл
// нельзя было записать мегабайт текста.
func trimContent(c *Content) {
	clamp := func(s *string, max int) { *s = clampString(*s, max) }

	clamp(&c.Meta.Title, 200)
	clamp(&c.Meta.Description, 400)

	for i := range c.Nav {
		clamp(&c.Nav[i].Label, 60)
		clamp(&c.Nav[i].Href, 400)
	}

	clamp(&c.Hero.Kicker, 80)
	clamp(&c.Hero.Title1, 120)
	clamp(&c.Hero.Title2, 120)
	clamp(&c.Hero.Lead, 800)
	clamp(&c.Hero.Note, 300)
	clamp(&c.Hero.Primary.Label, 80)
	clamp(&c.Hero.Primary.Href, 400)
	clamp(&c.Hero.Secondary.Label, 80)
	clamp(&c.Hero.Secondary.Href, 400)
	for i := range c.Hero.Facts {
		clamp(&c.Hero.Facts[i].Key, 60)
		clamp(&c.Hero.Facts[i].Value, 200)
	}
	c.Hero.Facts = limitSlice(c.Hero.Facts, 8)

	for i := range c.Stats {
		clamp(&c.Stats[i].Value, 40)
		clamp(&c.Stats[i].Label, 120)
	}
	c.Stats = limitSlice(c.Stats, 8)

	clamp(&c.Work.Kicker, 80)
	clamp(&c.Work.Title, 160)
	clamp(&c.Work.Lead, 800)

	c.Projects = limitSlice(c.Projects, 40)
	for i := range c.Projects {
		p := &c.Projects[i]
		clamp(&p.ID, 41)
		clamp(&p.Name, 120)
		clamp(&p.Kind, 80)
		clamp(&p.Tagline, 400)
		clamp(&p.Summary, 4000)
		clamp(&p.Demo, 400)
		clamp(&p.DemoText, 80)
		clamp(&p.Year, 40)
		clamp(&p.Status, 40)
		for j := range p.Bullets {
			clamp(&p.Bullets[j], 600)
		}
		p.Bullets = limitSlice(p.Bullets, 12)
		for j := range p.Tags {
			clamp(&p.Tags[j], 40)
		}
		p.Tags = limitSlice(p.Tags, 12)
	}

	clamp(&c.Approach.Kicker, 80)
	clamp(&c.Approach.Title, 160)
	clamp(&c.Approach.Lead, 800)
	c.Approach.Items = limitSlice(c.Approach.Items, 20)
	for i := range c.Approach.Items {
		clamp(&c.Approach.Items[i].Title, 200)
		clamp(&c.Approach.Items[i].Body, 1200)
	}

	clamp(&c.Contact.Kicker, 80)
	clamp(&c.Contact.Title, 160)
	clamp(&c.Contact.Lead, 800)
	clamp(&c.Contact.Note, 300)
	c.Contact.Items = limitSlice(c.Contact.Items, 10)
	for i := range c.Contact.Items {
		clamp(&c.Contact.Items[i].Label, 60)
		clamp(&c.Contact.Items[i].Value, 200)
		clamp(&c.Contact.Items[i].Href, 400)
		clamp(&c.Contact.Items[i].Copied, 60)
	}

	clamp(&c.Footer.Note, 300)
	clamp(&c.Footer.Built, 80)
}

func clampString(s string, max int) string {
	s = strings.TrimSpace(s)
	if len([]rune(s)) > max {
		s = string([]rune(s)[:max])
	}
	return s
}

func limitSlice[T any](items []T, max int) []T {
	if len(items) > max {
		return items[:max]
	}
	return items
}
