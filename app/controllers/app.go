package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	moreScripts := []string{
		"//maps.google.com/maps/api/js?sensor=true",
		"/public/js/gmap.js",
	}
	IsDashboard := true
	return c.Render(moreScripts, IsDashboard)
}
