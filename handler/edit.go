package handler

import (
	"database/sql"
	"net/http"
)

type newPartPage struct {
	Title string
	Error string
	Info  string
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
		if err := h.postNew(w, req); err != nil {
			contents.Error = err.Error()
			// TODO: fill form with entered data
		} else {
			contents.Info = "Part added"
		}
	}
	h.populateLists(&contents)
	err := h.tmpl.ExecuteTemplate(w, "edit.html", contents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h *Handler) populateLists(contents *newPartPage) {
	h.populateList("manufacturer", &contents.Manufacturers)
	h.populateList("category", &contents.Categories)
	h.populateList("location", &contents.Locations)
}

func (h *Handler) populateList(col string, l *[]string) {
	rows, err := h.db.Query("SELECT DISTINCT " + col + " FROM parts")
	if err == nil {
		for rows.Next() {
			var v string
			if err = rows.Scan(&v); err == nil {
				*l = append(*l, v)
			}
		}
	}
}

func (h *Handler) postNew(w http.ResponseWriter, req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	tx, err := h.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var catID sql.NullInt32
	cat := req.PostForm.Get("category")
	if cat != "" {
		catID, err = h.getCategoryID(tx, cat)
		if err != nil {
			return err
		}
		catID.Valid = true
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
		return err
	}
	return tx.Commit()
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
