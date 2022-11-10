package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Ahdeyyy/go-web/internal/config"
)

var functions = template.FuncMap{}

var appConfig *config.Config
var pagesPath string = "./web/templates/pages"
var layoutsPath string = "./web/templates/layouts"

// NewTemplates sets the config for the template package
func NewTemplates(a *config.Config) {
	appConfig = a
}

// RenderTemplate renders template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string) {

	var tc map[string]*template.Template

	if appConfig.Debug {
		tc, _ = CreateTemplateCache()
	} else {
		tc = appConfig.TemplateCache
	}

	t, ok := tc[tmpl]

	if !ok {
		log.Fatal("Could not get template from template cache")
	}

	buf := new(bytes.Buffer)

	err := t.Execute(buf, nil)

	if err != nil {
		log.Println(err)
	}

	_, err = buf.WriteTo(w)

	if err != nil {
		log.Println(err)
	}

}

// CreateTemplateCache creates a template cache as a map
func CreateTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.tmpl", pagesPath))

	if err != nil {
		return cache, err
	}

	for _, page := range pages {
		name := filepath.Base(page)

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)

		if err != nil {
			return cache, err
		}

		matches, err := filepath.Glob(fmt.Sprintf("%s/*.tmpl", layoutsPath))

		if err != nil {
			return cache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.tmpl", layoutsPath))

			if err != nil {
				return cache, err
			}
		}

		cache[name] = ts
	}

	return cache, nil

}
