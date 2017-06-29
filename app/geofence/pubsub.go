package geofence

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

/*
{
	"command":"set",
	"group":"594c71a3cd32852ac892e1b4",
	"detect":"outside",
	"hook":"warehouse",
	"key":"fleet",
	"time":"2017-06-23T09:13:28.8918921+07:00",
	"id":"truck1",
	"object":{"type":"Point","coordinates":[-112.269,33.5123]}
}
*/
func init() {
	go pubSubRedis()
}

func pubSubRedis() {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		log.Fatalf("Could not connect: %v\n", err)
	}
	log.Println("Oke")
	defer c.Close()
	psc := redis.PubSubConn{Conn: c}
	psc.PSubscribe("fence.*")
	for {
		switch v := psc.Receive().(type) {
		case redis.PMessage:
			log.Printf("%s - %s: message: %s\n", v.Pattern, v.Channel, v.Data)
			Result(v.Channel, string(v.Data[:len(v.Data)]))
		case redis.Subscription:
			log.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case error:
			return
		}
	}
}
