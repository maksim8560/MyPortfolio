package main

import "strings"

// supportedLangs — порядок важен: первый считается языком по умолчанию.
var supportedLangs = []string{"ru", "en"}

// pickLang определяет язык страницы.
//
// Приоритет: явный ?lang= (если поддерживается) → Accept-Language →
// язык по умолчанию. Неизвестный ?lang= не ломает страницу — просто
// игнорируется.
func pickLang(query, acceptLanguage string) string {
	if l := normalizeLang(query); l != "" && isSupported(l) {
		return l
	}
	for _, part := range strings.Split(acceptLanguage, ",") {
		tag := part
		if i := strings.IndexByte(tag, ';'); i >= 0 {
			tag = tag[:i]
		}
		if l := normalizeLang(tag); l != "" && isSupported(l) {
			return l
		}
	}
	return supportedLangs[0]
}

func normalizeLang(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	if i := strings.IndexAny(s, "-_"); i >= 0 {
		s = s[:i]
	}
	return s
}

func isSupported(lang string) bool {
	for _, l := range supportedLangs {
		if l == lang {
			return true
		}
	}
	return false
}
