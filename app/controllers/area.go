package controllers

import (
	"strconv"
	"strings"
	"time"

	"fmt"

	"github.com/kbse-mlg/gofence/app/geofence"
	"github.com/kbse-mlg/gofence/app/models"
	"github.com/kbse-mlg/gofence/app/modules/response"
	"github.com/kbse-mlg/gofence/app/routes"
	"github.com/kbse-mlg/gofence/utility/types"
	"github.com/revel/revel"
)

type Area struct {
	App
}

func (c Area) Index() revel.Result {
	moreStyles := []string{"/public/css/leaflet.css"}
	moreScripts := []string{"/public/js/leaflet.js"}
	IsAreas := true

	results, err := c.Txn.Select(models.Area{},
		`select * from "Area"`)
	if err != nil {
		panic(err)
	}

	var areas []*models.Area
	for _, r := range results {
		b := r.(*models.Area)
		areas = append(areas, b)
	}

	return c.Render(moreScripts, moreStyles, IsAreas, areas)
}

func (c Area) New() revel.Result {
	moreStyles := []string{"/public/css/leaflet.css", "/public/css/leaflet.pm.css"}
	moreScripts := []string{"/public/js/leaflet.js", "/public/js/leaflet.pm.min.js", "https://cdnjs.cloudflare.com/ajax/libs/leaflet-editable/1.0.0/Leaflet.Editable.min.js"}
	IsAreas := true

	return c.Render(moreScripts, moreStyles, IsAreas)
}

func (c Area) Edit(id int64) revel.Result {
	moreStyles := []string{"/public/css/leaflet.css", "/public/css/leaflet.pm.css"}
	moreScripts := []string{"/public/js/leaflet.js", "/public/js/leaflet.pm.min.js", "https://cdnjs.cloudflare.com/ajax/libs/leaflet-editable/1.0.0/Leaflet.Editable.min.js"}
	IsAreas := true
	var area *models.Area
	err := c.Txn.SelectOne(&area, `select * from "Area" WHERE "AreaID"=$1 LIMIT 1`, id)
	if err != nil {
		fmt.Println(err.Error())
		c.Flash.Error(err.Error())
	}

	return c.Render(moreScripts, moreStyles, IsAreas, area)
}

func (c Area) List(search string, size, page int) revel.Result {
	if page == 0 {
		page = 1
	}
	nextPage := page + 1
	search = strings.TrimSpace(search)

	var areas []*models.Area
	if search == "" {
		areas = loadAreas(c.Txn.Select(models.Area{},
			`select * from "Area" OFFSET $1 LIMIT $2`, (page-1)*size, size))
	} else {
		search = strings.ToLower(search)
		areas = loadAreas(c.Txn.Select(models.Area{},
			`select * from "Area" where lower(Name) like $1 or lower(Geodata) like $2
 OFFSET $3 LIMIT $4`, "%"+search+"%", "%"+search+"%", (page-1)*size, size))
	}

	return c.Render(areas, search, size, page, nextPage)
}

func (c Area) ListJson() revel.Result {

	size, _ := strconv.ParseInt(c.Params.Get("size"), 0, 64)
	search := c.Params.Route.Get("search")
	page, _ := strconv.ParseInt(c.Params.Get("page"), 0, 64)

	if page == 0 {
		page = 1
	}
	nextPage := page + 1

	if size == 0 {
		size = 10
	}
	search = strings.TrimSpace(search)

	var areas []*models.Area
	if search == "" {
		areas = loadAreas(c.Txn.Select(models.Area{},
			`select * from "Area" OFFSET $1 LIMIT $2`, (page-1)*size, size))
	} else {
		search = strings.ToLower(search)
		areas = loadAreas(c.Txn.Select(models.Area{},
			`select * from "Area" where lower(Name) like $1 or lower(Geodata) like $2
 OFFSET $3 LIMIT $4`, "%"+search+"%", "%"+search+"%", (page-1)*size, size))
	}

	result := models.AreaCollection{
		CurrentSearch: search,
		Areas:         areas,
		Size:          size,
		Page:          page,
		NextPage:      nextPage,
	}

	return c.RenderJSON(result)
}

func loadAreas(results []interface{}, err error) []*models.Area {
	if err != nil {
		panic(err)
	}
	var Areas []*models.Area
	for _, r := range results {
		Areas = append(Areas, r.(*models.Area))
	}
	return Areas
}

func (c Area) loadAreaById(id int) *models.Area {
	// var area *models.Area
	// err := c.Txn.SelectOne(&area, `SELECT * FROM "Area" WHERE "AreaID"=$1`, id)
	// if err != nil {
	// 	panic(err)
	// }
	// if area == nil {
	// 	return nil
	// }
	// return area

	area, err := c.Txn.Get(models.Area{}, id)
	if err != nil {
		panic(err)
	}

	if area == nil {
		return nil
	}

	return area.(*models.Area)
}

func (c Area) Show(id int) revel.Result {
	Area := c.loadAreaById(id)
	if Area == nil {
		return c.NotFound("Area %d does not exist", id)
	}
	title := Area.Name
	return c.Render(title, Area)
}

func (c Area) GetJson(id int) revel.Result {
	Area := c.loadAreaById(id)
	if Area == nil {
		return c.NotFound("Area %d does not exist", id)
	}
	return c.RenderJSON(Area)
}

func (c Area) SetHookWeb(id int) revel.Result {
	area := c.loadAreaById(id)
	if area == nil {
		c.Flash.Error(fmt.Sprintf("Area %d does not exist", id))
		return c.Redirect(routes.Area.Index())
	}
	activestring := c.Params.Form.Get(fmt.Sprintf("active%d", id))
	b, err := strconv.ParseBool(activestring)
	if err != nil {
		c.Flash.Error(err.Error())
		return c.Redirect(routes.Area.Index())
	}
	fmt.Println("-----> active ", b, area.Active)
	area.Active = b
	fmt.Println("-----> active updated", b, area.Active)
	_, err = c.Txn.Update(area)
	if err != nil {
		c.Flash.Error(err.Error())
		fmt.Println("-----> error", err.Error())
		return c.Redirect(routes.Area.Index())
	}

	if area.Active == true {
		geofence.SetFenceHook("", area.Name, area.Group, area.Geodata, "")
	} else {
		geofence.DeleteHook(area.Name)
	}

	return c.Redirect(routes.Area.Index())
}

func (c Area) SetHookAPI(id int) revel.Result {
	area := c.loadAreaById(id)

	var updatearea models.Area
	c.Params.BindJSON(&updatearea)
	area.Active = updatearea.Active
	if area == nil {
		return c.RenderJSON(response.ERROR(fmt.Sprintf("Area %d does not exist", id)))
	}

	if area.Active == true {
		geofence.SetFenceHook("", area.Name, area.Group, area.Geodata, "")
	} else {
		geofence.DeleteHook(area.Name)
	}

	return c.RenderJSON(response.OK())
}

func (c Area) ConfirmNew() revel.Result {

	geodata := c.Params.Form.Get("geodata")
	group := c.Params.Form.Get("group")
	name := c.Params.Form.Get("name")

	fmt.Printf("---> new data geodata:%s\nGroup:%s\nName:%s\n", geodata, group, name)
	now := types.DateTime{Int64: time.Now().UnixNano(), Valid: true}
	area := models.Area{0, name, geodata, 1, group, true, now, now}

	area.Validate(c.Validation)

	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(routes.Area.New)
	}

	// fmt.Println(geodata, group, area)
	err := c.Txn.Insert(&area)
	if err != nil {
		fmt.Println(err.Error())
	}

	geofence.SetFenceHook("fence", area.Name, area.Group, area.Geodata, ":6379")

	return c.Redirect(routes.Area.Index())
}

func (c Area) ConfirmEdit(id int) revel.Result {
	area := c.loadAreaById(id)
	if area == nil {
		return c.NotFound("Area %d does not exist", id)
	}

	geodata := c.Params.Form.Get("geodata")
	group := c.Params.Form.Get("group")

	area.Geodata = geodata
	area.Group = group

	// fmt.Println(geodata, group, area)
	// _, err := c.Txn.Query(`Update "Area" Set "Geodata"=$1, "Group"=$2 Where "AreaID"=$3`, geodata, group, id)
	_, err := c.Txn.Update(area)
	if err != nil {
		fmt.Println(err.Error())
	}

	geofence.SetFenceHook("", area.Name, area.Group, area.Geodata, "")

	return c.Redirect(routes.Area.Index())
}
