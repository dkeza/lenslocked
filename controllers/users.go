package controllers

import (
	"fmt"
	"lenslocked/views"
	"net/http"
)

func NewUsers() *Users {
	return &Users{
		NewView: views.NewView("bootstrap", "views/users/new.gohtml"),
	}
}

type Users struct {
	NewView *views.View
}

// New renders new user form
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

// Create new user
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "User created")
}
