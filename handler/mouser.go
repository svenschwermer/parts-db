package handler

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/svenschwermer/parts-db/config"
)

// Mouser forwards the request to the Mouser API server and returns its result
func (h *Handler) Mouser(w http.ResponseWriter, req *http.Request) {
	if h.auth.Required(w, req) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, req.Body); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pathElems := strings.SplitN(strings.TrimLeft(req.URL.Path, "/"), "/", 2)
	if len(pathElems) >= 1 {
		pathElems = pathElems[1:]
	}
	mouserURL := "https://api.mouser.com" +
		"/api/" + strings.Join(pathElems, "/") +
		"?apiKey=" + url.QueryEscape(config.Env.MouserAPIKey)

	resp, err := http.Post(mouserURL, "application/json", buf)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
