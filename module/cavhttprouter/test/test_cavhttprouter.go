package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	nd "goAgent"
	logger "goAgent/logger"
	"goAgent/module/cavhttprouter"
	"log"
	"net/http"
	"time"
)

func m1(bt uint64) {
	nd.Method_entry(bt, "a.b.m1")
	logger.TracePrint("m1 called")
	time.Sleep(2 * time.Millisecond)
	nd.Method_exit(bt, "a.b.m1")
}
func m3(bt uint64) {
	nd.Method_entry(bt, "a.b.m3")
	time.Sleep(2 * time.Millisecond)
	logger.TracePrint("m3 called")
	nd.Method_exit(bt, "a.b.m3")
}

func m2(bt uint64) {
	nd.Method_entry(bt, "a.b.m2")
	time.Sleep(2 * time.Millisecond)
	logger.TracePrint("m2 called")
	nd.Method_exit(bt, "a.b.m2")
}
func m4(bt uint64) {
	nd.Method_entry(bt, "a.b.m4")
	time.Sleep(2 * time.Millisecond)
	logger.TracePrint("m4 called")
	nd.Method_exit(bt, "a.b.m4")
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
	ctx := r.Context()
	bt := ctx.Value("CavissonTx").(uint64)
	m1(bt)
	m3(bt)

}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
	ctx := r.Context()
	bt := ctx.Value("CavissonTx").(uint64)
	m2(bt)
	m4(bt)

}

func main() {
	nd.Sdk_init()
	router := cavhttprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)
	log.Println("server up at 3000 port")
	defer nd.Sdk_free()
	log.Fatal(http.ListenAndServe(":3000", router))

}
