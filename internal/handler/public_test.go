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

func TestPublic_Contact(t *testing.T) {
	p := NewPublic()
	rr := httptest.NewRecorder()
	p.Contact(rr, newTestRequest(t, "/contact"))

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if !strings.Contains(rr.Body.String(), "Contact") {
		t.Errorf("contact body missing 'Contact' text")
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
