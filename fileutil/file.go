package fileutil

import (
	"encoding/json"
	"fmt"
	"github.com/logxxx/utils/log"
	"io"
	"io/fs"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
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

func FindFile(rootPath string, checkDirFn func(string) bool, checkFileFn func(filepath string) bool) (string, error) {
	subFiles, err := ioutil.ReadDir(rootPath)
	if err != nil {
		return "", err
	}

	if len(subFiles) <= 0 {
		return "", nil
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(subFiles), func(i, j int) {
		subFiles[i], subFiles[j] = subFiles[j], subFiles[i]
	})

	for _, subFile := range subFiles {
		subFilePath := filepath.Join(rootPath, subFile.Name())
		if subFile.IsDir() {

			if checkDirFn != nil && !checkDirFn(subFilePath) {
				continue
			}

			fileSubDirResult, err := FindFile(subFilePath, checkDirFn, checkFileFn)
			if err != nil {
				return "", err
			}
			if fileSubDirResult != "" {
				return fileSubDirResult, nil
			}
		} else {
			ok := checkFileFn(subFilePath)
			if ok {
				return subFilePath, nil
			}
		}
	}

	return "", nil

}

func CopyDir(srcDir, dstDir string) error {
	err := filepath.Walk(srcDir, func(srcPath string, f fs.FileInfo, err error) error {
		if err != nil {
			log.Errorf("CopyDir Walk err:%v srcPath:%v", err, srcPath)
			return err
		}

		//rebase path
		relPath, err := filepath.Rel(srcDir, srcPath)
		if err != nil {
			log.Errorf("CopyDir Walk Rel err:%v srcDir:%v srcPath:%v", err, srcDir, srcPath)
			return err
		}

		dstPath := filepath.Join(dstDir, relPath)

		if IsDir(srcPath) {
			//log.Debugf("CopyDir Mkdir:%v", dstPath)
			err := os.MkdirAll(dstPath, f.Mode())
			if err != nil {
				log.Errorf("CopyDir Walk Mkdir err:%v dstPath:%v mode:%v", err, dstPath, f.Mode())
				return err
			}
		} else {
			//log.Debugf("CopyDir CopyFile %v => %v", srcPath, dstPath)
			err := CopyFile(srcPath, dstPath, f.Mode())
			if err != nil {
				log.Errorf("CopyDir Walk CopyFile err:%v srcPath:%v dstPath:%v", err, dstPath, srcPath, dstPath)
				return err
			}
		}

		return nil

	})

	return err
}

func CopyFile(srcPath, dstPath string, pem fs.FileMode) error {
	log.Debugf("copy file from %s => %s", srcPath, dstPath)

	srcFileStat, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	if !srcFileStat.Mode().IsRegular() {
		return fmt.Errorf("not a regular file:%v", srcPath)
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	os.MkdirAll(filepath.Dir(dstPath), pem)

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	err = dstFile.Chmod(pem)
	if err != nil {
		return err
	}

	return nil
}

func IsDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	return fileInfo.IsDir()
}
