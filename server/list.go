package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/svenschwermer/parts-db/si"
)

type distributor struct {
	Name string
	URL  string
}

type partData struct {
	PartID       string
	PartNumber   string
	Manufacturer string
	Category     string
	Value        string
	UnitPrefix   string
	Unit         string
	Package      string
	Description  string
	Location     string
	Inventory    int32
	Distributors []distributor
}

func (s *Server) List(w http.ResponseWriter, req *http.Request) {
	if s.auth.RedirectIfRequired(w, req) {
		return
	}

	rows, err := s.db.Query(`
		SELECT p.id, pn, manufacturer, c.name, value, package, unit,
			description, location, inventory, d.name, d.url
		FROM parts p
		LEFT JOIN (
			SELECT id, name
			FROM categories
		) c ON p.category = c.id
		LEFT JOIN (
			SELECT part, name, url
			FROM distributors
		) d ON p.id = d.part
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var parts []*partData
	partLookup := make(map[string]*partData)
	for rows.Next() {
		var id, inventory sql.NullInt32
		var pn, mfr, cat, pkg, unit, desc, loc, distiName, distiURL sql.NullString
		val := new(si.Quantity)
		err = rows.Scan(&id, &pn, &mfr, &cat, val, &pkg, &unit, &desc, &loc, &inventory, &distiName, &distiURL)
		if err != nil {
			log.Println(err)
			continue
		}

		idString := fmt.Sprint(id.Int32)
		disti := distributor{
			Name: distiName.String,
			URL:  distiURL.String,
		}
		if p, ok := partLookup[idString]; ok {
			p.Distributors = append(p.Distributors, disti)
			continue
		}

		p := &partData{
			PartID:       idString,
			PartNumber:   pn.String,
			Manufacturer: mfr.String,
			Category:     cat.String,
			Unit:         unit.String,
			Package:      pkg.String,
			Description:  desc.String,
			Location:     loc.String,
			Inventory:    inventory.Int32,
		}
		if val.Valid {
			p.Value, p.UnitPrefix = val.Coeff, val.Prefix
		}
		if distiName.Valid && distiURL.Valid {
			p.Distributors = []distributor{disti}
		}
		parts = append(parts, p)
		partLookup[p.PartID] = p
	}

	tmplData := getTmplData("Parts List", parts)
	err = s.tmpl.ExecuteTemplate(w, "list.html", tmplData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) ChangeInventory(w http.ResponseWriter, req *http.Request) {
	if s.auth.Required(w, req) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	part := req.PostForm.Get("part")
	qty, err := strconv.Atoi(req.PostForm.Get("qty"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = s.db.Exec(`
		UPDATE parts
		SET inventory = ?
		WHERE id = ?
	`, qty, part)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Changing inventory for part (id=%s) to %d", part, qty)
}
