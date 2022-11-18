package render

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Ahdeyyy/go-web/internal/config"
	"github.com/Ahdeyyy/go-web/internal/models"
)

var functions = template.FuncMap{}

var appConfig *config.Config
var pagesPath string = "./web/templates/pages"
var layoutsPath string = "./web/templates/layouts"

// NewTemplates sets the config for the template package
func NewTemplates(a *config.Config) {
	appConfig = a
}

func AddDefaultContext(td *models.TemplateContext, r *http.Request) *models.TemplateContext {

	session := appConfig.Session

	if flash := session.Flashes("flash"); len(flash) > 0 {
		td.Flash = flash[0].(string)
	}
	if Error := session.Flashes("error"); len(Error) > 0 {
		td.Error = Error[0].(string)
	}

	if Warning := session.Flashes("warning"); len(Warning) > 0 {
		td.Warning = Warning[0].(string)
	}

	if session.Values["user_id"] != nil {
		td.IsAuthenticated = 1
	}

	return td
}

// RenderTemplate renders template
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, tContext *models.TemplateContext) {

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
	tContext = AddDefaultContext(tContext, r)
	appConfig.Session.Save(r, w)

	err := t.Execute(buf, tContext)

	if err != nil {
		log.Println(err)
		return
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
