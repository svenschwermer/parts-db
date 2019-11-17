package auth

import (
	"bytes"
	"crypto/sha256"
	"html/template"
	"net/http"
	"net/url"
	"time"

	"github.com/svenschwermer/parts-db/config"
	"github.com/svenschwermer/parts-db/session"
)

const (
	Path            = "/login"
	sessionValidity = 24 * time.Hour
)

type Server struct {
	tmpl   *template.Template
	pwHash [sha256.Size]byte
	sm     *session.Manager
}

func New(tmpl *template.Template) *Server {
	return &Server{
		tmpl:   tmpl,
		pwHash: sha256.Sum256([]byte(config.Env.SitePassword)),
		sm:     session.NewManager(),
	}
}

func (s *Server) Required(w http.ResponseWriter, req *http.Request) bool {
	cookie, err := req.Cookie("session")
	return err != nil || !s.sm.Valid(cookie.Value)
}

func (s *Server) RedirectIfRequired(w http.ResponseWriter, req *http.Request) bool {
	if s.Required(w, req) {
		referer := url.QueryEscape(req.URL.String())
		http.Redirect(w, req, Path+"?referer="+referer, http.StatusFound)
		return true
	}
	return false
}

func (s *Server) Login(w http.ResponseWriter, req *http.Request) {
	var errorString string

	if req.Method == http.MethodPost {
		if err := req.ParseForm(); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		hashed := sha256.Sum256([]byte(req.PostForm.Get("password")))
		if bytes.Equal(hashed[:], s.pwHash[:]) {
			redirect, err := url.QueryUnescape(req.FormValue("referer"))
			if err != nil {
				redirect = "/"
			}
			http.SetCookie(w, &http.Cookie{
				Name:  "session",
				Value: s.sm.New(sessionValidity),
			})
			http.Redirect(w, req, redirect, http.StatusFound)
			return
		}
		errorString = "ERROR: Invalid Password"
	}

	err := s.tmpl.ExecuteTemplate(w, "login.html", errorString)
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
