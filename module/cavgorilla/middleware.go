
package cavgorilla

import (
	"net/http"

	"github.com/gorilla/mux"

	"goAgent/module/cavhttp"
)

func Instrument(r *mux.Router) {
	m := Middleware()
	r.Use(m)
	r.NotFoundHandler = WrapNotFoundHandler(r.NotFoundHandler, m)
	r.MethodNotAllowedHandler = WrapMethodNotAllowedHandler(r.MethodNotAllowedHandler, m)
}

func WrapNotFoundHandler(h http.Handler, m mux.MiddlewareFunc) http.Handler {
	if h == nil {
		h = http.NotFoundHandler()
	}
	return m(h)
}
func WrapMethodNotAllowedHandler(h http.Handler, m mux.MiddlewareFunc) http.Handler {
	if h == nil {
		h = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.WriteHeader(http.StatusMethodNotAllowed)
		})
	}
	return m(h)
}

func Middleware() mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return cavhttp.Wrap(h)
	}
}


