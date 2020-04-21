package main

import (
       "fmt"
	"net/http"
       "log"
       logger "goAgent/logger"
	"github.com/gorilla/mux"
       "goAgent/module/cavgorilla"
       nd "goAgent"
       "time"
)

func m1(bt uint64){
        nd.Method_entry(bt, "a.b.m1")
        logger.TracePrint("m1 called")
        time.Sleep(2*time.Millisecond)    
        nd.Method_exit(bt, "a.b.m1")

}

func m3(bt uint64){
        nd.Method_entry(bt, "a.b.m3")
        time.Sleep(2*time.Millisecond) 
        logger.TracePrint("m3 called")    
        nd.Method_exit(bt, "a.b.m3")

}

func m2(bt uint64){
        nd.Method_entry(bt, "a.b.m2")
        time.Sleep(2*time.Millisecond) 
        logger.TracePrint("m2 called")    
        nd.Method_exit(bt, "a.b.m2")

}

func m4(bt uint64){
        nd.Method_entry(bt, "a.b.m4")
        time.Sleep(2*time.Millisecond) 
        logger.TracePrint("m4 called")    
        nd.Method_exit(bt, "a.b.m4")

}

func main() {
        nd.Sdk_init()
        r := mux.NewRouter()
        cavgorilla.Instrument(r)
        r.HandleFunc("/hello",handler).Methods("GET")
        r.HandleFunc("/bye",handler2).Methods("GET")
        r.Use(loggingMiddleware)
        log.Println("server up at 4041 port")
        defer nd.Sdk_free()
        http.ListenAndServe(":4041", r)

}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
        ctx := r.Context()
        bt := ctx.Value("CavissonTx").(uint64)
        m1(bt)
        m3(bt)
}
func handler2(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "bye World!")
        ctx := r.Context()
        bt := ctx.Value("CavissonTx").(uint64)
        m2(bt)
        m4(bt)
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        logger.TracePrint(r.RequestURI)
        next.ServeHTTP(w, r)
    })
}
