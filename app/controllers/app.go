package controllers

import (
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Index() revel.Result {
	moreStyles := []string{"css/leaflet.css"}
	moreScripts := []string{"js/leaflet.js"}
	IsDashboard := true
	return c.Render(moreScripts, moreStyles, IsDashboard)
}
