package handlers

import (
	"net/http"

	"github.com/Ahdeyyy/go-web/internal/config"
	"github.com/Ahdeyyy/go-web/internal/database"
	"github.com/Ahdeyyy/go-web/internal/database/sqlite"
	"github.com/Ahdeyyy/go-web/internal/driver"
	"github.com/Ahdeyyy/go-web/internal/models"
	"github.com/Ahdeyyy/go-web/internal/render"
)

var Dep *Dependency

// Dependency is the dependency for the handlers
type Dependency struct {
	// App is the application configuration
	App *config.Config
	DB  database.DBInterface
}

// NewDependency creates a new dependency for the handlers
func NewDependency(app *config.Config, db *driver.DB) *Dependency {
	return &Dependency{
		App: app,
		DB:  sqlite.NewSqliteInterface(db.SQL, app),
	}
}

// Init initializes the handlers
func Init(d *Dependency) {
	Dep = d
}

// Home is the home handler
func (d *Dependency) Home(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	// Get the user from the database
	user, _ := d.DB.GetUserByID(1)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	data := make(map[string]interface{})
	data["user"] = user

	users, _ := d.DB.AllUsers()
	data["users"] = users

	// Get the template from the cache
	render.RenderTemplate(w, r, "home.tmpl", &models.TemplateContext{
		Data: data,
	})
}

// User is the user handler it creates new users to the database
func (d *Dependency) User(w http.ResponseWriter, r *http.Request) {

	// Create a new user
	user := models.User{
		Username:  "ahdeyyy",
		Firstname: "Ahmed",
		Lastname:  "Eldeeb",
		Email:     "dewd@mail.com",
		Password:  "123456",
	}

	// Insert the user to the database
	i, err := d.DB.CreateUser(user)

	d.App.InfoLog.Println(i, err)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)

	// }

	data := make(map[string]interface{})
	data["user"] = user

	// redirect to the home page
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
