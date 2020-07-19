package server

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/svenschwermer/parts-db/si"
)

func (s *Server) Edit(w http.ResponseWriter, req *http.Request) {
	if s.auth.RedirectIfRequired(w, req) {
		return
	}

	if err := req.ParseForm(); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	switch req.Method {
	case http.MethodPost:
		w.WriteHeader(http.StatusInternalServerError)
	case http.MethodGet:
		if err := s.getEditPage(w, req); err != nil {
			log.Print(err)
		}
	default:
		http.Error(w, "unexpected method", http.StatusBadRequest)
	}
}

func (s *Server) getEditPage(w http.ResponseWriter, req *http.Request) error {
	p := new(newPartPage)
	p.PartID = req.Form.Get("part")
	tmplData := getTmplData("Edit Part", p)

	var id, inventory sql.NullInt32
	var pn, mfr, cat, pkg, unit, desc, loc sql.NullString
	val := new(si.Quantity)

	err := s.db.QueryRow(`
		SELECT p.id, pn, manufacturer, c.name, value, package, unit,
			description, location, inventory
		FROM parts p
		LEFT JOIN (
			SELECT id, name
			FROM categories
		) c ON p.category = c.id
		WHERE p.id = ?
	`, p.PartID).Scan(&id, &pn, &mfr, &cat, val, &pkg, &unit, &desc, &loc, &inventory)
	if err != nil {
		return err
	}

	p.PartNumber = pn.String
	p.Manufacturer = mfr.String
	p.Category = cat.String
	p.Unit = unit.String
	p.Package = pkg.String
	p.Description = desc.String
	p.Location = loc.String
	p.Inventory = inventory.Int32
	if val.Valid {
		p.Value, p.UnitPrefix = val.Coeff, val.Prefix
	}

	rows, err := s.db.Query(`
	  SELECT name, url
	  FROM distributors
		WHERE part = ?
	`, p.PartID)
	for rows.Next() {
		var name, url sql.NullString
		if err := rows.Scan(&name, &url); err != nil {
			return err
		}
		if name.Valid && url.Valid {
			p.Distributors = append(p.Distributors, distributor{Name: name.String, URL: url.String})
		}
	}

	return s.tmpl.ExecuteTemplate(w, "edit.html", tmplData)
}
