package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/Ahdeyyy/go-web/internal/config"
	"github.com/Ahdeyyy/go-web/internal/handlers"
	"github.com/Ahdeyyy/go-web/internal/routes"
)

// config is the configuration for the application
var appConfig config.Config

const portNumber = ":8080"

// main is the entry point for the application
func main() {

	// create a new server
	srv := &http.Server{
		Addr:    appConfig.Port,
		Handler: routes.Routes(&appConfig),
	}

	// start the server
	go func() {
		log.Println("Starting server on port", appConfig.Port)
		err := srv.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()

	// wait for a signal to shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutting down server...")

}

func init() {
	// set the configuration
	appConfig.Port = portNumber

	// initialize the handlers
	handlers.Init(handlers.NewDependency(&appConfig))
}
