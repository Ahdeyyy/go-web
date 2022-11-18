package main

import (
	"context"
	"crypto/rand"
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
	"github.com/gorilla/sessions"
)

// config is the configuration for the application
var appConfig config.Config
var sessionStore *sessions.CookieStore
var session *sessions.Session
var infoLog *log.Logger
var errorLog *log.Logger

// wait is the time to wait before shutting down
var wait time.Duration

var debug bool

// portNumber is the port number the server will listen on
const portNumber = ":8080"

// main is the entry point for the application
func main() {

	// run the application
	if err := run(); err != nil {
		log.Fatal(err)
	}

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

func run() error {

	// read flags
	flag.DurationVar(&wait, "timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.BoolVar(&debug, "debug", false, "enable debug mode")
	// dbHost := flag.String("dbhost", "localhost", "Database host")
	// dbName := flag.String("dbname", "", "Database name")
	// dbUser := flag.String("dbuser", "", "Database user")
	// dbPass := flag.String("dbpass", "", "Database password")
	// dbPort := flag.String("dbport", "5432", "Database port")
	// dbSSL := flag.String("dbssl", "disable", "Database ssl settings (disable, prefer, require)")

	flag.Parse()

	// change this to false when in production

	if debug {
		appConfig.Debug = true
	}

	log.Println("debug mode is ", debug)

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appConfig.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appConfig.ErrorLog = errorLog

	// declare the error
	var err error

	// session

	key := make([]byte, 32)
	_, err = rand.Read(key)
	if err != nil {
		return err
	}

	os.Setenv("SESSION_KEY", string(key))
	sessionStore = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	session = sessions.NewSession(sessionStore, "app-session")
	session.Values["count"] = 0

	// TODO set session options

	// set the configuration
	appConfig.Port = portNumber
	appConfig.TemplateCache, err = render.CreateTemplateCache()
	if err != nil {
		return err
	}
	appConfig.SessionStore = sessionStore
	appConfig.Session = session

	// initialize the handlers
	handlers.Init(handlers.NewDependency(&appConfig))
	render.NewTemplates(&appConfig)

	log.Printf("Server configs: %+v , shutdown timeout: %v", appConfig, wait)

	return nil
}
