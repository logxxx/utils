package fileutil

import "testing"

func TestWriteToFileWithRename(t *testing.T) {
	for i := 0; i < 10; i++ {
		err := WriteToFileWithRename([]byte("hello"), "./download/1", "test.jpg")
		if err != nil {
			t.Fatal(err)
		}
	}

}
