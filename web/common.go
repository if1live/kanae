package web

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"

	"github.com/if1live/kanae/kanaelib"
)

func renderJSON(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func renderTemplate(w http.ResponseWriter, tplfile string, ctx interface{}) {
	basePath := kanaelib.GetExecutablePath()
	fp := path.Join(basePath, "web", "templates", tplfile)
	tpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, ctx); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderStatic(w http.ResponseWriter, r *http.Request, target string) {
	cleaned := path.Clean(target)
	basePath := kanaelib.GetExecutablePath()
	fp := path.Join(basePath, "web", "static", cleaned)
	cleanedFp := path.Clean(fp)
	http.ServeFile(w, r, cleanedFp)
}
