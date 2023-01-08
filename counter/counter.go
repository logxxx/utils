package counter

import (
	"fmt"
	"gopkg.in/redis.v3"
	"time"
)

type Counter struct {
	*redis.Client
}

func (c *Counter) IsExist(req interface{}) bool {
	ok, err := c.SetNX(fmt.Sprintf("%v", req), time.Now().Unix(), time.Hour*24*7).Result()
	if err != nil {
		return false
	}
	return !ok //设置成功->不存在->exist=false
}

func EmptyCounter() *Counter {
	return &Counter{
		Client: redis.NewClient(&redis.Options{}),
	}
}

func NewCounter(addr, pwd string, db int64) (*Counter, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
	})

	testPing := rdb.Ping().String()
	if testPing != "PING: PONG" {
		return nil, fmt.Errorf("ping failed:%v", testPing)
	}

	counter := &Counter{Client: rdb}

	return counter, nil

}
