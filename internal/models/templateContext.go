package models

// change the module name to your own
import (
	"html/template"

	"github.com/Ahdeyyy/go-web/internal/forms"
)

// TemplateContext holds data sent from handlers to templates
type TemplateContext struct {
	StringMap       map[string]string
	IntMap          map[string]int
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	CSRFTemplateTag template.HTML
	Flash           string
	Warning         string
	Error           string
	Form            *forms.Form
	IsAuthenticated int
}
