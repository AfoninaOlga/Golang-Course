package xkcd

import (
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type UrlComic struct {
	Id         uint   `json:"num"`
	Url        string `json:"img"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
}

func GetComicResponse(url string) (comic UrlComic, err error) {
	c := http.Client{Timeout: time.Duration(1) * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Add("Accept", `application/json`)
	resp, err := c.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &comic)
	return
}

func GetComicsCount(url string) (cnt uint, err error) {
	comic, err := GetComicResponse(url)
	if err != nil {
		return
	}
	cnt = comic.Id
	return
}
