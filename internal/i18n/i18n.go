// Package i18n provides a minimal translation engine using English-as-key.
// Translations are loaded from embedded JSON files keyed by locale code.
// Call Init(lang) at startup, then use T(key) to translate.
// Only language files named by locale code (e.g., zh-CN.json) should exist in this directory.
package i18n

import (
	"embed"
	"encoding/json"
	"errors"
	"io/fs"
	"sync"
)

//go:embed *.json
var localeFS embed.FS

var (
	mu       sync.RWMutex
	messages map[string]string // current language translations
)

// Init loads translations for the given language code.
// lang should be "en", "zh-CN", etc.
// For "en", no file is loaded — T() returns keys as-is.
func Init(lang string) error {
	if lang == "" || lang == "en" {
		mu.Lock()
		messages = nil
		mu.Unlock()
		return nil
	}

	data, err := localeFS.ReadFile(lang + ".json")
	if err != nil {
		if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		// File doesn't exist, fall back to identity map (English)
		mu.Lock()
		messages = nil
		mu.Unlock()
		return nil
	}

	var dict map[string]string
	if err := json.Unmarshal(data, &dict); err != nil {
		return err
	}

	mu.Lock()
	messages = dict
	mu.Unlock()
	return nil
}

// SetLanguage is a convenience wrapper for Init, useful for hot-reload in tests.
func SetLanguage(lang string) error {
	return Init(lang)
}

// T translates key to the current language, or returns key itself if no translation found.
func T(key string) string {
	mu.RLock()
	dict := messages
	mu.RUnlock()

	if dict == nil {
		return key
	}
	if v, ok := dict[key]; ok {
		return v
	}
	return key
}
