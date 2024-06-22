package handler

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"net/url"
)

func (h *Handler) Search(w http.ResponseWriter, req *http.Request) {
	token, err := req.Cookie("token")
	if err != nil || token == nil {
		log.Println("error getting token from cookies:", err, "redirecting to login")
		http.Redirect(w, req, "/login", http.StatusSeeOther)
	}

	tmpl, err := template.ParseFiles("webserver/internal/handler/templates/comics.html")
	if err != nil {
		log.Println("error parsing comics template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

	if !req.URL.Query().Has("search") {
		log.Println("rendering search page")
		tmpl.Execute(w, nil)
		return
	}

	var data struct {
		Images     []string
		ImageCount int
		ErrMessage string
	}

	text := req.URL.Query().Get("search")

	if text == "" {
		data.ErrMessage = "Search query should not be empty. Enter keyword(s) and try again."
		log.Println("Empty search query")
		tmpl.Execute(w, data)
		return
	}

	text = url.QueryEscape(text)
	apiReq, err := http.NewRequest("GET", h.apiUrl+"/pics?search="+text, nil)
	if err != nil {
		log.Println("error creating search request:", err)
		data.ErrMessage = "Internal error. Try again later."
		tmpl.Execute(w, data)
		return

	}
	apiReq.Header.Set("Authorization", token.Value)

	resp, err := h.client.Do(apiReq)
	if err != nil {
		log.Println("error sending search request:", err)
		tmpl.Execute(w, data)
		return
	}
	if resp.StatusCode == http.StatusUnauthorized {
		log.Println("unauthorized status code with token:", token)
		http.Redirect(w, req, "/login", http.StatusSeeOther)
	}
	if resp.StatusCode != http.StatusOK {
		data.ErrMessage = "Internal error. Try again later."
		tmpl.Execute(w, data)
		return
	}

	var searchResult []struct {
		Url string `json:"url"`
	}
	err = json.NewDecoder(resp.Body).Decode(&searchResult)
	if err != nil {
		log.Println("error decoding search response:", err)
		data.ErrMessage = "Internal error. Try again later."
		tmpl.Execute(w, data)
		return
	}

	for _, v := range searchResult {
		log.Println(v)
		data.Images = append(data.Images, v.Url)
	}
	data.ImageCount = len(data.Images)
	log.Println(data.Images, data.ImageCount)

	if data.ImageCount == 0 {
		data.ErrMessage = "No comics found. Try other keywords."
	}

	tmpl.Execute(w, data)
	return
}
