package xkcd

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type UrlComic struct {
	Id         uint   `json:"num"`
	Url        string `json:"img"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
}

type Client struct {
	Url    string
	Client http.Client
}

func NewClient(url string, timeout time.Duration) Client {
	c := http.Client{Timeout: timeout}
	return Client{url, c}
}

func (c Client) GetComic(id int) (UrlComic, error) {
	return getComic(c, "/"+strconv.Itoa(id)+"/info.0.json")
}

func (c Client) GetComicsCount() (uint, error) {
	comic, err := getComic(c, "/info.0.json")
	if err != nil {
		return 0, err
	}
	return comic.Id, nil
}

func getComic(c Client, suffix string) (UrlComic, error) {
	resp, err := c.Client.Get(c.Url + suffix)
	if err != nil {
		return UrlComic{}, err
	}

	var comic UrlComic
	err = json.NewDecoder(resp.Body).Decode(&comic)
	return comic, err
}
