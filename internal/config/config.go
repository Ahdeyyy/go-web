package config

import "html/template"

// Config is the configuration for the application
type Config struct {

	// Port is the port the application will listen on
	Port string
	// Debug is whether the application is in debug mode
	Debug bool

	// TemplateCache is the template cache
	TemplateCache map[string]*template.Template
}
