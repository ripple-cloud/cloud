package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"time"

	"code.google.com/p/go.crypto/bcrypt"
	"github.com/ripple-cloud/cloud/data"
)

func loginPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "login", nil)
}

func signupPageHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "signup", nil)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	redirectPath := "/"

	login = data.User{
		Email: r.FormValue("email"),
	}

	// Check if user's email exists.
	// Reference: data package.
	if login.GetUserFrom(db).Email == "" {
		user = data.User{
			Username:  r.FormValue("username"),
			Email:     r.FormValue("email"),
			Password:  data.Encrypt(r.FormValue("password")),
			Token:     data.GenerateToken(),
			CreatedAt: time.Now(),
		}
		user.AddTo(db)

		redirectPath = "/home"
	} else {
		fmt.Fprint(w, "email is already taken")
	}

	http.Redirect(w, r, redirectPath, 302)
}

var login data.User
var user data.User

func loginHandler(w http.ResponseWriter, r *http.Request) {
	redirectPath := "/"

	login = data.User{
		Email:    r.FormValue("email"),
		Password: []byte(r.FormValue("password")),
	}

	if login.GetUserFrom(db).Email == login.Email {
		user = *login.GetUserFrom(db)

		err := bcrypt.CompareHashAndPassword(user.Password, []byte(login.Password))
		if err != nil {
			fmt.Println(err)
			fmt.Fprint(w, "wrong password or user does not exist")
			return
		}

		redirectPath = "/home"
	} else {
		// TODO: Better handling of error messages.
		fmt.Fprint(w, "wrong password or user does not exist")
	}

	http.Redirect(w, r, redirectPath, 302)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {

	user = data.User{
		Token: login.GetUserFrom(db).Token,
	}

	renderTemplate(w, "home", &user)
}

var t = make(map[string]*template.Template)

func init() {
	for _, tmpl := range []string{"home", "login", "signup"} {
		path := filepath.Join(cwd, "github.com/ripple-cloud/cloud/templates/"+tmpl+".html")

		t[tmpl] = template.Must(template.New("tmpl").ParseFiles(path))
	}
}

func renderTemplate(w http.ResponseWriter, tmpl string, user *data.User) {
	err := t[tmpl].ExecuteTemplate(w, tmpl+".html", user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
