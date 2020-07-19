package server

import (
	"database/sql"
	"html/template"
	"net/http"
)

type authHandler interface {
	Required(http.ResponseWriter, *http.Request) bool
	RedirectIfRequired(http.ResponseWriter, *http.Request) bool
}

type Server struct {
	*http.ServeMux
	tmpl *template.Template
	auth authHandler
	db   *sql.DB
}

func New(tmpl *template.Template, auth authHandler, db *sql.DB) (*Server, error) {
	s := &Server{http.NewServeMux(), tmpl, auth, db}
	if err := s.createTables(); err != nil {
		return nil, err
	}
	return s, nil
}

type nav struct {
	Label  string
	HRef   string
	Active bool
}

type tmplData struct {
	Title   string
	Error   string
	Info    string
	Nav     []nav
	Content interface{}
}

var defaultNav = []nav{
	{Label: "Parts List", HRef: "/list"},
	{Label: "New Part", HRef: "/new"},
}

func getTmplData(title string, content interface{}) *tmplData {
	d := &tmplData{
		Title:   title,
		Nav:     make([]nav, len(defaultNav)),
		Content: content,
	}
	copy(d.Nav, defaultNav)
	for i := range d.Nav {
		if d.Nav[i].Label == title {
			d.Nav[i].Active = true
			break
		}
	}
	return d
}
