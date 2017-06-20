package controllers

import (
	"github.com/revel/revel"
)

type Object struct {
	*revel.Controller
}

func (c Object) Index() revel.Result {
	moreStyles := []string{"css/leaflet.css"}
	moreScripts := []string{"js/leaflet.js"}
	IsObjects := true
	return c.Render(moreScripts, moreStyles, IsObjects)
}
