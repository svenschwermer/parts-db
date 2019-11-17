package handler

import (
	"database/sql"
	"html/template"
	"net/http"
)

type authHandler interface {
	Required(http.ResponseWriter, *http.Request) bool
	RedirectIfRequired(http.ResponseWriter, *http.Request) bool
}

type Handler struct {
	*http.ServeMux
	tmpl *template.Template
	auth authHandler
	db   *sql.DB
}

func New(tmpl *template.Template, auth authHandler, db *sql.DB) (*Handler, error) {
	h := &Handler{http.NewServeMux(), tmpl, auth, db}
	if err := h.createTables(); err != nil {
		return nil, err
	}
	return h, nil
}
