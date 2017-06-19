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
	return c.Render(moreScripts, moreStyles)
}
