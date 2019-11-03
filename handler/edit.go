package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

type newPartPage struct {
	Title string
	partData
	Manufacturers    []string
	Categories       []string
	Locations        []string
	DistributorNames []string
}

func (h *Handler) New(w http.ResponseWriter, req *http.Request) {
	if h.auth.RedirectIfRequired(w, req) {
		return
	}
	contents := newPartPage{Title: "New part"}
	contents.PartID = "new"
	if req.Method == http.MethodPost {
		if err := req.ParseForm(); err != nil {
			fmt.Fprint(w, err)
			return
		}
		tx, err := h.db.Begin()
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		defer tx.Rollback()

		var catID sql.NullInt32
		cat := req.PostForm.Get("category")
		if cat != "" {
			catID, err = h.getCategoryID(tx, cat)
			if err != nil {
				log.Println("/new:", err)
			} else {
				catID.Valid = true
			}
		}

		_, err = tx.Exec(`INSERT INTO parts
		                  (pn,manufacturer,category,value,package,description,location,inventory)
		                  VALUES
		                  (?,?,?,?,?,?,?,?)`,
			req.PostForm.Get("pn"), req.PostForm.Get("manufacturer"),
			catID, req.PostForm.Get("value"),
			req.PostForm.Get("package"), req.PostForm.Get("description"),
			req.PostForm.Get("location"), req.PostForm.Get("inventory"))
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		if err = tx.Commit(); err != nil {
			fmt.Fprint(w, err)
			return
		}
	}
	err := h.tmpl.ExecuteTemplate(w, "edit.html", contents)
	if err != nil {
		fmt.Fprint(w, err)
	}
}

func (h *Handler) Edit(w http.ResponseWriter, req *http.Request) {
	if h.auth.RedirectIfRequired(w, req) {
		return
	}

}

func (h *Handler) getCategoryID(tx *sql.Tx, category string) (id sql.NullInt32, err error) {
	getID, err := tx.Prepare(`SELECT id FROM categories WHERE name = ?`)
	if err != nil {
		return
	}
	err = getID.QueryRow(category).Scan(&id)
	if err == sql.ErrNoRows {
		_, err = tx.Exec(`INSERT INTO categories (name) VALUES (?)`, category)
	}
	if err != nil {
		return
	}
	err = getID.QueryRow(category).Scan(&id)
	return id, err
}
