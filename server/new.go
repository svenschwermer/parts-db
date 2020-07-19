package server

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/svenschwermer/parts-db/si"
)

var (
	errPartNumberRequired = errors.New("part number required")
)

type newPartPage struct {
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
	contents := new(newPartPage)
	contents.PartID = "new"
	tmplData := getTmplData("New Part", contents)
	if req.Method == http.MethodPost {
		if err := s.postNew(w, req); err != nil {
			tmplData.Error = err.Error()
			// TODO: fill form with entered data
		} else {
			tmplData.Info = "Part added"
		}
	}
	s.populateLists(contents)
	err := s.tmpl.ExecuteTemplate(w, "edit.html", tmplData)
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

	pn := getPostString(req, "pn")
	if !pn.Valid {
		return errPartNumberRequired
	}

	var catID sql.NullInt64
	cat := req.PostForm.Get("category")
	if cat != "" {
		catID, err = getCategoryID(tx, cat)
		if err != nil {
			return err
		}
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

	result, err := tx.Exec(`
		INSERT INTO parts
		(pn,manufacturer,category,value,unit,package,description,location,inventory)
		VALUES
		(?,?,?,?,?,?,?,?,?)`, pn,
		getPostString(req, "manufacturer"), catID, mag, unit,
		getPostString(req, "package"), getPostString(req, "description"),
		getPostString(req, "location"), req.PostForm.Get("inventory"))
	if err != nil {
		return err
	}
	partID, err := result.LastInsertId()
	if err != nil {
		return err
	}

	err = addDistis(tx, partID, req.PostForm)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func getCategoryID(tx *sql.Tx, category string) (id sql.NullInt64, err error) {
	err = tx.QueryRow("SELECT id FROM categories WHERE name = ?", category).Scan(&id)
	var result sql.Result
	if err == sql.ErrNoRows {
		result, err = tx.Exec("INSERT INTO categories (name) VALUES (?)", category)
		if err != nil {
			return
		}
		id.Int64, err = result.LastInsertId()
	}
	if err != nil {
		return
	}
	id.Valid = true
	return
}

func getPostString(req *http.Request, name string) sql.NullString {
	s := sql.NullString{String: req.PostForm.Get(name)}
	s.Valid = (s.String != "")
	return s
}

func addDistis(tx *sql.Tx, partID int64, postForm url.Values) error {
	for k, v := range postForm {
		if strings.HasPrefix(k, "disti_name_") {
			var id int
			n, err := fmt.Sscanf(k, "disti_name_%d", &id)
			if err != nil {
				return fmt.Errorf("invalid distributor name parameter: %s", err)
			}
			if n != 1 {
				return fmt.Errorf("invalid distributor name parameter: n=%d", n)
			}

			url := postForm.Get(fmt.Sprintf("disti_url_%d", id))
			_, err = tx.Exec(`
				INSERT INTO distributors (part,name,url)
				VALUES (?,?,?)`, partID, v[0], url)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
