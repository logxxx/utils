package counter_test

import (
	"github.com/logxxx/mybili_pkg/utils/counter"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCounter_Exist(t *testing.T) {

	c, err := counter.NewCounter("49.232.219.233:6379", "he1234", 1)
	if err != nil {
		t.Fatal(err)
	}

	req1 := time.Now().UnixNano()
	resp1 := c.IsExist(req1)
	assert.False(t, resp1)

	resp2 := c.IsExist(req1)
	assert.True(t, resp2)

	resp3 := c.IsExist(req1)
	assert.True(t, resp3)

}

func TestEmptyCounter(t *testing.T) {

	c := counter.EmptyCounter()
	resp1 := c.IsExist(123)
	t.Logf("resp1:%v", resp1)

	err := c.Set("hello", "world", 10*time.Hour).Err()
	t.Logf("err:%v", err)

}
