package handler

import (
	"disapp/internal/view/home"
	"net/http"
)

type homeHandler struct{}

func (h homeHandler) handleIndex(w http.ResponseWriter, r *http.Request) error {
	return home.Index().Render(r.Context(), w)
}
