package filekv

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/logxxx/utils"
	"github.com/logxxx/utils/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
)

var (
	ErrNotFound = errors.New("not found")
)

// FileKV 以文件的方式存储kv。
// 适用于读多写少(效率不高)，只有getset的场景。
type FileKV struct {
	fileName  string
	lock      sync.RWMutex
	encryptFn func([]byte) ([]byte, error)
	decryptFn func([]byte) ([]byte, error)
}

func NewFileKV(fileName string) *FileKV {
	dir := filepath.Dir(fileName)
	if dir != "" && dir != "." {
		os.MkdirAll(dir, 0755)
	}

	return &FileKV{
		fileName: fileName,
		encryptFn: func(data []byte) ([]byte, error) {
			return data, nil
		},
		decryptFn: func(data []byte) ([]byte, error) {
			return data, nil
		},
	}
}

func (w *FileKV) SetEncryptFn(encryptFn func(data []byte) ([]byte, error)) {
	w.encryptFn = encryptFn
}

func (w *FileKV) SetDecryptFn(decryptFn func(data []byte) ([]byte, error)) {
	w.decryptFn = decryptFn
}

func (w *FileKV) MustGet(key string, value interface{}) error {

	w.lock.RLock()
	defer w.lock.RUnlock()

	data, err := w.readFileAndDecrypt()
	if err != nil {
		log.Errorf("FileKV.MustGet readFileAndDecrypt err:%v fileName:%v", err, w.fileName)
		return err
	}

	rawValue, ok := data[key]
	if !ok || len(rawValue) < 0 || string(rawValue) == "{}" {
		return ErrNotFound
	}

	err = json.Unmarshal(rawValue, value)
	if err != nil {
		log.Errorf("FileKV.MustGet Unmarshal err:%v rawValue:%v", err, string(rawValue))
		return err
	}

	return nil
}

func (w *FileKV) GetStr(key string) string {
	str := ""
	w.Get(key, &str)
	return str
}

func (w *FileKV) Get(key string, value interface{}) error {

	err := w.MustGet(key, value)
	if err != nil {
		if err == ErrNotFound {
			return nil
		}
		return err
	}
	return nil
}

func (w *FileKV) CleanAll() error {
	w.lock.Lock()
	defer w.lock.Unlock()
	//log.Debugf("FileKV.RemoveFile:%v", w.fileName)
	if !utils.HasFile(w.fileName) {
		return nil
	}
	return os.Remove(w.fileName)

}

func (w *FileKV) Set(key string, value interface{}) error {
	return w.setWithLock(key, value)
}

func (w *FileKV) setWithLock(key string, value interface{}) error {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.set(key, value)
}

func (w *FileKV) set(key string, value interface{}) error {

	//log.Debugf("FileKV.Set start. fileName:%v key:%v value:%+v", w.fileName, key, utils.JsonToString(value))

	data, err := w.readFileAndDecrypt()
	if err != nil {
		log.Errorf("FileKV.Set readFileAndDecrypt err:%v fileName:%v", err, w.fileName)
		return err
	}

	var valueBytes = []byte("{}")
	if value != nil {
		valueBytes, err = json.Marshal(value)
		if err != nil {
			log.Errorf("FileKV.Set json.Marshal err:%v value:%+v", err, utils.JsonToString(value))
			return err
		}
	}

	if bytes.Compare(data[key], valueBytes) == 0 {
		//log.Debugf("FileKV.Set NO NEED! data[%v]:%v value:%v", key, string(data[key]), string(valueBytes))
		return nil
	}

	data[key] = valueBytes

	newData, err := json.Marshal(data)
	if err != nil {
		log.Errorf("FileKV.Set json.Marshal data err:%v data:%+v", err, utils.JsonToString(data))
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

	err = ioutil.WriteFile(w.fileName, newData, 0766)
	if err != nil {
		log.Errorf("FileKV.Set ioutil.WriteFile err:%v fileName:%v", err, w.fileName)
		return err
	}

	return nil
}

func (w *FileKV) readFileAndDecrypt() (allKVs map[string][]byte, err error) {

	allKVs = make(map[string][]byte)

	if !utils.HasFile(w.fileName) {
		return
	}

	fileData, err := ioutil.ReadFile(w.fileName)
	if err != nil {
		log.Errorf("FileKV.readFileAndDecrypt ioutil.ReadFile err:%v fileName:%v", err, w.fileName)
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
			err = nil //解析失败，屏蔽错误,直接返回空的map。以保证可用性
			return
		}
	}

	//log.Debugf("FileKV.readFileAndDecrypt fileName:%v fileData:%v fileData.len:%v", w.fileName, string(fileData), len(fileData))

	if len(fileData) <= 0 {
		return
	}

	err = json.Unmarshal(fileData, &allKVs)
	if err != nil {
		log.Errorf("FileKV.readFileAndDecrypt Unmarshal err:%v fileName:%v fileData:%v", err, w.fileName, string(fileData))
		err = nil //解析失败，屏蔽错误,直接返回空的map。以保证可用性
		return
	}

	return
}
