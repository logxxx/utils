package exister

import (
	"fmt"
	"github.com/logxxx/utils/fileutil"
	"os"
	"path/filepath"
	"sync"
)

type Exister struct {
	lock     sync.RWMutex
	filepath string
}

func NewExister(filePath string) *Exister {
	os.MkdirAll(filepath.Dir(filePath), 0755)
	return &Exister{
		filepath: filePath,
	}
}

func (ex *Exister) Set(key interface{}) {
	ex.lock.Lock()
	defer ex.lock.Unlock()
	data := ex.getData()
	data[fmt.Sprintf("%v", key)] = true
	ex.save(data)
}

func (ex *Exister) IsExist(key interface{}) bool {
	ex.lock.RLock()
	defer ex.lock.RUnlock()
	data := ex.getData()
	return data[fmt.Sprintf("%v", key)]
}

func (ex *Exister) save(data map[string]bool) {
	fileutil.WriteJsonToFile(data, ex.filepath)
}

func (ex *Exister) Delete(key interface{}) {
	ex.lock.Lock()
	defer ex.lock.Unlock()
	data := ex.getData()
	delete(data, fmt.Sprintf("%v", key))
	ex.save(data)
}

func (ex *Exister) Clean() {
	ex.save(map[string]bool{})
}

func (ex *Exister) getData() map[string]bool {
	data := make(map[string]bool, 0)
	if fileutil.HasFile(ex.filepath) {
		fileutil.ReadJsonFile(ex.filepath, &data)
	}
	return data
}
