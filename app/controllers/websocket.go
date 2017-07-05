package controllers

import (
	"encoding/json"

	"github.com/revel/revel"
	"golang.org/x/net/websocket"

	"github.com/kbse-mlg/gofence/app/geofence"
)

type WebSocket struct {
	*revel.Controller
}

type ClientData struct {
	Command string  `json:"cmd"`
	Name    string  `json:"name"`
	Geojson string  `json:"geojson"`
	Group   string  `json:"group"`
	Long    float64 `json:"long"`
	Lat     float64 `json:"lat"`
}

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
			if !ok {
				return nil
			}
			revel.TRACE.Println(msg)
			var m ClientData
			err := json.Unmarshal([]byte(msg), &m)
			if err == nil {
				// put error message
			}
			doProcess(&m)
		}
	}
}

func doProcess(cd *ClientData) {
	switch cmd := cd.Command; cmd {
	case geofence.POSITION:
		geofence.Position(cd.Name, cd.Lat, cd.Long)
	case geofence.SETHOOK:
		geofence.Sethook(cd.Name, cd.Geojson)
		geofence.SetGeofenceHook(cd.Name, cd.Group, cd.Geojson, ":6379")
	case geofence.DELHOOK:
		geofence.DeleteHook(cd.Name)
	default:
		revel.TRACE.Println("no process")
	}
}
