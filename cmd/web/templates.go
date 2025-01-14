package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"text/template"

	"github.com/Lanrey-waju/sunny-akins/ui"
)

type templateData struct {
	CurrentYear int
	Form        any
	Flash       string
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		fmt.Println(name)
		cache[name] = ts
	}

	return cache, nil
}
