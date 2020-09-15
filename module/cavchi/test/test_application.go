package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/phayes/freeport"
	nd "goAgent"
	md "goAgent/module/cavchi"
	"net/http"
	"strconv"
)

func m1(bt uint64) {
	nd.Method_entry(bt, "a.b.m1")
	nd.Method_exit(bt, "a.b.m1")
}

func m2(bt uint64) {
	nd.Method_entry(bt, "a.b.m2")
	nd.Method_exit(bt, "a.b.m2")
}

func m3(bt uint64) {
	nd.Method_entry(bt, "a.b.m3")
	nd.Method_exit(bt, "a.b.m3")
}
func m4(bt uint64) {
	nd.Method_entry(bt, "a.b.m4")
	nd.Method_exit(bt, "a.b.m4")
}

func m5(bt uint64) {
	nd.Method_entry(bt, "a.b.m5")
	nd.Method_exit(bt, "a.b.m5")
}

func m6(bt uint64) {
	nd.Method_entry(bt, "a.b.m6")
	nd.Method_exit(bt, "a.b.m6")
}

func m7(bt uint64) {
	nd.Method_entry(bt, "a.b.m7")
	nd.Method_exit(bt, "a.b.m7")
}
func m8(bt uint64) {
	nd.Method_entry(bt, "a.b.m8")
	nd.Method_exit(bt, "a.b.m8")
}

func op_check(w http.ResponseWriter, req *http.Request) {

	w.Write([]byte("root."))

	ctx := req.Context()

	bt := ctx.Value("CavissonTx").(uint64)

	m1(bt)

	m2(bt)

}

func op_root(w http.ResponseWriter, req *http.Request) {

	w.Write([]byte("Rascala bhAI"))

	ctx := req.Context()

	bt := ctx.Value("CavissonTx").(uint64)

	m3(bt)

	m4(bt)

}

func op_cavisson(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("root."))

	ctx := req.Context()

	bt := ctx.Value("CavissonTx").(uint64)

	m5(bt)

	m6(bt)

}

func op_hit(w http.ResponseWriter, req *http.Request) {

	w.Write([]byte("Rascala bhAI"))

	ctx := req.Context()

	bt := ctx.Value("CavissonTx").(uint64)

	m7(bt)

	m8(bt)

}

func main() {
	port, err := freeport.GetFreePort()

	if err != nil {

		fmt.Println(err)

	}

	fmt.Println(port)

	nd.Sdk_init()

	r := chi.NewRouter()

	r.Use(md.Middleware())

	r.Get("/check", op_check)

	r.Get("/root", op_root)

	r.Get("/cavisson", op_cavisson)

	r.Get("/hit", op_hit)

	defer nd.Sdk_free()

	http.ListenAndServe(":"+strconv.Itoa(port), r)
}
