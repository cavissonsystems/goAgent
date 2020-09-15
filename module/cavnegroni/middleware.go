package cavnegroni

import (
	"context"
	"net/http"

	"github.com/urfave/negroni"
	"goAgent/module/cavhttp"
)

func Middleware() negroni.Handler {
	m := &middleware{
		handler: cavhttp.Wrap(http.HandlerFunc(nextHandler)),
	}
	return m
}

type middleware struct {
	handler http.Handler
}

func (m *middleware) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	r = r.WithContext(context.WithValue(r.Context(), nextKey{}, next))
	m.handler.ServeHTTP(w, r)
}

type nextKey struct{}

func nextHandler(w http.ResponseWriter, r *http.Request) {
	next := r.Context().Value(nextKey{}).(http.HandlerFunc)
	next(w, r)
}
