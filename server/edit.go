package server

import "net/http"

func (s *Server) Edit(w http.ResponseWriter, req *http.Request) {
	if s.auth.RedirectIfRequired(w, req) {
		return
	}

}
