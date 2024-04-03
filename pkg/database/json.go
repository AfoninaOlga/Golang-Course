package database

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type Comic struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

type ComicMap map[int]Comic

func DisplayComicMap(cm ComicMap, cnt int) {
	for i := 1; i <= cnt; i++ {
		value := cm[i]
		fmt.Println(i, ":")
		fmt.Println("\turl:", value.Url)
		fmt.Println("\tkeywords:", value.Keywords)
	}
}

func WriteFile(path string, comicMap ComicMap, maxId int) error {
	if maxId > GetMaxIdFromDB(path) {
		//write maxId to file
		os.WriteFile(path+".max", []byte(strconv.Itoa(maxId)), 0644)
	}
	file, err := json.MarshalIndent(comicMap, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, file, 0644)
	return err
}

func ReadFile(path string) (cm ComicMap, err error) {
	if !fileExists(path) {
		return ComicMap{}, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &cm)
	return
}

func GetMaxIdFromDB(path string) int {
	path += ".max"
	if fileExists(path) {
		f, err := os.ReadFile(path)
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
