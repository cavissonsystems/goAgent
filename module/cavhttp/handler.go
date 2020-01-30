package cavhttp

import (
//	"context"
	"net/http"
        "log"
        nd "go-agent"
        // "strings"
       //cl "go-agent/example"
)

func m1(bt uint64) {
        log.Printf("invoke m1 method")
        nd.Method_entry(bt, "m1")
        log.Println("m1 called")
        nd.Method_exit(bt, "m1")
}
func Wrap(h http.Handler) http.Handler {
	if h == nil {
		panic("h == nil")
	}
	handler := &handler{
		handler:        h,
	//	requestName:    ServerRequestName,
	}
	/*for _, o := range o {
		o(handler)
	}*/
	log.Printf("Inside wrap func")
	return handler
}

type handler struct {
	handler          http.Handler
	//requestName      RequestNameFunc
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	nd.Sdk_init()
	log.Printf("invoke handle function")
        unique_id:="1"
        name := req.URL.Path            //"/admin/name"       strings.TrimLeft(req.URL.Path, "/")
        log.Printf("invoke middleware:%v\n",name)
        req = nd.Start_transacation(name,req)
        bt:= nd.Current_Transaction(req.Context())
	defer nd.BT_end(bt)
        nd.BT_store(bt,unique_id)
	m1(bt)
	h.handler.ServeHTTP(w, req)
	nd.Sdk_free()
}
