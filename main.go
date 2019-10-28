package main

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/svenschwermer/parts-db/session"
)

var (
	sitePasswordHash = sha256.Sum256([]byte(os.Getenv("SITE_PASSWORD")))
	templates        = template.Must(template.ParseFiles("login.html", "style.css"))
	sessions         = session.NewManager()
)

func needsAuth(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("session")
	if err == nil && sessions.Valid(cookie.Value) {
		return false
	}

	referer := url.QueryEscape(req.URL.String())
	http.Redirect(w, req, "/login?referer="+referer, http.StatusFound)
	return true
}

func home(w http.ResponseWriter, req *http.Request) {
	if needsAuth(w, req) {
		return
	}

	fmt.Fprintf(w, "home")
}

func login(w http.ResponseWriter, req *http.Request) {
	var err string

	if req.Method == http.MethodPost {
		if err := req.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		hashed := sha256.Sum256([]byte(req.PostForm.Get("password")))
		if bytes.Equal(hashed[:], sitePasswordHash[:]) {
			redirect, err := url.QueryUnescape(req.FormValue("referer"))
			if err != nil {
				redirect = "/"
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "session",
				Value: sessions.New(24 * time.Hour),
			})
			http.Redirect(w, req, redirect, http.StatusFound)
			return
		}
		err = "ERROR: Invalid Password"
	}

	templates.ExecuteTemplate(w, "login.html", err)
}

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/foo", home)
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":8080", nil)
	log.Fatal(err)
}
