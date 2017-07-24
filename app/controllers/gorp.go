package controllers

import (
	"database/sql"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/go-gorp/gorp"
	//_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
	db "github.com/revel/modules/db/app"
	r "github.com/revel/revel"

	"github.com/kbse-mlg/gofence/app/geofence"
	"github.com/kbse-mlg/gofence/app/models"
	"github.com/kbse-mlg/gofence/utility/types"
)

type Coord struct {
	Lat  float64 `json:"lat"`
	Long float64 `json:"long"`
}

var (
	Dbm   *gorp.DbMap
	Route []Coord = []Coord{
		{3.127437461, 101.676402334},
		{3.127849069, 101.676072003},
		{3.128427563, 101.675369265},
		{3.128834651, 101.674891831},
		{3.129235546, 101.674454212},
		{3.129632757, 101.674628975},
		{3.130141617, 101.674097897},
		{3.129954142, 101.673861863},
		{3.130334448, 101.673588278},
		{3.130120191, 101.673185946},
		{3.129980924, 101.672912361},
		{3.129713104, 101.673116209},
		{3.129488134, 101.673320057},
		{3.129263165, 101.673534634},
		{3.129032839, 101.673738481},
		{3.128834651, 101.673953058},
		{3.128598968, 101.674189093},
		{3.128373999, 101.674430491},
		{3.128202593, 101.674618246},
		{3.128047257, 101.674822094},
		{3.127849069, 101.675036671},
		{3.127715158, 101.675224425},
		{3.127575891, 101.675358536},
		{3.127425074, 101.675886512},
		{3.127120595, 101.675557019},
		{3.126638516, 101.675337078},
		{3.126177863, 101.675197603},
		{3.125743991, 101.674848916},
		{3.125347615, 101.674527051},
		{3.125010159, 101.674194457},
		{3.124811971, 101.673947694},
		{3.124624495, 101.674130084},
		{3.124485228, 101.674339296},
		{3.124303109, 101.67453778},
		{3.12402993, 101.674865009},
		{3.123853167, 101.675090315},
		{3.123665692, 101.675320985},
		{3.123537137, 101.675535562},
		{3.123392513, 101.675707223},
		{3.123183611, 101.675927164},
		{3.123017561, 101.676163198},
		{3.123296097, 101.676447513},
		{3.123590701, 101.676731827},
		{3.123896019, 101.676914217},
		{3.124169197, 101.677150251},
		{3.124404881, 101.677289726},
		{3.124892317, 101.677584769},
		{3.125245842, 101.677751066},
		{3.125658288, 101.677944185},
		{3.126161793, 101.678185584},
		{3.126681367, 101.67841089},
		{3.127104525, 101.678609373},
		{3.127431268, 101.678904416},
		{3.127940128, 101.679274561},
		{3.128320434, 101.679574968},
		{3.128593612, 101.67918873},
		{3.128781087, 101.678974153},
		{3.128513266, 101.678738119},
		{3.128202593, 101.678475263},
		{3.127806218, 101.678013923},
		{3.127527683, 101.677713515},
		{3.127297357, 101.677246811},
		{3.127275931, 101.6767962},
		{3.127415198, 101.67642069},
	}
)

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
	t.AddIndex("NameIndex", "Btree", []string{"Name"}).SetUnique(true)

	t = Dbm.AddTableWithName(models.Object{}, "Object").SetKeys(true, "ObjectID")
	setColumnSizes(t, map[string]int{
		"Name":  100,
		"Group": 100,
	})

	t = Dbm.AddTableWithName(models.MoveHistory{}, "MoveHistory").SetKeys(true, "HistoryID")
	t.ColMap("Object").Transient = true
	t = Dbm.AddTableWithName(models.Alert{}, "Alert").SetKeys(true, "ID")
	setColumnSizes(t, map[string]int{
		"Info": 100,
	})

	Dbm.TraceOn("[gorp]", r.INFO)
	Dbm.CreateTablesIfNotExists()
	// InsertData()

	// Dummy Moving
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	counter := 0

	var a1, a2, a3 *models.Object
	max := len(Route)
	fmt.Println("--->", max, a1, a2)

	go func() {
		for {
			select {
			case <-ticker.C:
				if counter == max {
					counter = 0
				}

				c1 := Route[counter]
				c2 := Route[max-counter-1]
				counter++
				fmt.Println("--->", counter, c1, c2)

				txn, err := Dbm.Begin()
				if err != nil {
					fmt.Println("----error", err.Error())
					panic(err)
				}
				// get data from db
				if a1 == nil {
					af1, err := txn.Get(models.Object{}, 1)
					if err == nil {
						a1 = af1.(*models.Object)
					}
				}

				if a2 == nil {
					af2, err := txn.Get(models.Object{}, 2)
					if err == nil {
						a2 = af2.(*models.Object)
					}
				}

				if a3 == nil {
					af3, err := txn.Get(models.Object{}, 3)
					if err == nil {
						a3 = af3.(*models.Object)
					}
				}

				a1.Lat = c1.Lat
				a1.Long = c1.Long
				a2.Lat = c2.Lat
				a2.Long = c2.Long
				a3.Lat = 3.1290962786081646
				a3.Long = 101.67458295822144
				checkStopped(a3, a3, a3.Name)

				txn.Update(a1)
				txn.Update(a2)
				if err := txn.Commit(); err != nil && err != sql.ErrTxDone {
					fmt.Println("----error", err.Error())
					txn.Rollback()
				}
				geofence.SetObject("A1", "Truck", c1.Lat, c1.Long)
				geofence.SetObject("A2", "Truck", c2.Lat, c2.Long)
				geofence.SetObject("A3", "Truck", 3.1290962786081646, 101.67458295822144)
				geofence.Position("A1", c1.Lat, c1.Long)
				geofence.Position("A2", c2.Lat, c2.Long)
				geofence.Position("A3", 3.1290962786081646, 101.67458295822144)
			case <-quit:
				fmt.Println("----stopped")
				ticker.Stop()
				// return
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

	now := types.DateTime{Int64: time.Now().UnixNano(), Valid: true}
	demoArea := &models.Area{0, "A1",
		`{"type":"FeatureCollection","features":[{"type":"Feature","properties":{"color":"red"},"geometry":{"type":"Polygon","coordinates":[[[101.67458295822144,3.1290962786081646],[101.67694330215454,3.127125114155911],[101.67750120162964,3.1273822227728822],[101.6785740852356,3.1284106566103684],[101.67887449264526,3.128796319039309],[101.67833805084227,3.129267684037543],[101.67923927307129,3.1302532647123797],[101.67709350585938,3.13222442328034],[101.67563438415527,3.130638926463226],[101.67522668838501,3.1299747311373816],[101.67458295822144,3.1290962786081646]]]}}]}`,
		1,
		"Truck",
		false,
		now,
		now}
	if err := Dbm.Insert(demoArea); err != nil {
		panic(err)
	}

	objects := []*models.Object{
		&models.Object{0, "Truck", "A1", 101.67458295822144, 3.1290962786081646, 1, now, now},
		&models.Object{0, "Truck", "A2", 101.67478295822144, 3.1290962786081646, 1, now, now},
		&models.Object{0, "Truck", "A3", 101.67478295822144, 3.1290962786081646, 1, now, now},
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
