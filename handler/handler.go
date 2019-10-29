package handler

import (
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
}

func New(tmpl *template.Template, auth authHandler) *Handler {
	h := &Handler{http.NewServeMux(), tmpl, auth}
	return h
}
