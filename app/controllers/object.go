package controllers

import (
	"strconv"
	"strings"
	"time"

	"github.com/kbse-mlg/gofence/app/geofence"
	"github.com/kbse-mlg/gofence/app/models"
	"github.com/kbse-mlg/gofence/app/modules/response"
	"github.com/kbse-mlg/gofence/app/routes"
	"github.com/revel/revel"
)

type Object struct {
	App
}

func (c Object) Index() revel.Result {
	moreStyles := []string{"public/css/leaflet.css"}
	moreScripts := []string{"public/js/leaflet.js"}
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

func (c Object) UpdatePosition(name string) revel.Result {
	var obj models.Object
	var existingObj models.Object
	err := c.Params.BindJSON(&obj)
	if err != nil {
		return c.RenderJSON(response.ERROR(err.Error()))
	}
	err = c.Txn.SelectOne(&existingObj, `SELECT * FROM "Object" WHERE "Name"=$1`, name)
	if err != nil {
		return c.RenderJSON(response.ERROR(err.Error()))
	}
	checkStopped(&obj, &existingObj, existingObj.Name)
	existingObj.Lat = obj.Lat
	existingObj.Long = obj.Long
	existingObj.Group = obj.Group

	c.Txn.Update(&existingObj)
	geofence.SetObject(name, obj.Group, obj.Lat, obj.Long)
	geofence.Position(name, obj.Lat, obj.Long)
	return c.RenderJSON(response.OK())
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

func checkStopped(obj1, obj2 *models.Object, name string) {
	if !obj2.SameLoc(obj1) {
		geofence.SetTsObject(name)
		return
	}

	if ts, err := geofence.GetTsObject(name); err == nil {
		// stopped time
		stopped := revel.Config.IntDefault("app.map.stoppedtime", 1)
		now := time.Now().Add(time.Minute * (time.Duration(stopped))).UnixNano()
		revel.TRACE.Println("========> ts vs now ", ts, now)
		if now >= ts {
			revel.TRACE.Println("---------> Stoped 1 minutes")
			SendStoppedEvent(name, obj1.Lat, obj1.Long)
			geofence.SetTsObject(name)
		}
	} else {
		geofence.SetTsObject(name)
	}
}

func (c Object) History(id, search string, size, page int) revel.Result {
	object, err := c.Txn.Get(models.Object{}, id)
	if err != nil {
		return c.Redirect(routes.Object.Index())
	}
	if page == 0 {
		page = 1
	}
	if size == 0 {
		size = 10
	}
	nextPage := page + 1
	prevPage := page
	search = strings.TrimSpace(search)

	var movehistories []*models.MoveHistory
	if search == "" {
		movehistories = c.loadLocationHistory(c.Txn.Select(models.MoveHistory{},
			`select * from "MoveHistory" WHERE "ObjectID"=$1 OFFSET $2 LIMIT $3`, id, (page-1)*size, size))
	} else {
		search = strings.ToLower(search)
		movehistories = c.loadLocationHistory(c.Txn.Select(models.MoveHistory{},
			`select * from "MoveHistory" where lower(Name) like $1 AND "ObjectID"=$2
 OFFSET $3 LIMIT $4`, "%"+search+"%", id, (page-1)*size, size))
	}

	if movehistories == nil || len(movehistories) < size {
		nextPage = -1
		prevPage = page - 1
	}

	return c.Render(object, movehistories, search, size, page, nextPage, prevPage)
}

func (c Object) loadLocationHistory(results []interface{}, err error) []*models.MoveHistory {
	if err != nil {
		revel.TRACE.Fatal(err.Error())
		return nil
	}
	tz, err := strconv.Atoi(c.Session["tz"])
	if err != nil {
		tz = -420
	}

	var Histories []*models.MoveHistory
	for _, r := range results {
		hist := r.(*models.MoveHistory)
		hist.Created.TimezoneOffset = tz
		hist.Created.Relative = false
		hist.StringCreated = hist.Created.ToString()
		Histories = append(Histories, hist)
	}
	return Histories
}
