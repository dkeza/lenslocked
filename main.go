package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"simplegallery/controllers"
	"simplegallery/email"
	"simplegallery/middleware"
	"simplegallery/models"
	"simplegallery/rand"
	"time"

	llctx "simplegallery/context"

	"golang.org/x/oauth2"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

func main() {
	prod := flag.Bool("prod", false, "Set true in production, to force reading of .config file.")
	flag.Parse()

	cfg := LoadConfig(*prod)
	dbCfg := cfg.Database

	serverPort := cfg.GetPort()

	services, err := models.NewServices(
		models.WithGorm(dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithGallery(),
		models.WithImage(),
		models.WithOAuth(),
	)

	must(err)

	defer services.Close()
	//services.DestructiveReset()
	services.AutoMigrate()

	mgCfg := cfg.Mailgun
	emailer := email.NewClient(
		email.WithSender("SimpleGallery Support", "support@sandboxedbdc3b36f894f5b8edeb7c47e599964.mailgun.org"),
		email.WithMailgun(mgCfg.Domain, mgCfg.APIKey, mgCfg.PublicAPIKey),
	)

	r := mux.NewRouter()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(services.User, emailer)
	galleriesC := controllers.NewGalleries(services.Gallery, services.Image, r)

	isProd := false
	b, err := rand.Bytes(32)
	must(err)
	csrfMw := csrf.Protect(b, csrf.Secure(isProd))

	userMw := middleware.User{
		UserService: services.User,
	}
	requireUserMw := middleware.RequireUser{
		User: userMw,
	}

	dbxOAuth := &oauth2.Config{
		ClientID:     cfg.Dropbox.ID,
		ClientSecret: cfg.Dropbox.Secret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cfg.Dropbox.AuthURL,
			TokenURL: cfg.Dropbox.TokenURL,
		},
		RedirectURL: "http://localhost:3333/oauth/dropbox/callback",
	}

	dbxRedirect := func(w http.ResponseWriter, r *http.Request) {
		state := csrf.Token(r)
		cookie := http.Cookie{
			Name:     "oauth_state",
			Value:    state,
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		url := dbxOAuth.AuthCodeURL(state)
		http.Redirect(w, r, url, http.StatusFound)
	}
	r.HandleFunc("/oauth/dropbox/connect", requireUserMw.ApplyFn(dbxRedirect))
	dbxCallback := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		state := r.FormValue("state")
		cookie, err := r.Cookie("oauth_state")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} else if cookie == nil || cookie.Value != state {
			http.Error(w, "Invalid state provided", http.StatusBadRequest)
			return
		}
		cookie.Value = ""
		cookie.Expires = time.Now()
		http.SetCookie(w, cookie)
		code := r.FormValue("code")
		token, err := dbxOAuth.Exchange(context.TODO(), code)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user := llctx.User(r.Context())
		existing, err := services.OAuth.Find(user.ID, models.OAuthDropbox)
		if err == models.ErrNotFound {
			// noop
		} else if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		} else {
			services.OAuth.Delete(existing.ID)
		}
		userOAuth := models.OAuth{
			UserID:  user.ID,
			Token:   *token,
			Service: models.OAuthDropbox,
		}
		err = services.OAuth.Create(&userOAuth)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "%+v", token)
	}
	r.HandleFunc("/oauth/dropbox/callback", requireUserMw.ApplyFn(dbxCallback))

	dbxQuery := func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		path := r.FormValue("path")

		user := llctx.User(r.Context())
		userOAuth, err := services.OAuth.Find(user.ID, models.OAuthDropbox)
		if err != nil {
			panic(err)
		}
		token := userOAuth.Token

		data := struct {
			Path string `json:"path"`
		}{
			Path: path,
		}
		dataBytes, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		client := dbxOAuth.Client(context.TODO(), &token)
		req, err := http.NewRequest(http.MethodPost, "https://api.dropboxapi.com/2/files/list_folder", bytes.NewReader(dataBytes))
		if err != nil {
			panic(err)
		}
		req.Header.Add("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()
		io.Copy(w, resp.Body)
	}
	r.HandleFunc("/oauth/dropbox/test", requireUserMw.ApplyFn(dbxQuery))

	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	r.HandleFunc("/login", usersC.Show).Methods("GET")
	r.HandleFunc("/login", usersC.Login).Methods("POST")
	r.HandleFunc("/logout", requireUserMw.ApplyFn(usersC.Logout)).Methods("POST")
	r.Handle("/forgot", usersC.ForgotPwView).Methods("GET")
	r.HandleFunc("/forgot", usersC.InitiateReset).Methods("POST")
	r.HandleFunc("/reset", usersC.ResetPw).Methods("GET")
	r.HandleFunc("/reset", usersC.CompleteReset).Methods("POST")

	// Assets
	assetHandler := http.FileServer(http.Dir("./assets/"))
	assetHandler = http.StripPrefix("/assets/", assetHandler)
	r.PathPrefix("/assets/").Handler(assetHandler)

	// Image routes
	imageHandler := http.FileServer(http.Dir("./images/"))
	r.PathPrefix("/images/").Handler(http.StripPrefix("/images/", imageHandler))

	// Gallery routes
	r.Handle("/galleries", requireUserMw.ApplyFn(galleriesC.Index)).Methods("Get").Name("index_gallery")
	r.Handle("/galleries/new", requireUserMw.Apply(galleriesC.New)).Methods("Get")
	r.HandleFunc("/galleries", requireUserMw.ApplyFn(galleriesC.Create)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/edit", requireUserMw.ApplyFn(galleriesC.Edit)).Methods("GET").Name("edit_gallery")
	r.HandleFunc("/galleries/{id:[0-9]+}/update", requireUserMw.ApplyFn(galleriesC.Update)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/delete", requireUserMw.ApplyFn(galleriesC.Delete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images", requireUserMw.ApplyFn(galleriesC.ImageUpload)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}/images/{filename}/delete", requireUserMw.ApplyFn(galleriesC.ImageDelete)).Methods("POST")
	r.HandleFunc("/galleries/{id:[0-9]+}", galleriesC.Show).Methods("GET").Name("show_gallery")

	fmt.Println("Server is listening on port " + serverPort)
	http.ListenAndServe(":"+serverPort, csrfMw(userMw.Apply(r)))
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
