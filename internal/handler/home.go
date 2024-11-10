package handler

import (
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"
)

type homeHandler struct{}

func (h homeHandler) handleIndex(w http.ResponseWriter, r *http.Request) error {
	name := "index.html"
	path := filepath.Join("./", "web", name)
	htmlTemp, err := template.New(name).ParseFiles(path)
	if err != nil {
		slog.Error("faild to pars html template", slog.Attr{
			Key:   "error",
			Value: slog.StringValue(err.Error()),
		})
	}
	htmlTemp.Execute(w, nil)

	return nil
}
