package database

import (
	"encoding/json"
	"fmt"
	"os"
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

func WriteFile(path string, comicMap ComicMap) error {
	file, err := json.MarshalIndent(comicMap, "", " ")
	if err != nil {
		return err
	}
	err = os.WriteFile(path, file, 0644)
	return err
}

func ReadFile(path string) (cm ComicMap, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &cm)
	return
}
