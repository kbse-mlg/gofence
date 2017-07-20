package controllers

import (
	"strconv"
	"strings"

	"github.com/kbse-mlg/gofence/app/models"
	"github.com/revel/revel"
)

type LocationHistory struct {
	App
}

func (c LocationHistory) ListJson() revel.Result {

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
	revel.TRACE.Printf("size:%d page:%d\n", size, page)

	var movehistories []*models.MoveHistory
	if search == "" {
		movehistories = loadLocationHistory(c.Txn.Select(models.MoveHistory{},
			`select * from "MoveHistory" OFFSET $1 LIMIT $2`, (page-1)*size, size))
	} else {
		search = strings.ToLower(search)
		movehistories = loadLocationHistory(c.Txn.Select(models.MoveHistory{},
			`select * from "MoveHistory" where lower(Name) like $1
 OFFSET $2 LIMIT $3`, "%"+search+"%", (page-1)*size, size))
	}

	result := models.MoveHistoriesCollection{
		CurrentSearch: search,
		MoveHistories: movehistories,
		Size:          size,
		Page:          page,
		NextPage:      nextPage,
	}

	return c.RenderJSON(result)
}

func loadLocationHistory(results []interface{}, err error) []*models.MoveHistory {
	if err != nil {
		revel.TRACE.Fatal(err.Error())
		return nil
	}
	var Histories []*models.MoveHistory
	for _, r := range results {
		Histories = append(Histories, r.(*models.MoveHistory))
	}
	return Histories
}
