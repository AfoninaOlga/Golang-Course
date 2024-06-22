package handler

import (
	"bytes"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
)

type FoundComic struct {
	Id    int    `json:"-"`
	Url   string `json:"url"`
	Count int    `json:"-"`
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type Handler struct {
	client http.Client
	apiUrl string
	ttl    time.Duration
}

func NewHandler(apiUrl string, ttl time.Duration) *Handler {
	return &Handler{ttl: ttl, apiUrl: apiUrl, client: http.Client{Timeout: 5 * time.Second}}
}

func (h *Handler) Login(w http.ResponseWriter, req *http.Request) {
	var data struct {
		ErrMessage string
	}
	tmpl, err := template.ParseFiles("webserver/internal/handler/templates/login.html")
	if err != nil {
		log.Println("Error parsing login template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
	if req.Method == http.MethodGet {
		tmpl.Execute(w, data)
		return
	}

	if req.Method == http.MethodPost {
		err := req.ParseForm()
		if err != nil {
			log.Println("error parsing login form:", err)
			data.ErrMessage = "Invalid input data. Try again."
			tmpl.Execute(w, data)
			return
		}

		name := req.Form.Get("name")
		password := req.Form.Get("password")
		user := []byte(`{"name":"` + name + `", "password":"` + password + `"}`)

		resp, err := http.Post(h.apiUrl+"/login", "application/json", bytes.NewBuffer(user))
		if err != nil {
			log.Printf("error getting token: %w", err)
			data.ErrMessage = "Incorrect username or password. Try again."
			tmpl.Execute(w, data)
			return
		}

		if resp.StatusCode == http.StatusOK {
			var token struct {
				Token string `json:"token"`
			}
			defer resp.Body.Close()
			err = json.NewDecoder(resp.Body).Decode(&token)
			if err != nil {
				log.Printf("error decoding token: %v\n", err)
				data.ErrMessage = "Incorrect input data. Try again."
				tmpl.Execute(w, data)
				return
			}
			cookie := http.Cookie{
				Name:   "token",
				Value:  token.Token,
				MaxAge: int(h.ttl),
			}
			http.SetCookie(w, &cookie)
			http.Redirect(w, req, "/comics", http.StatusSeeOther)

		} else {
			log.Println("error requesting /login:", resp.Status)
			if resp.StatusCode == http.StatusUnauthorized {
				data.ErrMessage = "Wrong username or password. Try again."
			}
			if resp.StatusCode >= 500 {
				data.ErrMessage = "Internal error. Try again."
			}
			tmpl.Execute(w, data)
			return
		}
	}
}
