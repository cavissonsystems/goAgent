
package cavhttprouter

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
        nd  "goAgent"
	"goAgent/module/cavhttp"
  logger "goAgent/logger"

)


func Wrap(h httprouter.Handle, route string) httprouter.Handle {

	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
              unique_id :="1"
                name := req.URL.Path
                req = nd.Start_transacation(name,req)
                bt := nd.Current_Transaction(req.Context())
                defer nd.BT_end(bt)
                nd.BT_store(bt, unique_id)
       
                h(w, req, p)
	}
}

func WrapNotFoundHandler(h http.Handler) http.Handler {
	if h == nil {
                logger.ErrorPrint("Error : Handler not found")

		h = http.NotFoundHandler()
	}
	return wrapHandlerUnknownRoute(h)
}

func WrapMethodNotAllowedHandler(h http.Handler) http.Handler {
	if h == nil {
                logger.ErrorPrint("Error : method not found")
		h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		})
	}
	return wrapHandlerUnknownRoute(h)
}

func wrapHandlerUnknownRoute(h http.Handler) http.Handler {
	return cavhttp.Wrap(h)
}


