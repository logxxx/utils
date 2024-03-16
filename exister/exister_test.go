package exister

import (
	"testing"
	"time"
)

func TestNewExister(t *testing.T) {
	ex := NewExister("test_exister")
	req1 := time.Now().UnixNano()
	resp1 := ex.IsExist(req1)
	if resp1 != false {
		panic("")
	}

	ex.Set(req1)
	resp2 := ex.IsExist(req1)
	if resp2 != true {
		panic("")
	}

	ex.Delete(req1)
	resp3 := ex.IsExist(req1)
	if resp3 != false {
		panic("")
	}

	ex.Set(req1)
	resp4 := ex.IsExist(req1)
	if resp4 != true {
		panic("")
	}

	req5 := time.Now().UnixNano()
	resp5 := ex.IsExist(req5)
	if resp5 != false {
		panic("")
	}

	ex.Set(req5)

	resp6 := ex.IsExist(req5)
	if resp6 != true {
		panic("")
	}

	//ex.Clean()
}
