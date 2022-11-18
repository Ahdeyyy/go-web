package handlers

import (
	"log"
	"net/http"

	"github.com/Ahdeyyy/go-web/internal/config"
	"github.com/Ahdeyyy/go-web/internal/render"
)

var Dep *Dependency

// Dependency is the dependency for the handlers
type Dependency struct {
	// App is the application configuration
	App *config.Config
}

// NewDependency creates a new dependency for the handlers
func NewDependency(app *config.Config) *Dependency {
	return &Dependency{
		App: app,
	}
}

// Init initializes the handlers
func Init(d *Dependency) {
	Dep = d
}

// Home is the home handler
func (d *Dependency) Home(w http.ResponseWriter, r *http.Request) {

	session, err := d.App.SessionStore.Get(r, "app-session")

	if err != nil {
		log.Println(err)
	}

	session.Values["test"] = "t"

	err = session.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the template from the cache
	render.RenderTemplate(w, r, "home.tmpl")
}
