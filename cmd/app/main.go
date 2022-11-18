package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/Ahdeyyy/go-web/internal/config"
	"github.com/Ahdeyyy/go-web/internal/handlers"
	"github.com/Ahdeyyy/go-web/internal/render"
	"github.com/Ahdeyyy/go-web/internal/routes"
)

// config is the configuration for the application
var appConfig config.Config

const portNumber = ":8080"

// main is the entry point for the application
func main() {

	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	// create a new server mux and register the routes
	srv := &http.Server{
		Addr:         appConfig.Port,
		WriteTimeout: time.Second * 15, // 15 seconds write timeout
		ReadTimeout:  time.Second * 15, // 15 seconds read timeout
		IdleTimeout:  time.Second * 60, // 60 seconds idle timeout
		Handler:      routes.Routes(&appConfig),
	}

	// start the server
	go func() {
		log.Println("Starting server on port", appConfig.Port)
		err := srv.ListenAndServe()
		if err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)

}

func init() {

	// change this to false when in production
	appConfig.Debug = true

	// declare the error
	var err error

	// set the configuration
	appConfig.Port = portNumber
	appConfig.TemplateCache, err = render.CreateTemplateCache()
	if err != nil {
		log.Fatal("couldn't create template cache", err)
	}

	// initialize the handlers
	handlers.Init(handlers.NewDependency(&appConfig))
	render.NewTemplates(&appConfig)

}
