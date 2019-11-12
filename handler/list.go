package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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

func (h *Handler) List(w http.ResponseWriter, req *http.Request) {
	if h.auth.RedirectIfRequired(w, req) {
		return
	}

	rows, err := h.db.Query(`
		SELECT p.id, pn, manufacturer, c.name, value, package, description, location, inventory
		FROM parts p
		LEFT JOIN (
			SELECT id, name
			FROM categories
		) c ON p.category = c.id
	`)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var parts []partData
	for rows.Next() {
		var id, inventory sql.NullInt32
		var pn, mfr, cat, val, pkg, desc, loc sql.NullString
		if err = rows.Scan(&id, &pn, &mfr, &cat, &val, &pkg, &desc, &loc, &inventory); err != nil {
			log.Println(err)
		} else {
			parts = append(parts, partData{
				PartID:       fmt.Sprint(id.Int32),
				PartNumber:   pn.String,
				Manufacturer: mfr.String,
				Category:     cat.String,
				Value:        val.String,
				UnitPrefix:   "",
				Unit:         "",
				Package:      pkg.String,
				Description:  desc.String,
				Location:     loc.String,
				Inventory:    inventory.Int32,
			})
		}
	}

	err = h.tmpl.ExecuteTemplate(w, "list.html", parts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type changeInventoryReq struct {
	Part  string
	Delta int
}

func (h *Handler) ChangeInventory(w http.ResponseWriter, req *http.Request) {
	if h.auth.Required(w, req) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if err := req.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(req.Body)
	reqContent := new(changeInventoryReq)
	if err := decoder.Decode(reqContent); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, err := h.db.Exec(`
		UPDATE parts
		SET inventory = inventory + ?
		WHERE id = ?
		  AND (inventory + ?) >= 0
	`, reqContent.Delta, reqContent.Part, reqContent.Delta)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var newInventory int
	err = h.db.QueryRow(`
		SELECT inventory
		FROM parts
		WHERE id = ?
	`, reqContent.Part).Scan(&newInventory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Changing inventory for part (id=%v) to %d", reqContent.Part, newInventory)
	w.Write([]byte(fmt.Sprint(newInventory)))
}
