package handler

import (
	"encoding/json"
	"net/http"

	"github.com/svenschwermer/parts-db/mouser"
)

// Mouser issues a request to Mouser's search API and returns the details in
// JSON format
func (h *Handler) Mouser(w http.ResponseWriter, req *http.Request) {
	if h.auth.Required(w, req) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	part, err := mouser.GetPart(req.PostFormValue("pn"))
	if err != nil {
		// TODO: error handling
		return
	}

	json.NewEncoder(w).Encode(part)
}
