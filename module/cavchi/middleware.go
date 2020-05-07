package cavchi

import (
	"net/http"
        md  "goAgent/module/cavhttp"
)

func Middleware() func(http.Handler) http.Handler {
       return func(h http.Handler) http.Handler {
		return md.Wrap(
			h,
		)
	}
}


