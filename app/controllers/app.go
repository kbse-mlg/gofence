package controllers

import (
	"github.com/kbse-mlg/gofence/app/models"
	"github.com/revel/revel"
)

type App struct {
	GorpController
}

func (c App) Index() revel.Result {
	moreScripts := []string{
		"https://maps.googleapis.com/maps/api/js?key=AIzaSyBJFNacrQSkWIUsbZjLw4wHo0yyF9DDrgE",
		"/public/js/app/socket.js",
	}
	moreStyles := []string{
		"https://unpkg.com/leaflet@1.0.3/dist/leaflet.css",
	}
	IsDashboard := true

	results, err := c.Txn.Select(models.Object{},
		`select * from "Object"`)
	if err != nil {
		panic(err)
	}

	var objects []*models.Object
	for _, r := range results {
		b := r.(*models.Object)
		objects = append(objects, b)
	}

	return c.Render(moreScripts, moreStyles, IsDashboard, objects)
}
