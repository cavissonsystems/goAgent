package cavhttp

import (
	nd "goAgent"
	logger "goAgent/logger"
	"net/http"
)

func Wrap(h http.Handler) http.Handler {
	if h == nil {
		panic("h == nil")
		logger.ErrorPrint("Error : handler not found")

	}

	handler := &handler{

		handler: h,
	}
	return handler
}

type handler struct {
	handler http.Handler
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	unique_id := "1"

	name := req.URL.Path

	req = nd.Start_transacation(name, req)

	bt := nd.Current_Transaction(req.Context())

	defer nd.BT_end(bt)

	nd.BT_store(bt, unique_id)

	h.handler.ServeHTTP(w, req)
}
