package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const GITHUB_API_KEY = ""
const GITHUB_USERNAME = ""

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/login", func(w http.ResponseWriter, r *http.Request) {
		var usr User
		dec := json.NewDecoder(r.Body)
		err := dec.Decode(&usr)
		if err != nil {
			w.Write([]byte("error"))
		}

		for _, u := range users {
			if u.Username == usr.Username && u.Password == usr.Password {
				tokStr := []byte(u.Username + ":" + u.Password)
				tokEnc := base64.StdEncoding.EncodeToString(tokStr)
				fmt.Printf("user logged in, user=%s, token=%s", u.Username, tokEnc)

				w.WriteHeader(200)
				w.Write([]byte("{\"token\":\"" + tokEnc + "\"}"))
				return
			}
		}

		w.WriteHeader(404)
		w.Write([]byte("user not found"))
	})

	mux.HandleFunc("/api/list", func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			w.WriteHeader(500)
			w.Write([]byte("no token"))
			return
		}
		fmt.Printf("got token: %s", token)
		ghUsername := r.URL.Query().Get("gh-username")
		fmt.Printf("got ghUsername: %s", ghUsername)

		tokStr, err := base64.StdEncoding.DecodeString(token)
		if err != nil {
			panic(err)
		}
		tokParts := strings.Split(string(tokStr), ":")
		usernameMatch := ""
		passwdMatch := ""
		for _, u := range users {
			if u.Username == tokParts[0] {
				usernameMatch = u.Username
			}
			if u.Password == tokParts[1] {
				passwdMatch = u.Password
			}
		}
		if usernameMatch == "" {
			w.WriteHeader(500)
			w.Write([]byte("bad username"))
			return
		}
		if passwdMatch == "" {
			w.WriteHeader(500)
			w.Write([]byte("bad password"))
			return
		}

		var items []Item
		req, _ := http.NewRequest("GET", fmt.Sprintf("https://api.github.com/users/%s/repos", ghUsername), nil)
		req.SetBasicAuth(GITHUB_USERNAME, GITHUB_API_KEY)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
		}

		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		err = json.Unmarshal(respBody, &items)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}

		enc, err := json.Marshal(items)
		if err != nil {
			panic("uh oh")
		}

		_, err = w.Write(enc)
		if err != nil {
			log.Printf("Failed sending HTTP response body: %v", err)
		}
	})

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./public"))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	})

	(&http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}).ListenAndServe()
}

var users = []User{
	User{Username: "admin", Password: "password"},
}

type User struct {
	Username string
	Password string
}

type Item struct {
	Repo  string `json:"name"`
	Stars int    `json:"stargazers_count"`
}
