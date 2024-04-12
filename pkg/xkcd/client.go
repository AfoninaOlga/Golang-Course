package xkcd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

type UrlComic struct {
	Id         int    `json:"num"`
	Url        string `json:"img"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
}

type Client struct {
	url    string
	client http.Client
}

func NewClient(url string, timeout time.Duration) Client {
	c := http.Client{Timeout: timeout}
	return Client{url, c}
}

func (c Client) GetComic(id int) (UrlComic, error) {
	return getComic(c, "/"+strconv.Itoa(id)+"/info.0.json")
}

func (c Client) GetComicsCount() (int, error) {
	comic, err := getComic(c, "/info.0.json")
	if err != nil {
		return 0, err
	}
	return comic.Id, nil
}

func getComic(c Client, suffix string) (UrlComic, error) {
	url := c.url + suffix
	resp, err := c.client.Get(url)
	if err != nil {
		return UrlComic{}, err
	}

	if resp.StatusCode != 200 {
		return UrlComic{}, fmt.Errorf("Error getting %v, StatusCode=%v", url, resp.StatusCode)
	}

	var comic UrlComic
	err = json.NewDecoder(resp.Body).Decode(&comic)
	return comic, err
}
