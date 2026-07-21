package i18n

import (
	"net/http/httptest"
	"testing"
)

func newTestTranslator(t *testing.T) *Translator {
	t.Helper()
	tr, err := Load("en")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return tr
}

func TestT(t *testing.T) {
	tr := newTestTranslator(t)

	if got := tr.T("en", "nav.home"); got != "Home" {
		t.Errorf("en nav.home = %q, want %q", got, "Home")
	}
	if got := tr.T("de", "nav.home"); got != "Startseite" {
		t.Errorf("de nav.home = %q, want %q", got, "Startseite")
	}
	// Unknown language falls back to the default language.
	if got := tr.T("fr", "nav.home"); got != "Home" {
		t.Errorf("fr nav.home = %q, want fallback %q", got, "Home")
	}
	// Unknown key falls back to the key itself.
	if got := tr.T("en", "does.not.exist"); got != "does.not.exist" {
		t.Errorf("unknown key = %q, want the key itself", got)
	}
}

func TestResolve(t *testing.T) {
	tr := newTestTranslator(t)

	// Query parameter wins.
	r := httptest.NewRequest("GET", "/?lang=de", nil)
	if got := tr.Resolve(r); got != "de" {
		t.Errorf("query lang: got %q, want de", got)
	}

	// Cookie is used when no query parameter is set.
	r = httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Cookie", CookieName+"=de")
	if got := tr.Resolve(r); got != "de" {
		t.Errorf("cookie lang: got %q, want de", got)
	}

	// Accept-Language header, including region codes like de-AT.
	r = httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Language", "fr-FR,de-AT;q=0.8")
	if got := tr.Resolve(r); got != "de" {
		t.Errorf("accept-language: got %q, want de", got)
	}

	// Nothing supported falls back to the default language.
	r = httptest.NewRequest("GET", "/", nil)
	r.Header.Set("Accept-Language", "fr,es")
	if got := tr.Resolve(r); got != "en" {
		t.Errorf("fallback: got %q, want en", got)
	}
}
