package handler

import "net/http"

func (h *Handler) Edit(w http.ResponseWriter, req *http.Request) {
	if h.auth.RedirectIfRequired(w, req) {
		return
	}

}
