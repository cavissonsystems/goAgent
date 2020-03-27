package main

import (
	"fmt"
	"net/http"
       "log"
	"github.com/gorilla/mux"
       "goAgent/module/cavgorilla"
  nd    "goAgent"
)

func main() {
        nd.Sdk_init()
	r := mux.NewRouter()
        cavgorilla.Instrument(r)
        r.HandleFunc("/hello",handler).Methods("GET")
        r.HandleFunc("/bye",handler2).Methods("GET")
        r.Use(loggingMiddleware)
        log.Println("server up")
	http.ListenAndServe(":4041", r)
        nd.Sdk_free()

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
func handler2(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "bye World!")
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.RequestURI)
        next.ServeHTTP(w, r)
    })
}
