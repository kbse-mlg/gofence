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
	Data    string  `json:"data"`
	Long    float64 `json:"long"`
	Lat     float64 `json:"lat"`
}

func (c WebSocket) SetHook(name, geojson string, ws *websocket.Conn) revel.Result {
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

	// Now listen for new events from either the websocket or the chatroom.
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
			var m ClientData
			err := json.Unmarshal([]byte(msg), &m)
			if err == nil {
				// put error message
			}

			doProcess(&m)
		}
	}
	return nil
}

func doProcess(clientData *ClientData) {
	switch cmd := clientData.Command; cmd {
	case geofence.POSITION:
		geofence.Position(clientData.Data, 10, 10)
	case geofence.SETHOOK:
	case geofence.DELHOOK:
	case geofence.JOIN:
	case geofence.LEAVE:
	case geofence.RESULT:
	default:
		geofence.Position("lalalala", 10, 10)
	}
}
