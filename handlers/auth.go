package handlers

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func ShowRegisterPage(c *gin.Context) {
	tpl := pongo2.Must(pongo2.FromFile("templates/register.html"))

	output, err := tpl.Execute(pongo2.Context{})
	if err != nil {
		c.String(500, "Template rendering error")
		return
	}

	c.Header("Content-Type", "text/html")
	c.String(200, output)
}
