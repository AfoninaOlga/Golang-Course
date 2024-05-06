package domain

type Comic struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

type UrlComic struct {
	Id         int    `json:"num"`
	Url        string `json:"img"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Title      string `json:"title"`
}

type FoundComic struct {
	Id    int
	Url   string
	Count int
}
