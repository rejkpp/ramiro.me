// Package handler holds HTTP handlers for the public, admin, and members areas.
package handler

import (
	"net/http"

	"github.com/rejkpp/ramiro.me/templates/layout"
	"github.com/rejkpp/ramiro.me/templates/pages"
)

// Public holds dependencies for the public-site handlers.
type Public struct{}

// NewPublic constructs a Public handler group.
func NewPublic() *Public {
	return &Public{}
}

// Home renders the public landing page.
func (p *Public) Home(w http.ResponseWriter, r *http.Request) {
	render(w, r, pages.Home())
}

// About renders the about page.
func (p *Public) About(w http.ResponseWriter, r *http.Request) {
	render(w, r, pages.About())
}

// Projects renders the projects index.
func (p *Public) Projects(w http.ResponseWriter, r *http.Request) {
	render(w, r, pages.Projects())
}

// Booking renders the booking page.
func (p *Public) Booking(w http.ResponseWriter, r *http.Request) {
	render(w, r, pages.Booking())
}

// Brand renders the brand playbook page.
func (p *Public) Brand(w http.ResponseWriter, r *http.Request) {
	render(w, r, pages.Brand())
}

// Menu renders the navigation menu partial.
func (p *Public) Menu(w http.ResponseWriter, r *http.Request) {
	render(w, r, layout.Menu())
}

// MenuClose renders an empty partial to close the menu.
func (p *Public) MenuClose(w http.ResponseWriter, r *http.Request) {
	render(w, r, layout.MenuClose())
}
