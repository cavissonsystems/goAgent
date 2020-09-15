package main

import (
	"encoding/json"
	"fmt"
	"github.com/phayes/freeport"
	wr "go-agent/module/cavhttp"
	"log"
	"net/http"
	"strconv"
)

type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

func returnAllArticles(w http.ResponseWriter, r *http.Request) {
	articles := []Article{
		Article{Title: "Hello", Desc: "Article Description", Content: "Article Content"},
		Article{Title: "Hello 2", Desc: "Article Description", Content: "Article Content"},
	}

	fmt.Println("Endpoint Hit: returnAllArticles")
	json.NewEncoder(w).Encode(articles)
}

func main() {
	port, err := freeport.GetFreePort()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strconv.Itoa(port))
	h := http.HandlerFunc(homePage)
	h1 := wr.Wrap(h)
	fmt.Println(h1)
	http.Handle("/", h1)
	// add our articles route and map it to our
	//http.HandleFunc("/articles", returnAllArticles)
	// h :=wr.Wrap(http.HandlerFunc(homePage))
	//wr.Wrap(h)
	//h.ServeHTTP()
	http.HandleFunc("/articles", returnAllArticles)
	log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
}
