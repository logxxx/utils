package exist

import (
	"os"
	"path/filepath"
	"strings"
)

type Exister struct {
	FilePath  string
	cache     map[string]bool
	writeStep int
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
	if e.cache[key] {
		return true
	}
	e.cache[key] = true
	e.writeFile()
	return false
}

func (e *Exister) writeFile() {
	os.MkdirAll(filepath.Dir(e.FilePath), 0755)
	content := []string{}
	for key := range e.cache {
		content = append(content, key)
	}
	os.WriteFile(e.FilePath, []byte(strings.Join(content, "\n")), 0766)
}
