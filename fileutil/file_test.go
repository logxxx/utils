package fileutil

import (
	"fmt"
	"github.com/logxxx/utils"
	log "github.com/sirupsen/logrus"
	rand2 "math/rand"
	"os"
	"strings"
	"testing"
	"time"
)

func TestWriteToFileWithRename(t *testing.T) {
	for i := 0; i < 10; i++ {
		newPath, err := WriteToFileWithRename([]byte("hello"), "./download/1/test.jpg")
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%v newPath:%v", i, newPath)
	}

}

func TestFindFile(t *testing.T) {
	result, err := FindFile("N:\\source", func(filepath string) bool {
		if strings.Contains(filepath, "微信录屏") {
			return false
		}
		return true
	}, func(filepath string) bool {

		if !strings.Contains(filepath, ".jpg") {
			return false
		}
		stat, err := os.Stat(filepath)
		if err != nil {
			return false
		}
		log.Infof("path:%v size:%v", filepath, utils.GetShowSize(stat.Size()))
		if stat.Size() < 1*1024*1024 {
			return false
		}
		return true
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("result:%v", result)
}

func TestGetUniqFilePath(t *testing.T) {
	filePath := "./test.txt"
	WriteToFile([]byte("hehe"), filePath)
	filePath2 := GetUniqFilePath(filePath)
	t.Logf("filePath2:%v", filePath2)
	WriteToFile([]byte("haha"), filePath2)
	filePath3 := GetUniqFilePath(filePath)
	t.Logf("filePath3:%v", filePath3)
}

func TestAppendToFile(t *testing.T) {
	fileNames := []string{
		fmt.Sprintf("%v.txt", time.Now().UnixNano()),
		fmt.Sprintf("test/%v.txt", time.Now().UnixNano()),
	}

	for _, fileName := range fileNames {
		defer func() {
			os.RemoveAll(fileName)
		}()
		rand := fmt.Sprintf("%v%v", time.Now().UnixNano(), rand2.Int63())
		have := []string{"hello", "xxx", " world", rand}
		want := "helloxxx world" + rand
		for _, e := range have {
			err := AppendToFile(fileName, e)
			if err != nil {
				t.Fatal(err)
			}
		}
		result, err := os.ReadFile(fileName)
		if err != nil {
			t.Fatal(err)
		}
		if string(result) != want {
			t.Fatalf("result:%v want:%v", string(result), want)
		}
	}

}
