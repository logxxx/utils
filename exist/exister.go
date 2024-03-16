package exist

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Exister struct {
	FilePath  string
	cacheLock sync.Mutex
	cache     map[string]bool
}

func NewExister(filePath string) *Exister {
	resp := &Exister{
		FilePath: filePath,
		cache:    map[string]bool{},
	}

	resp.loadFile()

	return resp

}

func (e *Exister) loadFile() {
	file, err := os.ReadFile(e.FilePath)
	if err != nil {
		return
	}
	for _, key := range strings.Split(string(file), "\n") {
		if key == "" {
			continue
		}
		e.cache[key] = true
	}
}

func (e *Exister) Has(key string) bool {
	if key == "" {
		return false
	}
	e.cacheLock.Lock()
	defer e.cacheLock.Unlock()
	if e.cache[key] {
		return true
	}
	e.cache[key] = true
	e.writeFile(key)
	return false
}

func (e *Exister) writeFile(key string) {
	os.MkdirAll(filepath.Dir(e.FilePath), 0755)

	f, err := os.OpenFile(e.FilePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	f.WriteString(key + "\n")
}
