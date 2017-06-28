package controllers

import (
	"database/sql"

	"github.com/go-gorp/gorp"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	db "github.com/revel/modules/db/app"
	r "github.com/revel/revel"

	"github.com/kbse-mlg/gofence/app/models"
)

var (
	Dbm *gorp.DbMap
)

func InitDB() {
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.PostgresDialect{}}

	setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
		for col, size := range colSizes {
			t.ColMap(col).MaxSize = size
		}
	}

	t := Dbm.AddTable(models.Area{}).SetKeys(true, "AreaID")
	setColumnSizes(t, map[string]int{
		"Name":    100,
		"Geodata": 4096,
		"Type":    6,
	})

	t = Dbm.AddTable(models.Object{}).SetKeys(true, "ObjectID")
	setColumnSizes(t, map[string]int{
		"Name": 100,
		"Long": 11,
		"Lat":  11,
		"Type": 6,
	})

	Dbm.TraceOn("[gorp]", r.INFO)
	Dbm.CreateTables()

	demoArea := &models.Area{0, "A1",
		`{
				"type": "FeatureCollection",
				"features": [
						{
						"type": "Feature",
						"properties": {
								"color":"red"
						},
						"geometry": {
								"type": "Polygon",
								"coordinates": [
								[
										[
										101.67458295822144,
										3.1290962786081646
										],
										[
										101.67694330215454,
										3.127125114155911
										],
										[
										101.67750120162964,
										3.1273822227728822
										],
										[
										101.6785740852356,
										3.1284106566103684
										],
										[
										101.67887449264526,
										3.128796319039309
										],
										[
										101.67833805084227,
										3.129267684037543
										],
										[
										101.67923927307129,
										3.1302532647123797
										],
										[
										101.67709350585938,
										3.13222442328034
										],
										[
										101.67563438415527,
										3.130638926463226
										],
										[
										101.67522668838501,
										3.1299747311373816
										],
										[
										101.67458295822144,
										3.1290962786081646
										]
								]
								]
						}
						}
				]
		}`, 1}
	if err := Dbm.Insert(demoArea); err != nil {
		panic(err)
	}

	objects := []*models.Object{
		&models.Object{0, "A1", 101.67458295822144, 3.1290962786081646, 1},
	}

	for _, obj := range objects {
		if err := Dbm.Insert(obj); err != nil {
			panic(err)
		}
	}
}

type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

func (c *GorpController) Begin() r.Result {
	txn, err := Dbm.Begin()
	if err != nil {
		panic(err)
	}
	c.Txn = txn
	return nil
}

func (c *GorpController) Commit() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Commit(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}

func (c *GorpController) Rollback() r.Result {
	if c.Txn == nil {
		return nil
	}
	if err := c.Txn.Rollback(); err != nil && err != sql.ErrTxDone {
		panic(err)
	}
	c.Txn = nil
	return nil
}
