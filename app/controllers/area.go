package controllers

import (
	"strings"

	"github.com/kbse-mlg/gofence/app/models"
	"github.com/revel/revel"
)

type Area struct {
	App
}

func (c Area) Index() revel.Result {
	moreStyles := []string{"css/leaflet.css"}
	moreScripts := []string{"js/leaflet.js"}
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

func (c Area) List(search string, size, page int) revel.Result {
	if page == 0 {
		page = 1
	}
	nextPage := page + 1
	search = strings.TrimSpace(search)

	var areas []*models.Area
	if search == "" {
		areas = loadAreas(c.Txn.Select(models.Area{},
			`select * from Area limit ?, ?`, (page-1)*size, size))
	} else {
		search = strings.ToLower(search)
		areas = loadAreas(c.Txn.Select(models.Area{},
			`select * from Area where lower(Name) like ? or lower(Geodata) like ?
 limit ?, ?`, "%"+search+"%", "%"+search+"%", (page-1)*size, size))
	}

	return c.Render(areas, search, size, page, nextPage)
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
	h, err := c.Txn.Get(models.Area{}, id)
	if err != nil {
		panic(err)
	}
	if h == nil {
		return nil
	}
	return h.(*models.Area)
}

func (c Area) Show(id int) revel.Result {
	Area := c.loadAreaById(id)
	if Area == nil {
		return c.NotFound("Area %d does not exist", id)
	}
	title := Area.Name
	return c.Render(title, Area)
}
