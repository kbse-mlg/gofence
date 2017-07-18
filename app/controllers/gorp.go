package controllers

import (
	"database/sql"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-gorp/gorp"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	db "github.com/revel/modules/db/app"
	r "github.com/revel/revel"

	"github.com/kbse-mlg/gofence/app/geofence"
	"github.com/kbse-mlg/gofence/app/models"
)

var (
	Dbm *gorp.DbMap
)

type Coord struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

func InitDB() {
	db.Init()
	Dbm = &gorp.DbMap{Db: db.Db, Dialect: gorp.PostgresDialect{}}
	setColumnSizes := func(t *gorp.TableMap, colSizes map[string]int) {
		for col, size := range colSizes {
			t.ColMap(col).MaxSize = size
		}
	}

	t := Dbm.AddTableWithName(models.User{}, "User").SetKeys(true, "UserID")
	t.ColMap("Password").Transient = true
	setColumnSizes(t, map[string]int{
		"Username": 20,
		"Name":     100,
	})

	t = Dbm.AddTableWithName(models.Area{}, "Area").SetKeys(true, "AreaID")
	setColumnSizes(t, map[string]int{
		"Name":    100,
		"Geodata": 4096,
		"Type":    6,
		"Group":   100,
	})

	t = Dbm.AddTableWithName(models.Object{}, "Object").SetKeys(true, "ObjectID")
	setColumnSizes(t, map[string]int{
		"Name":  100,
		"Group": 100,
	})

	t = Dbm.AddTableWithName(models.MoveHistory{}, "MoveHistory").SetKeys(true, "ObjectID")
	t.ColMap("Object").Transient = true

	Dbm.TraceOn("[gorp]", r.INFO)
	Dbm.CreateTables()
	// InsertData()

	// Dummy Moving
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	inc := 0.0001
	direction, counter := 1, 0

	a1 := Coord{Lat: 3.1270, Long: 101.6772}
	a2 := Coord{Lat: 3.1299, Long: 101.6738}
	max := 17
	go func() {
		for {
			select {
			case <-ticker.C:
				if counter < max*1 {
					direction = 1
				} else if counter < max*2 {
					direction = 2
				} else if counter < max*3 {
					direction = 3
				} else if counter < max*4 {
					direction = 4
				} else {
					direction = 1
					counter = 0
					a1 = Coord{Lat: 3.1270, Long: 101.6772}
					a2 = Coord{Lat: 3.1299, Long: 101.6738}
				}
				counter++
				// do stuff
				switch direction {
				case 1:
					a1.Long -= inc
					a2.Long -= inc
				case 2:
					a1.Lat -= inc
					a2.Lat -= inc
				case 3:
					a1.Long += inc
					a2.Long += inc
				case 4:
					a1.Long += inc
					a2.Long += inc
				}

				txn, err := Dbm.Begin()
				if err != nil {
					panic(err)
				}
				txn.Exec(`UPDATE "Object" SET "Lat"=$1, "Long"=$2 WHERE "Name"=$3`, a1.Lat, a1.Long, "A1")
				txn.Exec(`UPDATE "Object" SET "Lat"=$1, "Long"=$2 WHERE "Name"=$3`, a2.Lat, a2.Long, "A2")
				if err := txn.Commit(); err != nil && err != sql.ErrTxDone {
					txn.Rollback()
				}
				geofence.SetObject("A1", "Truck", a1.Lat, a1.Long)
				geofence.SetObject("A2", "Truck", a2.Lat, a2.Long)
				geofence.Position("A1", a1.Lat, a1.Long)
				geofence.Position("A2", a2.Lat, a2.Long)
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func InsertData() {
	bcryptPassword, _ := bcrypt.GenerateFromPassword(
		[]byte("demo"), bcrypt.DefaultCost)
	demoUser := &models.User{0, "Demo User", "demo", "demo", bcryptPassword}
	if err := Dbm.Insert(demoUser); err != nil {
		panic(err)
	}

	demoArea := &models.Area{0, "A1",
		`{"type":"FeatureCollection","features":[{"type":"Feature","properties":{"color":"red"},"geometry":{"type":"Polygon","coordinates":[[[101.67458295822144,3.1290962786081646],[101.67694330215454,3.127125114155911],[101.67750120162964,3.1273822227728822],[101.6785740852356,3.1284106566103684],[101.67887449264526,3.128796319039309],[101.67833805084227,3.129267684037543],[101.67923927307129,3.1302532647123797],[101.67709350585938,3.13222442328034],[101.67563438415527,3.130638926463226],[101.67522668838501,3.1299747311373816],[101.67458295822144,3.1290962786081646]]]}}]}`,
		1,
		"Truck",
		false,
		time.Now().UnixNano(),
		time.Now().UnixNano()}
	if err := Dbm.Insert(demoArea); err != nil {
		panic(err)
	}

	objects := []*models.Object{
		&models.Object{0, "Truck", "A1", 101.67458295822144, 3.1290962786081646, 1, time.Now().UnixNano(), time.Now().UnixNano()},
		&models.Object{0, "Truck", "A2", 101.67478295822144, 3.1290962786081646, 1, time.Now().UnixNano(), time.Now().UnixNano()},
	}

	for _, obj := range objects {
		if err := Dbm.Insert(obj); err != nil {
			panic(err)
		}
	}
}

// GorpController controller wrapper which contains transaction database
type GorpController struct {
	*r.Controller
	Txn *gorp.Transaction
}

// Begin execute data
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
