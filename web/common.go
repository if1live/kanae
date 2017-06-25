package web

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path"
	"strings"

	"errors"

	"io"

	"path/filepath"

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

func makeFuncMap() template.FuncMap {
	fmap := template.FuncMap{
		"title":    strings.Title,
		"upper":    strings.ToUpper,
		"floatstr": kanaelib.ToFloatStr,
	}
	return fmap
}

func renderTemplate(w io.Writer, tplfile string, ctx interface{}) error {
	basePath := kanaelib.GetExecutablePath()
	fp := path.Join(basePath, "web", "templates", tplfile)
	return renderTemplateCore(w, fp, ctx)
}

func renderLayoutTemplate(w io.Writer, layoutfile, tplfile string, v interface{}) error {
	basePath := kanaelib.GetExecutablePath()
	lp := path.Join(basePath, "web", "templates", layoutfile)
	fp := path.Join(basePath, "web", "templates", tplfile)
	return renderLayoutTemplateCore(w, lp, fp, v)
}

func renderTemplateCore(w io.Writer, fp string, v interface{}) error {
	_, tplfile := filepath.Split(fp)
	fmap := makeFuncMap()
	tpl, err := template.New(tplfile).Funcs(fmap).ParseFiles(fp)
	if err != nil {
		return err
	}
	return tpl.Execute(w, v)
}
func renderLayoutTemplateCore(w io.Writer, lp, fp string, v interface{}) error {
	_, layoutfile := filepath.Split(lp)
	fmap := makeFuncMap()
	tpl, err := template.New(layoutfile).Funcs(fmap).ParseFiles(lp, fp)
	if err != nil {
		return err
	}
	return tpl.Execute(w, v)
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
