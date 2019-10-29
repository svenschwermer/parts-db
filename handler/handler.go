package handler

import (
	"database/sql"
	"html/template"
	"net/http"
)

const (
	categoriesSchema = `CREATE TABLE IF NOT EXISTS categories (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		parent INTEGER,
		FOREIGN KEY(parent) REFERENCES categories(id));`
	partsSchema = `CREATE TABLE IF NOT EXISTS parts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		pn TEXT,
		manufacturer TEXT,
		category INTEGER,
		value REAL,
		package TEXT,
		description TEXT,
		location TEXT,
		inventory INTEGER,
		FOREIGN KEY(category) REFERENCES categories(id));`
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

func (h *Handler) createTables() error {
	_, err := h.db.Exec(categoriesSchema)
	if err != nil {
		return err
	}
	_, err = h.db.Exec(partsSchema)
	return err
}
