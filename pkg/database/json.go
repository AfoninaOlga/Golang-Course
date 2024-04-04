package database

import (
	"encoding/json"
	"os"
	"strconv"
)

type Comic struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

type JsonDatabase struct {
	comics map[int]Comic
	path   string
	maxId  int
}

func (jb *JsonDatabase) Init(path string) (err error) {
	jb.path = path
	jb.maxId = getMaxId(path)
	if fileExists(jb.path) {
		var cm map[int]Comic
		var data []byte
		data, err = os.ReadFile(jb.path)
		if err != nil {
			return
		}
		err = json.Unmarshal(data, &cm)
		jb.comics = cm
	} else {
		jb.comics = map[int]Comic{}
	}
	return
}

func (jb *JsonDatabase) Flush() (err error) {
	if jb.maxId > getMaxId(jb.path) {
		//write maxId to file
		err = os.WriteFile(jb.path+".max", []byte(strconv.Itoa(jb.maxId)), 0644)
		if err != nil {
			return err
		}
		var file []byte
		file, err = json.MarshalIndent(jb.comics, "", " ")
		if err != nil {
			return err
		}
		err = os.WriteFile(jb.path, file, 0644)
		return
	}
	return
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

func (jb *JsonDatabase) GetMaxId() int {
	return jb.maxId
}

func getMaxId(dbPath string) int {
	idPath := dbPath + ".max"
	if fileExists(idPath) {
		f, err := os.ReadFile(idPath)
		if err != nil {
			return 0
		}
		res, err := strconv.Atoi(string(f))
		if err != nil {
			return 0
		}
		return res
	} else {
		return 0
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
