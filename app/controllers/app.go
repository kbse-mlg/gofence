package controllers

import (
	"github.com/kbse-mlg/gofence/app/models"
	"github.com/kbse-mlg/gofence/app/routes"
	"github.com/revel/revel"
)

type App struct {
	GorpController
}

func (c App) connected() *models.User {
	if c.ViewArgs["user"] != nil {
		return c.ViewArgs["user"].(*models.User)
	}
	if username, ok := c.Session["user"]; ok {
		return c.getUser(username)
	}
	return nil
}

func (c App) getUser(username string) *models.User {
	users, err := c.Txn.Select(models.User{}, `select * from "User" where "Username" = $1`, username)
	if err != nil {
		panic(err)
	}
	if len(users) == 0 {
		return nil
	}
	return users[0].(*models.User)
}

func (c App) Index() revel.Result {
	moreScripts := []string{
		"https://maps.googleapis.com/maps/api/js?key=AIzaSyBJFNacrQSkWIUsbZjLw4wHo0yyF9DDrgE",
		"/public/js/app/socket.js",
		"/public/js/custom_js/pnotify.custom.min.js",
		"http://underscorejs.org/underscore-min.js",
	}
	moreStyles := []string{
		"https://unpkg.com/leaflet@1.0.3/dist/leaflet.css",
		"/public/css/custom_css/pnotify.custom.min.css",
		"/public/css/timeline.css",
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

func (c App) checkUser() revel.Result {
	if user := c.connected(); user == nil {
		c.Flash.Error("Please log in first")
		return c.Redirect(routes.App.Index())
	}
	return nil
}
