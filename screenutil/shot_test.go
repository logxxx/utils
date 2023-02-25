package screen

import (
	"io/ioutil"
	"testing"
)

func TestDrawScreenOnly(t *testing.T) {
	resp, err := ShotRect(1167, 242, 1230, 312)
	if err != nil {
		t.Fatal(err)
	}
	ioutil.WriteFile("./resp1.jpg", resp, 0777)
}
