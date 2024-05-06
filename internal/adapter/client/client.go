package client

import (
	"encoding/json"
	"fmt"
	"github.com/AfoninaOlga/xkcd/internal/core/domain"
	"net/http"
	"strconv"
	"time"
)

type Client struct {
	url    string
	client http.Client
}

func NewClient(url string, timeout time.Duration, connectionLimit int) *Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = connectionLimit
	t.MaxConnsPerHost = connectionLimit
	t.MaxIdleConnsPerHost = connectionLimit
	c := http.Client{Timeout: timeout, Transport: t}
	return &Client{url, c}
}

func (c Client) GetComic(id int) (domain.UrlComic, error) {
	return getComic(c, "/"+strconv.Itoa(id)+"/info.0.json")
}

func (c Client) GetComicsCount() (int, error) {
	comic, err := getComic(c, "/info.0.json")
	if err != nil {
		return 0, err
	}
	return comic.Id, nil
}

func getComic(c Client, suffix string) (domain.UrlComic, error) {
	url := c.url + suffix
	resp, err := c.client.Get(url)
	if err != nil {
		return domain.UrlComic{}, err
	}

	if resp.StatusCode != 200 {
		return domain.UrlComic{}, fmt.Errorf("Error getting %v, StatusCode=%v", url, resp.StatusCode)
	}

	var comic domain.UrlComic
	err = json.NewDecoder(resp.Body).Decode(&comic)
	return comic, err
}
