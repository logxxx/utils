package fileutil

import (
	"encoding/json"
	"fmt"
	"github.com/logxxx/utils/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func IsExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil && !os.IsExist(err) {
		return false
	}
	return true
}

func WriteToFileWithRename(data []byte, dir, fileName string) (string, error) {
	dir, fileName = getValidPath(dir, fileName)
	newPath := filepath.Join(dir, fileName)
	return newPath, WriteToFile(data, newPath)
}

func getValidPath(dir, fileNameWithExt string) (string, string) {

	if !HasFile(filepath.Join(dir, fileNameWithExt)) {
		return dir, fileNameWithExt
	}

	i := 0
	fileExt := filepath.Ext(fileNameWithExt)
	fileName := strings.TrimRight(fileNameWithExt, fileExt)
	for {
		i++
		fileNameWithExt = fmt.Sprintf("%v_%v%v", fileName, i, fileExt)

		if !HasFile(filepath.Join(dir, fileNameWithExt)) {
			return dir, fileNameWithExt
		}

	}

}

func WriteToFile(data []byte, filePath string) error {
	fileDir := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)
	file, _, err := GetOrCreateFile(fileDir, fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func GetOrCreateFile(fileDir, fileName string) (*os.File, int64, error) {

	filePath := filepath.Join(fileDir, fileName)

	if HasFile(filePath) {
		os.Remove(filePath)
	}

	err := os.MkdirAll(fileDir, 0755)
	if err != nil {
		log.Errorf("getOrCreateFile MkdirAll err:%v dir:%v", err, fileDir)
		return nil, 0, err
	}

	file, err := os.Create(filePath)
	if err != nil {
		log.Errorf("getOrCreateFile os.Create err:%v path:%v", err, filePath)
		return nil, 0, err
	}

	return file, 0, nil

}

func HasFile(path string) bool {
	if _, err := os.Stat(path); err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func ReadJsonFile(fileName string, obj interface{}) error {
	content, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	err = json.Unmarshal(content, obj)
	if err != nil {
		return err
	}
	return nil
}

func WriteJsonToFile(obj interface{}, filePath string) error {
	data, err := json.MarshalIndent(obj, " ", "")
	if err != nil {
		return err
	}

	err = WriteToFile(data, filePath)
	if err != nil {
		return err
	}

	return nil
}
