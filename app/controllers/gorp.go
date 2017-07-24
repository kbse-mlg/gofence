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
		{2.737913927, 101.677829772},
		{2.73894272, 101.67877391},
		{2.737313798, 101.68006137},
		{2.736370738, 101.680662185},
		{2.734913279, 101.681606323},
		{2.733713017, 101.682636291},
		{2.734484614, 101.68383792},
		{2.735684875, 101.685812026},
		{2.737142333, 101.687957793},
		{2.738171126, 101.689588577},
		{2.739199917, 101.691476852},
		{2.740314441, 101.692249328},
		{2.741343231, 101.69370845},
		{2.742114823, 101.695425063},
		{2.743658005, 101.696970016},
		{2.74460106, 101.698514968},
		{2.745201186, 101.700489074},
		{2.745972775, 101.701862365},
		{2.747173024, 101.703664809},
		{2.748459004, 101.704866439},
		{2.749573519, 101.70623973},
		{2.750688033, 101.708042175},
		{2.751545351, 101.708986312},
		{2.7524884, 101.710702926},
		{2.752831326, 101.711904556},
		{2.753774374, 101.712934524},
		{2.754374495, 101.713964492},
		{2.754374495, 101.714736968},
		{2.752659863, 101.715595275},
		{2.751545351, 101.716453582},
		{2.750173642, 101.717054397},
		{2.748973396, 101.717741042},
		{2.747601685, 101.718856841},
		{2.7465729, 101.71988681},
		{2.745715579, 101.718771011},
		{2.745201186, 101.717311889},
		{2.743915202, 101.715595275},
		{2.742543485, 101.714221984},
		{2.741771893, 101.711904556},
		{2.739885778, 101.709758788},
		{2.739371383, 101.708299667},
		{2.736443075, 101.703701019},
		{2.739971511, 101.702034026},
		{2.742629217, 101.700832397},
		{2.744086667, 101.699716598},
		{2.743315076, 101.699029952},
		{2.742457753, 101.697828323},
		{2.741600428, 101.695940048},
		{2.740485906, 101.694309264},
		{2.73928565, 101.69250682},
		{2.737913927, 101.690876037},
		{2.736542203, 101.688215286},
		{2.734570347, 101.686584502},
		{2.733541551, 101.684267074},
		{2.732512754, 101.68375209},
		{2.731226757, 101.684267074},
		{2.729340626, 101.683923751},
		{2.728140359, 101.683580428},
		{2.731226757, 101.682035476},
		{2.733627284, 101.680147201},
		{2.736027806, 101.679117233},
		{2.737313798, 101.677915603},
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
	InsertData()

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
				geofence.SetObject("A3", "Truck", 2.737313798, 101.677915603)
				geofence.Position("A1", c1.Lat, c1.Long)
				geofence.Position("A2", c2.Lat, c2.Long)
				geofence.Position("A3", 2.737313798, 101.677915603)
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
	demoArea := []*models.Area{
		&models.Area{0, "A1",
			`{"type":"Feature","properties":{},"geometry":{"type":"Polygon","coordinates":[[[101.68494030833247,2.738029130778824],[101.68742939829826,2.7340425534504305],[101.6932658851147,2.7376861999169506],[101.69073387980464,2.740901172892378],[101.68494030833247,2.738029130778824]]]}}`,
			1,
			"Truck",
			false,
			now,
			now},
		&models.Area{0, "A2",
			`{"type":"Feature","properties":{},"geometry":{"type":"Polygon","coordinates":[[[101.70365139842036,2.751317625959215],[101.70678421854976,2.7462594419877515],[101.71339318156245,2.7506746376610627],[101.71004578471187,2.7554756085787466],[101.70365139842036,2.751317625959215]]]}}`,
			1,
			"Truck",
			false,
			now,
			now},
	}
	for _, ar := range demoArea {
		if err := Dbm.Insert(ar); err != nil {
			panic(err)
		}
	}

	objects := []*models.Object{
		&models.Object{0, "Truck", "A1", 101.71339318156245, 2.7506746376610627, 1, now, now},
		&models.Object{0, "Truck", "A2", 101.71339318156245, 2.7506746376610627, 1, now, now},
		&models.Object{0, "Truck", "A3", 101.71339318156245, 2.7506746376610627, 1, now, now},
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
