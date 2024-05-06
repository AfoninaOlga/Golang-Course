package domain

type Comic struct {
	Url      string
	Keywords []string
}

type UrlComic struct {
	Id         int    `json:"num"`
	Url        string `json:"img"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Title      string `json:"title"`
}
