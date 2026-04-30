package i18n

import (
	"testing"
)

func TestT_ReturnsKeyWhenNoTranslation(t *testing.T) {
	// When no language is loaded, T() should return the key itself
	result := T("Hello World")
	if result != "Hello World" {
		t.Errorf("expected 'Hello World', got '%s'", result)
	}
}

func TestT_ReturnsTranslationWhenLoaded(t *testing.T) {
	// Load zh-CN and verify translation
	err := Init("zh-CN")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	result := T("Queued")
	if result != "排队中" {
		t.Errorf("expected '排队中', got '%s'", result)
	}
}

func TestT_FallbackWhenKeyMissing(t *testing.T) {
	err := Init("zh-CN")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	// This key doesn't exist in zh-CN.json
	result := T("Some nonexistent string")
	if result != "Some nonexistent string" {
		t.Errorf("expected fallback to key, got '%s'", result)
	}
}

func TestT_EnglishModeReturnsKey(t *testing.T) {
	err := Init("en")
	if err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	result := T("Queued")
	if result != "Queued" {
		t.Errorf("expected 'Queued', got '%s'", result)
	}
}

func TestSetLanguage(t *testing.T) {
	err := SetLanguage("zh-CN")
	if err != nil {
		t.Fatalf("SetLanguage failed: %v", err)
	}
	result := T("Queued")
	if result != "排队中" {
		t.Errorf("expected '排队中', got '%s'", result)
	}
}
