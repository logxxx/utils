package filequeue

import (
	"encoding/json"
	"errors"
	"github.com/logxxx/utils/fileutil"
	"github.com/logxxx/utils/log"
	"io/ioutil"
	"os"
	"strings"
	"sync"
)

var (
	_fileQueueMap = make(map[string]*FileQueue, 0)

	ErrEmpty = errors.New("empty queue")
)

type FileQueue struct {
	lock     sync.RWMutex
	fileName string
}

func GetFileQueue(fileName string) *FileQueue {
	q := _fileQueueMap[fileName]
	if q != nil {
		return q
	}
	_fileQueueMap[fileName] = NewFileQueue(fileName)

	return _fileQueueMap[fileName]
}

func NewFileQueue(fileName string) *FileQueue {
	return &FileQueue{
		lock:     sync.RWMutex{},
		fileName: fileName,
	}
}

func saveToFile(fileName string, obj []json.RawMessage) error {

	data, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	newData := strings.ReplaceAll(string(data), "},", "},\n") //便于查看

	err = ioutil.WriteFile(fileName, []byte(newData), 0766)
	if err != nil {
		return err
	}

	return nil
}

func readFile(fileName string) ([]json.RawMessage, error) {

	if !fileutil.HasFile(fileName) {
		return nil, nil
	}

	allRows := make([]json.RawMessage, 0)

	fileData, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Errorf("FileQueue.Push ioutil.ReadFile err:%v fileName:%v", err, fileName)
		return nil, err
	}

	if len(fileData) <= 0 {
		return nil, nil
	}

	err = json.Unmarshal(fileData, &allRows)
	if err != nil {
		log.Errorf("FileQueue.Push Unmarshal err:%v fileName:%v", err, fileName)
		return nil, err
	}

	return allRows, nil

}

func (q *FileQueue) Clean() {

	q.lock.Lock()
	defer q.lock.Unlock()

	os.Remove(q.fileName)
}

func (q *FileQueue) Push(obj interface{}) error {

	q.lock.Lock()
	defer q.lock.Unlock()

	rows, err := readFile(q.fileName)
	if err != nil {
		log.Errorf("Push readFile err:%v", err)
		return err
	}

	newRow, err := json.Marshal(obj)
	if err != nil {
		log.Errorf("Push Marshal err:%v", err)
		return err
	}

	rows = append(rows, newRow)

	err = saveToFile(q.fileName, rows)
	if err != nil {
		log.Errorf("Push saveToFile err:%v", err)
		return err
	}

	return nil

}

func (q *FileQueue) Pop(obj interface{}) error {
	err := q.MustPop(obj)
	if err != nil && err != ErrEmpty {
		return err
	}
	return nil
}

func (q *FileQueue) MustPop(obj interface{}) error {

	q.lock.Lock()
	defer q.lock.Unlock()

	rows, err := readFile(q.fileName)
	if err != nil {
		log.Errorf("Pop readFile err:%v", err)
		return err
	}

	if len(rows) == 0 {
		return ErrEmpty
	}

	err = json.Unmarshal(rows[0], obj)
	if err != nil {
		log.Errorf("Pop Unmarshal err:%v", err)
		return err
	}

	newRows := rows[1:]
	err = saveToFile(q.fileName, newRows)
	if err != nil {
		log.Errorf("Pop saveToFile err:%v", err)
		return err
	}

	return nil
}
