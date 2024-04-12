package database

import (
	"encoding/json"
	"os"
	"sync"
)

type Comic struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

type JsonDatabase struct {
	comics map[int]Comic
	path   string
	maxId  int
	mtx    *sync.Mutex
}

func New(path string) (JsonDatabase, error) {
	var jb JsonDatabase
	err := jb.init(path)
	return jb, err
}

func (jb *JsonDatabase) init(path string) error {
	jb.path = path
	jb.comics = map[int]Comic{}
	jb.maxId = 0
	jb.mtx = &sync.Mutex{}

	if fileExists(jb.path) {
		var cm map[int]Comic
		data, err := os.ReadFile(jb.path)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(data, &cm); err != nil {
			return err
		}
		jb.comics = cm
		for id := range cm {
			if id > jb.maxId {
				jb.maxId = id
			}
		}
	}
	return nil
}

func (jb *JsonDatabase) Flush() (err error) {
	var file []byte
	file, err = json.MarshalIndent(jb.comics, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(jb.path, file, 0644)
	return
}

func (jb *JsonDatabase) FlushParallel() error {
	jb.mtx.Lock()
	defer jb.mtx.Unlock()
	err := jb.Flush()
	return err
}

func (jb *JsonDatabase) GetAll() map[int]Comic {
	return jb.comics
}

func (jb *JsonDatabase) AddComic(id int, c Comic) {
	jb.comics[id] = c
	if id > jb.maxId {
		jb.maxId = id
	}
}

func (jb *JsonDatabase) AddComicParallel(id int, c Comic) {
	jb.mtx.Lock()
	jb.AddComic(id, c)
	jb.mtx.Unlock()
}

func (jb *JsonDatabase) GetMaxId() int {
	return jb.maxId
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (jb *JsonDatabase) GetMissingIds() []int {
	var ids []int
	for i := 1; i < jb.maxId; i++ {
		if _, ok := jb.comics[i]; !ok && i != 404 {
			ids = append(ids, i)
		}
	}
	return ids
}
