package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Article struct {
	Name   string `json:"name"`
	Author string `json:"author"`
}

var Articles = []Article{
	{Name: "Black Swan", Author: "Nassim Nicholas Taleb"},
	{Name: "47 Laws Of Power", Author: "Robert Greene"},
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	if name := query.Get("name"); name != "" {
		var results []Article
		for _, article := range Articles {
			if strings.Contains(strings.ToLower(article.Name), strings.ToLower(name)) {
				results = append(results, article)
			}
		}
		json.NewEncoder(w).Encode(results)
		return
	}
	json.NewEncoder(w).Encode(Articles[:min(len(Articles), 50)])
}

func AddArticle(w http.ResponseWriter, r *http.Request) {
	var newArticle Article
	err := json.NewDecoder(r.Body).Decode(&newArticle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	Articles = append(Articles, newArticle)
	json.NewEncoder(w).Encode(Articles)
}

func UpdateArticle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	var updatedArticle Article
	err := json.NewDecoder(r.Body).Decode(&updatedArticle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	for i, article := range Articles {
		if strings.EqualFold(article.Name, name) {
			Articles[i] = updatedArticle
			break
		}
	}
	json.NewEncoder(w).Encode(Articles)
}

func DeleteArticle(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["name"]
	for i, article := range Articles {
		if strings.EqualFold(article.Name, name) {
			Articles = append(Articles[:i], Articles[i+1:]...)
			break
		}
	}
	json.NewEncoder(w).Encode(Articles)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	log.Println("[+] Started Server!")
	router := mux.NewRouter()
	router.HandleFunc("/articles", GetArticle).Methods("GET")
	router.HandleFunc("/articles", AddArticle).Methods("POST")
	router.HandleFunc("/articles/{name}", UpdateArticle).Methods("PUT")
	router.HandleFunc("/articles/{name}", DeleteArticle).Methods("DELETE")
	http.ListenAndServe(":8000", router)
}
