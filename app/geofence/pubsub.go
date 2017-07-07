package geofence

import (
	"fmt"

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

var (
	redispool         *redis.Pool
	tile38Pool        *redis.Pool
	redisHookTemplate string = "redis://%s/%s"
)

func init() {
	go initPubSub()
}

func newPool(address string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

func initPubSub() {
	redispool = newPool(":6379")
	tile38Pool = newPool(":9851")
	for {
		c := redispool.Get()
		defer c.Close()
		psc := redis.PubSubConn{Conn: c}
		psc.PSubscribe("fence.*")

		for c.Err() == nil {
			switch v := psc.Receive().(type) {
			case redis.PMessage:
				fmt.Printf("> %s - %s: message: %s\n", v.Pattern, v.Channel, v.Data)
				Result(v.Channel, string(v.Data[:len(v.Data)]))
			case redis.Subscription:
				fmt.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				return
			}
		}
		fmt.Println(c.Err())
		c.Close()
	}
}

// SetObject set object name
func SetObject(name, group string, lat, long float64) {
	c := tile38Pool.Get()
	defer c.Close()
	fmt.Println("SET %s %s POINT %d %d", group, name, lat, long)
	ret, err := c.Do("SET", group, name, "POINT", lat, long)
	if err != nil {
		fmt.Printf("%v -- %v", ret, err)
	}
}

// SetFenceHook set webhook to redis
func SetFenceHook(name, group, geojson, redisAddress string) {
	c := tile38Pool.Get()
	defer c.Close()
	fmt.Println(group, geojson)
	ret, err := c.Do("SETHOOK", name, fmt.Sprintf(redisHookTemplate, redisAddress, name), "WITHIN", group, "FENCE", "OBJECT", geojson)
	if err != nil {
		fmt.Printf("%v -- %v", ret, err)
	}
}

// DelFenceHook delete webhook to redis
func DelFenceHook(name string) {
	c := tile38Pool.Get()
	defer c.Close()
	ret, err := c.Do("DELHOOK", name)
	if err != nil {
		fmt.Printf("%d -- %v", ret, err)
	}
}

// DelAllHook delete all webhook to redis
func DelAllHook() {
	c := tile38Pool.Get()
	defer c.Close()
	ret, err := c.Do("PDELHOOK", "*")
	if err != nil {
		fmt.Printf("%d -- %v", ret, err)
	}
}