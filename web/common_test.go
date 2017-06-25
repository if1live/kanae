package web

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFuncMapTemplate(t *testing.T) {
	t.Parallel()

	tempDir, err := ioutil.TempDir("", "kanae-web")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	d := []byte(`{{.Title | upper}}`)
	tp := filepath.Join(tempDir, "simple.html")
	ioutil.WriteFile(tp, d, 0644)

	ctx := map[string]string{
		"Title": "kanae",
	}

	w := new(bytes.Buffer)
	err = renderTemplateCore(w, tp, ctx)
	if err != nil {
		t.Fatalf("renderTemplateCore: %v", err)
	}
	assert.Equal(t, "KANAE", w.String())
}

func TestFuncMapNestedTemplate(t *testing.T) {
	t.Parallel()

	tempDir, err := ioutil.TempDir("", "kanae-web")
	if err != nil {
		t.Fatalf("tempDir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	layout := []byte(`{{ template "content" .}}`)
	lp := filepath.Join(tempDir, "layout.html")
	ioutil.WriteFile(lp, layout, 0644)

	content := []byte(`{{ define "content" }}{{ .Title|upper }}{{ end }}`)
	tp := filepath.Join(tempDir, "content.html")
	ioutil.WriteFile(tp, content, 0644)

	ctx := map[string]string{
		"Title": "kanae",
	}

	w := new(bytes.Buffer)
	err = renderLayoutTemplateCore(w, lp, tp, ctx)
	if err != nil {
		t.Fatalf("renderLayoutTemplateCore: %v", err)
	}
	assert.Equal(t, "KANAE", w.String())
}
