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
