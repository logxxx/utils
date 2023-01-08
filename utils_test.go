package utils

import "testing"

func TestGetShowSize(t *testing.T) {

	req1 := int64(100 * 1024 * 1024)
	resp1 := GetShowSize(req1)
	t.Logf("resp1:%v", resp1)

}
