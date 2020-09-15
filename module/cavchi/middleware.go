package cavchi

import (
	md "goAgent/module/cavhttp"
	"net/http"
)

func Middleware() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return md.Wrap(
			h,
		)
	}
}
