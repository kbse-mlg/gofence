package geofence

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/revel/revel"
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
	timestampObject   string = "ts:obj:%s"
)

func init() {
	go initPubSub()
}

func newPool(address string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		MaxActive:   500, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				revel.TRACE.Fatal(err.Error())
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

		// ct := tile38Pool.Get()
		// defer ct.Close()

		psc := redis.PubSubConn{Conn: c}
		psc.PSubscribe("fence.*")

		// experiment keep alive redis
		// go func() {
		// 	for {
		// 		c.Do("PING")
		// 		ct.Do("PING")
		// 		time.Sleep(20 * time.Second)
		// 	}
		// }()

		for c.Err() == nil {
			switch v := psc.Receive().(type) {
			case redis.PMessage:
				revel.TRACE.Printf("> %s - %s: message: %s\n", v.Pattern, v.Channel, v.Data)
				Result(v.Channel, string(v.Data[:len(v.Data)]))
			case redis.Subscription:
				revel.TRACE.Printf("%s: %s %d\n", v.Channel, v.Kind, v.Count)
			case error:
				return
			}
		}
		revel.TRACE.Println(c.Err())
		c.Close()
	}
}

// SetObject set object name
func SetObject(name, group string, lat, long float64) error {
	c := tile38Pool.Get()
	defer c.Close()
	revel.TRACE.Printf("SET %s %s POINT %f %f", group, name, lat, long)
	ret, err := c.Do("SET", group, name, "POINT", lat, long)
	if err != nil {
		revel.TRACE.Printf("%v -- %v", ret, err)
		return err
	}

	return nil
}

// SetFenceHook set webhook to redis
func SetFenceHook(prefix, name, group, geojson, redisAddress string) error {
	if redisAddress == "" {
		redisAddress = ":6379"
	}

	if prefix == "" {
		prefix = "fence"
	}
	pname := fmt.Sprintf("%s.%s", prefix, name)
	//tile38 command
	c := tile38Pool.Get()
	defer c.Close()
	revel.TRACE.Println("SET HOOK", group, geojson)
	ret, err := c.Do("SETHOOK", pname, fmt.Sprintf(redisHookTemplate, redisAddress, pname), "WITHIN", group, "FENCE", "OBJECT", geojson)
	if err != nil {
		revel.TRACE.Printf("%v -- %v", ret, err)
		return err
	}

	return nil
}

// DelFenceHook delete webhook to redis
func DelFenceHook(name string) error {
	c := tile38Pool.Get()
	defer c.Close()
	ret, err := c.Do("DELHOOK", name)
	if err != nil {
		revel.TRACE.Printf("%d -- %v", ret, err)
		return err
	}
	return nil
}

// DelAllHook delete all webhook to redis
func DelAllHook() error {
	c := tile38Pool.Get()
	defer c.Close()
	ret, err := c.Do("PDELHOOK", "*")
	if err != nil {
		revel.TRACE.Printf("%d -- %v", ret, err)
		return err
	}

	return nil
}

func SetTsObject(name string) error {
	c := redispool.Get()
	defer c.Close()
	_, err := c.Do("SET", fmt.Sprintf(timestampObject, name), time.Now().UnixNano())

	return err
}

func GetTsObject(name string) (int64, error) {
	c := redispool.Get()
	defer c.Close()
	ts, err := redis.Int64(c.Do("GET", fmt.Sprintf(timestampObject, name)))

	if err != nil {
		return 0, err
	}

	return ts, nil
}
