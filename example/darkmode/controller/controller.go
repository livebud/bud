package controller

import (
	"net/http"
)

type Controller struct {
	Writer http.ResponseWriter
}

func (c *Controller) Index() {}

func (c *Controller) Create(theme string) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:  "theme",
		Value: theme,
	})
}
