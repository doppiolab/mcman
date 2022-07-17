package server

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type templateRenderer struct {
	templates *template.Template
}

func (t *templateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	// NOTE(hayeon): add global metadata if required
	//               https://echo.labstack.com/guide/templates/#advanced---calling-echo-from-templates
	return t.templates.ExecuteTemplate(w, name, data)
}
