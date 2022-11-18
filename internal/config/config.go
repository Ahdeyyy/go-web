package config

import (
	"html/template"
	"log"

	"github.com/gorilla/sessions"
)

// Config is the configuration for the application
type Config struct {
	Port          string                        // Port to listen on
	Debug         bool                          // Debug mode
	UseCache      bool                          // Use cache
	InfoLog       *log.Logger                   // Info log
	ErrorLog      *log.Logger                   // Error log
	SessionStore  *sessions.CookieStore         // SessionStore
	Session       *sessions.Session             // Session
	TemplateCache map[string]*template.Template // Template cache
}
