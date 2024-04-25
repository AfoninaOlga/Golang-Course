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
	index  map[string][]int
}

const indexPath = "index.json"

func New(path string) (*JsonDatabase, error) {
	var jb JsonDatabase
	err := jb.init(path)
	return &jb, err
}

func (jb *JsonDatabase) init(path string) error {
	jb.path = path
	jb.comics = make(map[int]Comic)
	jb.index = make(map[string][]int)
	jb.maxId = 0
	jb.mtx = &sync.Mutex{}

	if fileExists(jb.path) {
		data, err := os.ReadFile(jb.path)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(data, &jb.comics); err != nil {
			return err
		}

		for id := range jb.comics {
			if id > jb.maxId {
				jb.maxId = id
			}
		}
	}
	if fileExists(indexPath) {
		data, err := os.ReadFile(indexPath)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(data, &jb.index); err != nil {
			return err
		}
	} else {
		for id := range jb.comics {
			for _, k := range jb.comics[id].Keywords {
				jb.index[k] = append(jb.index[k], id)
			}
		}
	}
	return nil
}

func (jb *JsonDatabase) flush() error {
	file, err := json.MarshalIndent(jb.comics, "", " ")
	if err != nil {
		return err
	}
	if err = os.WriteFile(jb.path, file, 0644); err != nil {
		return err
	}

	file, err = json.MarshalIndent(jb.index, "", " ")
	if err != nil {
		return err
	}
	return os.WriteFile(indexPath, file, 0644)
}

func (jb *JsonDatabase) Flush() error {
	jb.mtx.Lock()
	defer jb.mtx.Unlock()
	return jb.flush()
}

func (jb *JsonDatabase) GetAll() map[int]Comic {
	return jb.comics
}

func (jb *JsonDatabase) GetIndex() map[string][]int {
	return jb.index
}

func (jb *JsonDatabase) addComic(id int, c Comic) error {
	jb.comics[id] = c
	for _, k := range c.Keywords {
		jb.index[k] = append(jb.index[k], id)
	}

	if id > jb.maxId {
		jb.maxId = id
	}

	if len(jb.comics)%50 == 0 {
		return jb.flush()
	}
	return nil
}

func (jb *JsonDatabase) AddComic(id int, c Comic) error {
	jb.mtx.Lock()
	defer jb.mtx.Unlock()
	return jb.addComic(id, c)
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

func (jb *JsonDatabase) Exists(id int) bool {
	jb.mtx.Lock()
	defer jb.mtx.Unlock()
	_, ok := jb.comics[id]
	return ok
}

func (jb *JsonDatabase) Size() int {
	return len(jb.comics)
}
