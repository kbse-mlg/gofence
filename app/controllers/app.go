package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	moreScripts := []string{
		"https://maps.googleapis.com/maps/api/js?key=AIzaSyBJFNacrQSkWIUsbZjLw4wHo0yyF9DDrgE",
		"/public/js/gmap.js",
	}
	IsDashboard := true
	return c.Render(moreScripts, IsDashboard)
}
