package xkcd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func NewClient(url string) Client {
	c := http.Client{}
	return Client{url, c}
}

func (c Client) GetComicResponse(id int) (comic UrlComic, err error) {
	url := fmt.Sprintf("%v/%v/info.0.json", c.Url, id)

	//hack to get comic with max id
	if id == -1 {
		url = c.Url + "/info.0.json"
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req.Header.Add("Accept", `application/json`)
	resp, err := c.Client.Do(req)
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

func (c Client) GetComicsCount() (cnt uint, err error) {
	comic, err := c.GetComicResponse(-1)
	if err != nil {
		return
	}
	cnt = comic.Id
	return
}
