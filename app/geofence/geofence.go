package geofence

import (
	"container/list"
	"time"
)

type Event struct {
	Type      string // "create" "destroy" "position", "list", or "object"
	Name      string
	Timestamp int    // Unix timestamp (secs)
	Data      string // What the user said (if Type == "message")
	Long      float64
	Lat       float64
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
	publish <- newEvent("join", name, "", 0, 0)
}

func Leave(name string) {
	publish <- newEvent("leave", name, "", 0, 0)
}

func Sethook(name, geojson string) {
	publish <- newEvent("create", name, geojson, 0, 0)
}

func DeleteHook(name string) {
	publish <- newEvent("destroy", name, "", 0, 0)
}

func Position(name string, long, lat float64) {
	publish <- newEvent("position", name, "", long, lat)
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

// This function loops forever, handling the chat room pubsub
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
