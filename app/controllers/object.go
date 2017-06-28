package controllers

import (
	"github.com/kbse-mlg/gofence/app/models"
	"github.com/revel/revel"
)

type Object struct {
	App
}

func (c Object) Index() revel.Result {
	moreStyles := []string{"css/leaflet.css"}
	moreScripts := []string{"js/leaflet.js"}
	IsObjects := true

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
	return c.Render(moreScripts, moreStyles, IsObjects, objects)
}
