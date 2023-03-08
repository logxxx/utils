package filekv

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/utils/log"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var (
	_fileKV = &FileKV{
		encryptFn: func(data []byte) ([]byte, error) { return data, nil },
		decryptFn: func(data []byte) ([]byte, error) { return data, nil },
	}
	ErrNotFound = errors.New("not found")
)

// FileKV 以文件的方式存储kv。
// 适用于读多写少(效率不高)，只有getset的场景。
type FileKV struct {
	lock      sync.RWMutex
	encryptFn func([]byte) ([]byte, error)
	decryptFn func([]byte) ([]byte, error)
}

func GetFileKV() *FileKV {
	return _fileKV
}

func (w *FileKV) SetEncryptFn(encryptFn func(data []byte) ([]byte, error)) {
	w.encryptFn = encryptFn
}

func (w *FileKV) SetDecryptFn(decryptFn func(data []byte) ([]byte, error)) {
	w.decryptFn = decryptFn
}

func (w *FileKV) MustGet(fileName, key string, value interface{}) error {

	w.lock.RLock()
	defer w.lock.RUnlock()

	data, err := w.readFileAndDecrypt(fileName)
	if err != nil {
		log.Errorf("FileKV.MustGet readFileAndDecrypt err:%v fileName:%v", err, fileName)
		return err
	}

	rawValue, ok := data[key]
	if !ok || len(rawValue) < 0 {
		return ErrNotFound
	}

	err = json.Unmarshal(rawValue, value)
	if err != nil {
		log.Errorf("FileKV.MustGet Unmarshal err:%v rawValue:%v", err, string(rawValue))
		return err
	}

	return nil
}

func (w *FileKV) GetString(fileName, key string) string {
	value := ""
	_ = w.Get(fileName, key, &value)
	return value
}

func (w *FileKV) Get(fileName, key string, value interface{}) error {

	err := w.MustGet(fileName, key, value)
	if err != nil {
		if err == ErrNotFound {
			return nil
		}
		return err
	}
	return nil
}

func (w *FileKV) RemoveFile(fileName string) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	log.Debugf("FileKV.RemoveFile:%v", fileName)
	if !utils.HasFile(fileName) {
		return nil
	}
	return os.Remove(fileName)

}

func (w *FileKV) Set(fileName, key string, value interface{}) error {
	return w.setWithLock(fileName, key, value)
}

func (w *FileKV) setWithLock(fileName, key string, value interface{}) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.set(fileName, key, value)
}

func (w *FileKV) set(fileName, key string, value interface{}) error {

	//log.Debugf("FileKV.Set start. fileName:%v key:%v value:%+v", fileName, key, utils.JsonToString(value))

	data, err := w.readFileAndDecrypt(fileName)
	if err != nil {
		log.Errorf("FileKV.Set readFileAndDecrypt err:%v fileName:%v", err, fileName)
		return err
	}

	var valueBytes = []byte("{}")
	if value != nil {
		valueBytes, err = json.Marshal(value)
		if err != nil {
			log.Errorf("FileKV.Set json.Marshal err:%v value:%+v", err, value)
			return err
		}
	}

	if bytes.Compare(data[key], valueBytes) == 0 {
		//log.Debugf("FileKV.Set NO NEED! data[%v]:%v value:%v", key, string(data[key]), string(valueBytes))
		return nil
	}

	data[key] = valueBytes
	data["update_time"] = []byte(time.Now().Format("2006/01/02 15:04:05"))

	newData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		log.Errorf("FileKV.Set json.Marshal data err:%v data:%+v", err, data)
		return err
	}

	//加密
	if w.encryptFn != nil {
		newData, err = w.encryptFn(newData)
		if err != nil {
			log.Errorf("FileKV.Set encryptFn err:%v data:%v", err, string(newData))
			return err
		}
	}

	err = fileutil.WriteToFile(newData, fileName)
	if err != nil {
		log.Errorf("FileKV.Set WriteToFile err:%v fileName:%v", err, fileName)
		return err
	}

	return nil
}

func (w *FileKV) readFileAndDecrypt(fileName string) (allKVs map[string][]byte, err error) {

	allKVs = make(map[string][]byte)

	if !utils.HasFile(fileName) {
		return
	}

	fileData, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Errorf("FileKV.readFileAndDecrypt ioutil.ReadFile err:%v fileName:%v", err, fileName)
		return
	}

	if len(fileData) <= 0 {
		return
	}

	//解密
	if w.decryptFn != nil {
		fileData, err = w.decryptFn(fileData)
		if err != nil {
			log.Errorf("FileKV.readFileAndDecrypt decryptFn err:%v fileData:%v", err, string(fileData))
			return
		}
	}

	err = json.Unmarshal(fileData, &allKVs)
	if err != nil {
		log.Errorf("FileKV.readFileAndDecrypt Unmarshal err:%v fileData:%v", err, string(fileData))
		return
	}

	return
}
