package main

import (
       "fmt"
	"net/http"
       "log"
      logger "goAgent/logger"
	"github.com/gorilla/mux"
       "goAgent/module/cavgorilla"
  nd    "goAgent"
)

func m1(bt uint64){
        nd.Method_entry(bt, "m1")
        logger.TracePrint("m1 called")    
        nd.Method_exit(bt, "m1")

}

func m2(bt uint64){
        nd.Method_entry(bt, "m2")
        logger.TracePrint("m2 called")    
        nd.Method_exit(bt, "m2")

}

func main() {
        nd.Sdk_init()
        r := mux.NewRouter()
        cavgorilla.Instrument(r)
        r.HandleFunc("/hello",handler).Methods("GET")
        r.HandleFunc("/bye",handler2).Methods("GET")
        r.Use(loggingMiddleware)
        log.Println("server up at 4041 port")
        http.ListenAndServe(":4041", r)
        nd.Sdk_free()

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
        ctx := r.Context()
        bt := ctx.Value("CavissonTx").(uint64)
        m1(bt)
}
func handler2(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "bye World!")
        ctx := r.Context()
        bt := ctx.Value("CavissonTx").(uint64)
        m2(bt)
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        logger.TracePrint(r.RequestURI)
        next.ServeHTTP(w, r)
    })
}
