package main

import (
	"fmt"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
        "goAgent/module/cavnegroni"
        nd "goAgent"
        logger "goAgent/logger"

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

func m5(bt uint64){
        nd.Method_entry(bt, "a.b.m5")
        time.Sleep(2*time.Millisecond) 
        logger.TracePrint("m5 called")    
        nd.Method_exit(bt, "a.b.m5")

}

func m6(bt uint64){ 
        nd.Method_entry(bt, "a.b.m6")
        time.Sleep(2*time.Millisecond) 
        logger.TracePrint("m6 called")    
        nd.Method_exit(bt, "a.b.m6")

}

func Car(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>car</h1>")
        ctx := r.Context()
        bt := ctx.Value("CavissonTx").(uint64)
        m1(bt)
        m3(bt)

}

func Plane(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, "<h1>plane</h1>")
        ctx := r.Context()
        bt := ctx.Value("CavissonTx").(uint64)
        m2(bt)
        m4(bt)

}

func Middleware1(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
  //     t1 := time.Now()
       next(rw, r)
 //      t2 := time.Now()
       ctx := r.Context()
       bt := ctx.Value("CavissonTx").(uint64)
       m5(bt)
       m6(bt)


}

func main() {
        nd.Sdk_init()
	mux := mux.NewRouter()
        n := negroni.New()
//	n := negroni.New(negroni.HandlerFunc(Middleware1))
        n.Use(cavnegroni.Middleware())
        n.Use(negroni.HandlerFunc(Middleware1))
	mux.Path("/car").HandlerFunc(Car).Methods("GET").Name("Car")
	mux.Path("/plane").HandlerFunc(Plane).Methods("GET").Name("Plane")
	n.UseHandler(mux)
        defer nd.Sdk_free()
	n.Run(":1234")
        
}
