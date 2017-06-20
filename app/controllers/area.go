package controllers

import (
	"github.com/revel/revel"
)

type Area struct {
	*revel.Controller
}

func (c Area) Index() revel.Result {
	moreStyles := []string{"css/leaflet.css"}
	moreScripts := []string{"js/leaflet.js"}
	IsAreas := true
	return c.Render(moreScripts, moreStyles, IsAreas)
}
