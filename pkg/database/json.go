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

type ComicMap map[string]Comic

func DisplayComicMap(cm ComicMap) {
	for key, value := range cm {
		fmt.Println(key, ":")
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
