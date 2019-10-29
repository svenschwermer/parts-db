package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type distributor struct {
	Name string
	URL  string
}

type partListRow struct {
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
	Inventory    int
	Distributors []distributor
}

func (h *Handler) List(w http.ResponseWriter, req *http.Request) {
	if h.auth.RedirectIfRequired(w, req) {
		return
	}
	err := h.tmpl.ExecuteTemplate(w, "list.html", []partListRow{
		{
			PartID:       "12",
			PartNumber:   "CRCW020141K2FNED",
			Manufacturer: "Vishay",
			Distributors: []distributor{
				{
					Name: "Mouser",
				},
				{
					Name: "Digi-Key",
				},
			},
		},
	})
	if err != nil {
		w.Write([]byte(err.Error()))
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
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(req.Body)
	reqContent := new(changeInventoryReq)
	if err := decoder.Decode(reqContent); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.Write([]byte(fmt.Sprint(reqContent)))
}
