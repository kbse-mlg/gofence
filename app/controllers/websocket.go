package controllers

import (
	"encoding/json"

	"github.com/go-gorp/gorp"
	"github.com/kbse-mlg/gofence/app/geofence"
	"github.com/revel/revel"
	"golang.org/x/net/websocket"
)

type WebSocket struct {
	App
}

type ClientData struct {
	Command string  `json:"cmd"`
	Name    string  `json:"name"`
	Geojson string  `json:"geojson"`
	Group   string  `json:"group"`
	Long    float64 `json:"long"`
	Lat     float64 `json:"lat"`
}

var (
	newInternalWS = make(chan ClientData)
)

func (c WebSocket) Geofence(name string, ws *websocket.Conn) revel.Result {
	// Join the room.
	subscription := geofence.Subscribe()
	defer subscription.Cancel()

	geofence.Join(name)
	defer geofence.Leave(name)

	// Send down the archive.
	for _, event := range subscription.Archive {
		if websocket.JSON.Send(ws, &event) != nil {
			// They disconnected
			return nil
		}
	}

	// In order to select between websocket messages and subscription events, we
	// need to stuff websocket events into a channel.
	newMessages := make(chan string)
	go func() {
		var msg string
		for {
			err := websocket.Message.Receive(ws, &msg)
			if err != nil {
				close(newMessages)
				return
			}
			newMessages <- msg
		}
	}()

	// Now listen for new events from either the websocket
	for {
		select {
		case event := <-subscription.New:
			if websocket.JSON.Send(ws, &event) != nil {
				// They disconnected.
				return nil
			}
		case msg, ok := <-newMessages:
			// If the channel is closed, they disconnected.
			revel.TRACE.Println("-->", msg, ok)
			if !ok {
				return nil
			}
			revel.TRACE.Println(msg)
			var m ClientData
			err := json.Unmarshal([]byte(msg), &m)
			if err == nil {
				// put error message
			}
			doProcess(c.Txn, &m)
			// case cd, ok := <-newInternalWS:
			// 	revel.TRACE.Println("-->", cd, ok)
			// 	if !ok {
			// 		return nil
			// 	}
			// 	doProcess(c.Txn, &cd)
		}
	}
}

func doProcess(txn *gorp.Transaction, cd *ClientData) {
	switch cmd := cd.Command; cmd {
	case geofence.POSITION:
		updatePos(txn, cd.Name, cd.Lat, cd.Long)
		geofence.SetObject(cd.Name, cd.Group, cd.Lat, cd.Long)
		geofence.Position(cd.Name, cd.Lat, cd.Long)
	case geofence.SETHOOK:
		geofence.SetFenceHook(cd.Name, cd.Group, cd.Geojson, ":6379")
	case geofence.DELHOOK:
		geofence.DeleteHook(cd.Name)
	case geofence.STOPPED:
		geofence.Stopped(cd.Name, cd.Lat, cd.Long)
	default:
		revel.TRACE.Println("no process")
	}
}

func updatePos(txn *gorp.Transaction, name string, lat, long float64) {
	revel.TRACE.Println(name, lat, long)
	_, err := txn.Exec(`UPDATE "Object" SET "Lat"=$1, "Long"=$2 WHERE "Name"=$3`, lat, long, name)
	if err != nil {
		revel.TRACE.Println("--pos-----", err.Error())
	}
}

func updatePosById(txn *gorp.Transaction, id int64, lat, long float64) {
	_, err := txn.Exec(`UPDATE "Object" SET "Lat"=$1, "Long"=$2 where "ObjectID"=$3`, lat, long, id)
	if err != nil {
		revel.TRACE.Println("--id-----", err.Error())
	}
}

func newClientData(cmd, name, data, group string, long, lat float64) ClientData {
	return ClientData{cmd, name, data, group, long, lat}
}

func SendStoppedEvent(name string, long, lat float64) {
	newInternalWS <- newClientData(geofence.STOPPED, name, "", "", long, lat)
}
