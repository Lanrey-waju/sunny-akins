package main

import (
	"fmt"
	"path/filepath"
	"text/template"
)

type templateData struct {
	CurrentYear int
	Form        any
	Flash       string
}

func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob("../../ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}
	for _, page := range pages {
		name := filepath.Base(page)

		files := []string{
			"../../ui/html/base.html",
			page,
			"../../ui/html/partials/contact_form.html",
			"../../ui/html/partials/flash.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}

		fmt.Println(name)
		cache[name] = ts
	}

	return cache, nil
}
