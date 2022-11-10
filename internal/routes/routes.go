package routes

import (
	"net/http"

	"github.com/Ahdeyyy/go-web/internal/config"
	"github.com/Ahdeyyy/go-web/internal/handlers"
	"github.com/gorilla/mux"
)

func Routes(app *config.Config) http.Handler {
	mux := mux.NewRouter()

	// serve static files
	fileServer := http.FileServer(http.Dir("./web/static/"))
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fileServer))

	// set the routes
	mux.HandleFunc("/", handlers.Dep.Home)

	return mux

}
