package main

import (
	"bytes"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"
)

func (app *application) serverError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	if err := app.config.ErrorLog.Output(2, trace); err != nil {
		app.config.ErrorLog.Fatal(err)
	}

	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

// func (app *application) notFound(w http.ResponseWriter) {
// 	app.clientError(w, http.StatusNotFound)
// }

func (app *application) newTemplateData(r *http.Request) *templateData {
	return &templateData{
		CurrentYear: time.Now().Year(),
	}
}

// render is a helper method that renders the templates from the cache
func (app *application) render(w http.ResponseWriter, status int, page string, data *templateData) {
	// Retrieve appropriate template set from the cache based on the page name (like  'index.html).
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	buf := new(bytes.Buffer)

	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, err)
		return
	}

	w.WriteHeader(status)

	if _, err := buf.WriteTo(w); err != nil {
		app.serverError(w, err)
	}
}

func (app *application) renderPartial(
	w http.ResponseWriter,
	page string,
	templateString string,
	data *templateData,
) {
	ts, ok := app.templateCache[page]
	if !ok {
		err := fmt.Errorf("the template %s does not exist", page)
		app.serverError(w, err)
		return
	}

	err := ts.ExecuteTemplate(w, templateString, data)
	if err != nil {
		app.serverError(w, err)
	}
}

func (app *application) background(fn func()) {
	app.wg.Add(1)

	go func() {
		// recover any background panic
		defer func() {
			defer app.wg.Done()

			if err := recover(); err != nil {
				app.config.ErrorLog.Print(err)
			}
		}()
		// run the arbitrary background task
		fn()
	}()
}
