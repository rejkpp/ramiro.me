package handler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newTestRequest constructs an httptest request for the given path.
func newTestRequest(t *testing.T, path string) *http.Request {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatalf("build request: %v", err)
	}
	return req
}

func TestPublic_Home(t *testing.T) {
	p := NewPublic()
	rr := httptest.NewRecorder()
	p.Home(rr, newTestRequest(t, "/"))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	ct := rr.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html prefix", ct)
	}
	body := rr.Body.String()
	if !strings.Contains(body, "<html") {
		t.Errorf("body missing <html tag: %s", body)
	}
	if !strings.Contains(body, "htmx.min.js") {
		t.Errorf("body missing htmx script reference")
	}

	// New homepage sections
	for _, want := range []string{
		"Intelligence",
		"Intuition",
		"Impact",
		"I build software that thinks",
		"Book a session",
		"See my work",
		"AI and agents",
		"Algorithmic trading",
		"Accounting software",
		"Clairvoyant Meditation",
		"Currently building",
		"ORC",
		"Ramiro",
		"AI integrations",
		"Meditation",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("home body missing %q", want)
		}
	}
}

func TestPublic_About(t *testing.T) {
	p := NewPublic()
	rr := httptest.NewRecorder()
	p.About(rr, newTestRequest(t, "/about"))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if !strings.Contains(rr.Body.String(), "About") {
		t.Errorf("about body missing 'About' text")
	}
}

func TestPublic_Projects(t *testing.T) {
	p := NewPublic()
	rr := httptest.NewRecorder()
	p.Projects(rr, newTestRequest(t, "/projects"))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if !strings.Contains(rr.Body.String(), "Projects") {
		t.Errorf("projects body missing 'Projects' text")
	}
}

func TestPublic_Booking(t *testing.T) {
	p := NewPublic()
	rr := httptest.NewRecorder()
	p.Booking(rr, newTestRequest(t, "/booking"))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	body := rr.Body.String()
	for _, want := range []string{
		"Book a paid call",
		"30-minute call",
		"60-minute call",
		"CALENDLY_PLACEHOLDER_30_MIN",
		"CALENDLY_PLACEHOLDER_60_MIN",
		"Book 30 minutes",
		"Book 60 minutes",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("booking body missing %q", want)
		}
	}
	for _, unwanted := range []string{
		"<form",
		"Send a message",
		"Send message",
		"15-minute",
	} {
		if strings.Contains(body, unwanted) {
			t.Errorf("booking body unexpectedly contains %q", unwanted)
		}
	}
}

func TestPublic_Menu(t *testing.T) {
	p := NewPublic()
	rr := httptest.NewRecorder()
	p.Menu(rr, newTestRequest(t, "/partials/menu"))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	body := rr.Body.String()

	// Menu should contain nav links
	if !strings.Contains(body, `href="/about"`) {
		t.Errorf("menu body missing about link")
	}

	// Menu should NOT contain the old "Close" text button
	if strings.Contains(body, ">Close</") {
		t.Errorf("menu body still contains Close text link — should be removed")
	}

	// Menu should include an out-of-band swap to replace the hamburger with an X button
	if !strings.Contains(body, `hx-swap-oob="true"`) {
		t.Errorf("menu body missing hx-swap-oob for button swap")
	}
	if !strings.Contains(body, `id="menu-btn"`) {
		t.Errorf("menu body missing id=menu-btn for OOB swap target")
	}

	// The OOB-swapped button should point to menu-close
	if !strings.Contains(body, `/partials/menu-close`) {
		t.Errorf("menu body missing /partials/menu-close endpoint in OOB button")
	}

	// Menu should include an overlay for click-outside-to-close
	if !strings.Contains(body, `id="menu-overlay"`) {
		t.Errorf("menu body missing overlay for click-outside-to-close")
	}
}

func TestPublic_MenuClose(t *testing.T) {
	p := NewPublic()
	rr := httptest.NewRecorder()
	p.MenuClose(rr, newTestRequest(t, "/partials/menu-close"))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	body := rr.Body.String()

	// MenuClose should include an out-of-band swap to restore the hamburger icon
	if !strings.Contains(body, `hx-swap-oob="true"`) {
		t.Errorf("menu-close body missing hx-swap-oob for button restore")
	}
	if !strings.Contains(body, `id="menu-btn"`) {
		t.Errorf("menu-close body missing id=menu-btn for OOB swap target")
	}

	// The restored button should point back to /partials/menu
	if !strings.Contains(body, `/partials/menu`) {
		t.Errorf("menu-close body missing /partials/menu endpoint in restored button")
	}
}

func TestPublic_Brand(t *testing.T) {
	p := NewPublic()
	rr := httptest.NewRecorder()
	p.Brand(rr, newTestRequest(t, "/brand"))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	ct := rr.Header().Get("Content-Type")
	if !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q, want text/html prefix", ct)
	}
	body := rr.Body.String()

	// Section 1: Color system — base colors
	for _, want := range []string{
		"#0a0a0f", "#141420", "#1e1e2e", "#f0f0f5", "#9090a0",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing base color %q", want)
		}
	}

	// Section 1: Primary accent (amber)
	for _, want := range []string{
		"#f59e0b", "#fbbf24",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing amber accent color %q", want)
		}
	}

	// Section 1: Secondary accent (pink)
	for _, want := range []string{
		"#ec4899", "#f472b6",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing pink accent color %q", want)
		}
	}

	// No violet references should remain
	if strings.Contains(body, "#7c3aed") {
		t.Errorf("body still contains old violet color #7c3aed")
	}

	// Section 2: Usage rules
	if !strings.Contains(body, "Usage Rules") {
		t.Errorf("body missing Usage Rules section")
	}

	// Gradient reference
	if !strings.Contains(body, "linear-gradient") {
		t.Errorf("body missing gradient reference")
	}

	// Must include the color-switcher script (loaded via base layout)
	if !strings.Contains(body, "color-switcher.js") {
		t.Errorf("body missing color-switcher.js script reference")
	}

	// Section 3: Typography
	for _, want := range []string{
		"Space Grotesk", "Inter", "JetBrains Mono",
	} {
		if !strings.Contains(body, want) {
			t.Errorf("body missing font name %q", want)
		}
	}

	// Section 4: Usage examples — buttons, card, link
	if !strings.Contains(body, "bg-accent") {
		t.Errorf("body missing primary button example with bg-accent")
	}
}

// TestPublic_NonGETMethods confirms the handlers accept any method at the
// handler level — route-level method constraints are enforced by Chi, not
// inside the handler. This test pins the contract so a future refactor
// doesn't silently change it.
func TestPublic_HandlerIsMethodAgnostic(t *testing.T) {
	p := NewPublic()
	for _, method := range []string{http.MethodGet, http.MethodPost, http.MethodPut} {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(method, "/", nil)
		if err != nil {
			t.Fatalf("build request: %v", err)
		}
		p.Home(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("method %s: status = %d, want %d", method, rr.Code, http.StatusOK)
		}
	}
}
