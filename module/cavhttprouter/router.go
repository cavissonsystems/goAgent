
package cavhttprouter

import (
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Router struct {
	*httprouter.Router
}

func New() *Router {
	router := httprouter.New()
	router.NotFound = WrapNotFoundHandler(router.NotFound )
	router.MethodNotAllowed = WrapMethodNotAllowedHandler(router.MethodNotAllowed )
	return &Router{
		Router: router,

	}
}

func (r *Router) DELETE(path string, handle httprouter.Handle) {
	r.Router.DELETE(path, Wrap(handle, path))
}

func (r *Router) GET(path string, handle httprouter.Handle) {
	r.Router.GET(path, Wrap(handle, path))
}

func (r *Router) HEAD(path string, handle httprouter.Handle) {
	r.Router.HEAD(path, Wrap(handle, path))
}

func (r *Router) Handle(method, path string, handle httprouter.Handle) {
	r.Router.Handle(method, path, Wrap(handle, path))
}

func (r *Router) HandlerFunc(method, path string, handler http.HandlerFunc) {
	r.Handler(method, path, handler)
}

func (r *Router) Handler(method, path string, handler http.Handler) {
	r.Handle(method, path, func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		ctx := req.Context()
		ctx = context.WithValue(ctx, httprouter.ParamsKey, p)
		req = req.WithContext(ctx)
		handler.ServeHTTP(w, req)
	})
}

func (r *Router) OPTIONS(path string, handle httprouter.Handle) {
	r.Router.OPTIONS(path, Wrap(handle, path))
}

func (r *Router) PATCH(path string, handle httprouter.Handle) {
	r.Router.PATCH(path, Wrap(handle, path))
}

func (r *Router) POST(path string, handle httprouter.Handle) {
	r.Router.POST(path, Wrap(handle, path))
}

func (r *Router) PUT(path string, handle httprouter.Handle) {
	r.Router.PUT(path, Wrap(handle, path))
}

