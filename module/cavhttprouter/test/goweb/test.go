package main

import (
	"fmt"
	"github.com/justinas/alice"
	hr "goAgent/module/cavhttprouter/test/goweb/httprouterwrapper"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	nd "goAgent"
	logger "goAgent/logger"
	"goAgent/module/cavhttprouter"
)

func m1(bt uint64) {
	nd.Method_entry(bt, "m1")
	logger.TracePrint("m1 called")
	nd.Method_exit(bt, "m1")
}

func m2(bt uint64) {
	nd.Method_entry(bt, "m2")
	logger.TracePrint("m2 called")
	nd.Method_exit(bt, "m2")
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
	ctx := r.Context()
	bt := ctx.Value("CavissonTx").(uint64)
	m1(bt)

}

func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
	ctx := r.Context()
	bt := ctx.Value("CavissonTx").(uint64)
	m2(bt)

}

func IndexGET(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "bye")
}
func main() {
	nd.Sdk_init()
	router := cavhttprouter.New()
	router.GET("/", Index)

	router.GET("/bye", hr.Handler(alice.
		New().
		ThenFunc(IndexGET)))

	router.GET("/hello/:name", Hello)
	log.Println("server up at 3000 port")
	log.Fatal(http.ListenAndServe(":3000", router))
	nd.Sdk_free()
}
