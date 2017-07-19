package geofence

import (
	"fmt"
	"time"

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

		ct := tile38Pool.Get()
		defer ct.Close()

		psc := redis.PubSubConn{Conn: c}
		psc.PSubscribe("fence.*")

		// experiment keep alive redis
		go func() {
			for {
				c.Do("PING")
				ct.Do("PING")
				time.Sleep(20 * time.Second)
			}
		}()

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
func SetObject(name, group string, lat, long float64) error {
	c := tile38Pool.Get()
	defer c.Close()
	fmt.Printf("SET %s %s POINT %f %f", group, name, lat, long)
	ret, err := c.Do("SET", group, name, "POINT", lat, long)
	if err != nil {
		fmt.Printf("%v -- %v", ret, err)
		return err
	}

	return nil
}

// SetFenceHook set webhook to redis
func SetFenceHook(name, group, geojson, redisAddress string) error {
	if redisAddress == "" {
		redisAddress = ":9851"
	}

	c := tile38Pool.Get()
	defer c.Close()
	fmt.Println("SET HOOK", group, geojson)
	ret, err := c.Do("SETHOOK", name, fmt.Sprintf(redisHookTemplate, redisAddress, name), "WITHIN", group, "FENCE", "OBJECT", geojson)
	if err != nil {
		fmt.Printf("%v -- %v", ret, err)
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
		fmt.Printf("%d -- %v", ret, err)
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
		fmt.Printf("%d -- %v", ret, err)
		return err
	}

	return nil
}
