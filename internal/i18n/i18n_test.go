package i18n

import (
	"net/http/httptest"
	"testing"
)

func newTestTranslator(t *testing.T) *Translator {
	t.Helper()
	tr, err := Load("de")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	return tr
}

func TestT(t *testing.T) {
	tr := newTestTranslator(t)

	if got := tr.T("de", "nav.welcome"); got != "Willkommen" {
		t.Errorf("de nav.welcome = %q, want %q", got, "Willkommen")
	}
	// Unknown language falls back to the default language.
	if got := tr.T("fr", "nav.welcome"); got != "Willkommen" {
		t.Errorf("fr nav.welcome = %q, want fallback %q", got, "Willkommen")
	}
	// Unknown key falls back to the key itself.
	if got := tr.T("de", "does.not.exist"); got != "does.not.exist" {
		t.Errorf("unknown key = %q, want the key itself", got)
	}
}

func TestList(t *testing.T) {
	tr := newTestTranslator(t)

	items := tr.List("de", "home.about.paragraphs")
	if len(items) != 3 {
		t.Fatalf("home.about.paragraphs: got %d items, want 3", len(items))
	}
	if items := tr.List("de", "does.not.exist"); len(items) != 0 {
		t.Errorf("unknown list key: got %d items, want 0", len(items))
	}
}

func TestResolve(t *testing.T) {
	tr := newTestTranslator(t)

	// Cookie sets the language when supported.
	r := httptest.NewRequest("GET", "/", nil)
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
	if got := tr.Resolve(r); got != "de" {
		t.Errorf("fallback: got %q, want de", got)
	}
}
