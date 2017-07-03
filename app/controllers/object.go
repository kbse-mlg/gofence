package controllers

import (
	"strconv"
	"strings"

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

func (c Object) ListJson() revel.Result {

	size, _ := strconv.ParseInt(c.Params.Route.Get("size"), 0, 64)
	search := c.Params.Route.Get("search")
	page, _ := strconv.ParseInt(c.Params.Route.Get("page"), 0, 64)

	if page == 0 {
		page = 1
	}

	if size == 0 {
		size = 10
	}
	search = strings.TrimSpace(search)

	var Objects []*models.Object
	if search == "" {
		Objects = loadObjects(c.Txn.Select(models.Object{},
			`select * from "Object" OFFSET $1 LIMIT $2`, (page-1)*size, size))
	} else {
		search = strings.ToLower(search)
		Objects = loadObjects(c.Txn.Select(models.Object{},
			`select * from "Object" where lower(Name) like $1 or lower(Geodata) like $2
 OFFSET $3 LIMIT $4`, "%"+search+"%", "%"+search+"%", (page-1)*size, size))
	}

	result := models.ObjectCollection{
		CurrentSearch: search,
		Objects:       Objects,
		Size:          size,
		Page:          page,
	}

	return c.RenderJSON(result)
}

func loadObjects(results []interface{}, err error) []*models.Object {
	if err != nil {
		panic(err)
	}
	var Objects []*models.Object
	for _, r := range results {
		Objects = append(Objects, r.(*models.Object))
	}
	return Objects
}
