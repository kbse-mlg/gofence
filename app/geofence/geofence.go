package geofence

import (
	"container/list"
	"time"
)

const (
	// JOIN subscribed event
	JOIN = "join"

	// LEAVE end subscribe
	LEAVE = "leave"

	// POSITION update event name
	POSITION = "position"

	// SETHOOK to register webhook to pubsub server using tile38
	SETHOOK = "sethook"

	// DELHOOK to delete assigned webhook
	DELHOOK = "delhook"

	// RESULT get the result of query
	RESULT = "result"

	//STOPPED get notification if object stopped
	STOPPED = "stopped"
)

// Event is a type for communication
//
type Event struct {
	Type      string  `json:"type"` // "create" "destroy" "position", "list", or "object"
	Name      string  `json:"name"`
	Timestamp int     `json:"timestamp"` // Unix timestamp (secs)
	Data      string  `json:"data"`      // What the user said (if Type == "message")
	Long      float64 `json:"long"`
	Lat       float64 `json:"lat"`
}

type Subscription struct {
	Archive []Event      // All the events from the archive.
	New     <-chan Event // New events coming in.
}

// Owner of a subscription must cancel it when they stop listening to events.
func (s Subscription) Cancel() {
	unsubscribe <- s.New // Unsubscribe the channel.
	drain(s.New)         // Drain it, just in case there was a pending publish.
}

func newEvent(typ, name, data string, long, lat float64) Event {
	return Event{typ, name, int(time.Now().Unix()), data, long, lat}
}

func Subscribe() Subscription {
	resp := make(chan Subscription)
	subscribe <- resp
	return <-resp
}

func Join(name string) {
	publish <- newEvent(JOIN, name, "", 0, 0)
}

func Leave(name string) {
	publish <- newEvent(LEAVE, name, "", 0, 0)
}

func Sethook(name, geojson string) {
	publish <- newEvent(SETHOOK, name, geojson, 0, 0)
}

func DeleteHook(name string) {
	publish <- newEvent(DELHOOK, name, "", 0, 0)
}

func Position(name string, long, lat float64) {
	publish <- newEvent(POSITION, name, "", long, lat)
}

func Stopped(name string, long, lat float64) {
	publish <- newEvent(STOPPED, name, "", long, lat)
}

func Result(name, data string) {
	publish <- newEvent(RESULT, name, data, 0, 0)
}

const archiveSize = 10

var (
	// Send a channel here to get room events back.  It will send the entire
	// archive initially, and then new messages as they come in.
	subscribe = make(chan (chan<- Subscription), 10)
	// Send a channel here to unsubscribe.
	unsubscribe = make(chan (<-chan Event), 10)
	// Send events here to publish them.
	publish = make(chan Event, 10)
)

// This function loops forever, handling the geofence pubsub
func geofence() {
	archive := list.New()
	subscribers := list.New()

	for {
		select {
		case ch := <-subscribe:
			var events []Event
			for e := archive.Front(); e != nil; e = e.Next() {
				events = append(events, e.Value.(Event))
			}
			subscriber := make(chan Event, 10)
			subscribers.PushBack(subscriber)
			ch <- Subscription{events, subscriber}

		case event := <-publish:
			for ch := subscribers.Front(); ch != nil; ch = ch.Next() {
				ch.Value.(chan Event) <- event
			}
			if archive.Len() >= archiveSize {
				archive.Remove(archive.Front())
			}
			archive.PushBack(event)

		case unsub := <-unsubscribe:
			for ch := subscribers.Front(); ch != nil; ch = ch.Next() {
				if ch.Value.(chan Event) == unsub {
					subscribers.Remove(ch)
					break
				}
			}
		}
	}
}

func init() {
	go geofence()
}

// Helpers

// Drains a given channel of any messages.
func drain(ch <-chan Event) {
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		default:
			return
		}
	}
}
