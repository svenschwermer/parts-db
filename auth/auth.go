package auth

import (
	"bytes"
	"crypto/sha256"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/svenschwermer/parts-db/session"
)

const (
	Path            = "/login"
	sessionValidity = 24 * time.Hour
)

type Handler struct {
	tmpl   *template.Template
	pwHash [sha256.Size]byte
	sm     *session.Manager
}

func New(tmpl *template.Template, password string) *Handler {
	return &Handler{
		tmpl:   tmpl,
		pwHash: sha256.Sum256([]byte(password)),
		sm:     session.NewManager(),
	}
}

func (h *Handler) Required(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("session")
	return err != nil || !h.sm.Valid(cookie.Value)
}

func (h *Handler) RedirectIfRequired(w http.ResponseWriter, req *http.Request) bool {
	if h.Required(w, req) {
		referer := url.QueryEscape(req.URL.String())
		http.Redirect(w, req, Path+"?referer="+referer, http.StatusFound)
		return true
	}
	return false
}

func (h *Handler) Login(w http.ResponseWriter, req *http.Request) {
	var errorString string

	if req.Method == http.MethodPost {
		if err := req.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		hashed := sha256.Sum256([]byte(req.PostForm.Get("password")))
		if bytes.Equal(hashed[:], h.pwHash[:]) {
			redirect, err := url.QueryUnescape(req.FormValue("referer"))
			if err != nil {
				redirect = "/"
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "session",
				Value: h.sm.New(sessionValidity),
			})
			http.Redirect(w, req, redirect, http.StatusFound)
			return
		}
		errorString = "ERROR: Invalid Password"
	}

	err := h.tmpl.ExecuteTemplate(w, "login.html", errorString)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
