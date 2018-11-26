package controllers

import (
	"fmt"
	"lenslocked/models"
	"lenslocked/views"
	"net/http"
)

func NewUsers(us *models.UserService) *Users {
	return &Users{
		NewView:   views.NewView("bootstrap", "users/new"),
		LoginView: views.NewView("bootstrap", "users/login"),
		us:        us,
	}
}

// Users
type Users struct {
	NewView   *views.View
	LoginView *views.View
	us        *models.UserService
}

// New renders new user form
// GET /signup
func (u *Users) New(w http.ResponseWriter, r *http.Request) {
	u.NewView.Render(w, nil)
}

type SignupForm struct {
	Name     string `schema:"name"`
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Create new user
// POST /signup
func (u *Users) Create(w http.ResponseWriter, r *http.Request) {
	var form SignupForm

	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	user := models.User{
		Name:     form.Name,
		Email:    form.Email,
		Password: form.Password,
	}
	if err := u.us.Create(&user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, user)
}

// Show shows login form
// GET /login
func (u *Users) Show(w http.ResponseWriter, r *http.Request) {
	u.LoginView.Render(w, nil)
}

type LoginForm struct {
	Email    string `schema:"email"`
	Password string `schema:"password"`
}

// Login user
// POST /login
func (u *Users) Login(w http.ResponseWriter, r *http.Request) {
	var form LoginForm

	if err := parseForm(r, &form); err != nil {
		panic(err)
	}

	fmt.Fprintln(w, form)
}
