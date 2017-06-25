package web

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"

	"errors"

	"github.com/if1live/kanae/kanaelib"
)

func renderErrorJSON(w http.ResponseWriter, err error, errcode int) {
	type Response struct {
		Error string `json:"error"`
	}
	resp := Response{
		Error: err.Error(),
	}
	data, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(errcode)
	w.Write(data)
}

func renderJSON(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
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
		renderErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, ctx); err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
	}
}

func renderLayoutTemplate(w http.ResponseWriter, tplfile string, v interface{}) {
	basePath := kanaelib.GetExecutablePath()
	lp := path.Join(basePath, "web", "templates", "layout.html")
	fp := path.Join(basePath, "web", "templates", tplfile)

	tpl, err := template.ParseFiles(lp, fp)
	if err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, v); err != nil {
		renderErrorJSON(w, err, http.StatusInternalServerError)
	}
}

func renderStatic(w http.ResponseWriter, r *http.Request, target string) {
	cleaned := path.Clean(target)
	basePath := kanaelib.GetExecutablePath()
	fp := path.Join(basePath, "web", "static", cleaned)
	cleanedFp := path.Clean(fp)
	http.ServeFile(w, r, cleanedFp)
}

func checkPostRequest(w http.ResponseWriter, r *http.Request) bool {
	type Response struct {
		Error string `json:"error"`
	}
	if r.Method != "POST" {
		renderErrorJSON(w, errors.New("only post allowed"), http.StatusBadRequest)
		return false
	}
	return true
}
