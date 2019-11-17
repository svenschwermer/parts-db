package server

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	"github.com/svenschwermer/parts-db/si"
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

func (s *Server) New(w http.ResponseWriter, req *http.Request) {
	if s.auth.RedirectIfRequired(w, req) {
		return
	}
	contents := newPartPage{Title: "New part"}
	contents.PartID = "new"
	if req.Method == http.MethodPost {
		if err := s.postNew(w, req); err != nil {
			contents.Error = err.Error()
			// TODO: fill form with entered data
		} else {
			contents.Info = "Part added"
		}
	}
	s.populateLists(&contents)
	err := s.tmpl.ExecuteTemplate(w, "edit.html", contents)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) populateLists(contents *newPartPage) {
	s.populateList("manufacturer", &contents.Manufacturers)
	s.populateList("category", &contents.Categories)
	s.populateList("location", &contents.Locations)
}

func (s *Server) populateList(col string, l *[]string) {
	rows, err := s.db.Query("SELECT DISTINCT " + col + " FROM parts")
	if err == nil {
		for rows.Next() {
			var v string
			if err = rows.Scan(&v); err == nil {
				*l = append(*l, v)
			}
		}
	}
}

func (s *Server) postNew(w http.ResponseWriter, req *http.Request) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var catID sql.NullInt32
	cat := req.PostForm.Get("category")
	if cat != "" {
		catID, err = getCategoryID(tx, cat)
		if err != nil {
			return err
		}
		catID.Valid = true
	}

	var mag, unit sql.NullString
	valUnit := req.PostForm.Get("value")
	if strings.TrimSpace(valUnit) != "" {
		val, err := si.Parse(valUnit)
		if err != nil {
			return fmt.Errorf("failed to parse value: %v", err)
		}
		mag.String, mag.Valid = val.Mag.String(), true
		unit.String, unit.Valid = val.Unit, true
	}

	_, err = tx.Exec(`
		INSERT INTO parts
		(pn,manufacturer,category,value,unit,package,description,location,inventory)
		VALUES
		(?,?,?,?,?,?,?,?,?)`, getPostString(req, "pn"),
		getPostString(req, "manufacturer"), catID, mag, unit,
		getPostString(req, "package"), getPostString(req, "description"),
		getPostString(req, "location"), req.PostForm.Get("inventory"))
	if err != nil {
		return err
	}
	return tx.Commit()
}

func getCategoryID(tx *sql.Tx, category string) (id sql.NullInt32, err error) {
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

func getPostString(req *http.Request, name string) sql.NullString {
	s := sql.NullString{String: req.PostForm.Get("pn")}
	s.Valid = (s.String != "")
	return s
}
