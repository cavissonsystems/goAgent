package cavhttp

import (
	"net/http"
        nd "goAgent"

)

func WrapClient(c *http.Client) *http.Client {
	if c == nil {
		c = http.DefaultClient
	}
	copied :=  *c
	copied.Transport = WrapRoundTripper(copied.Transport)
        return &copied
}

func WrapRoundTripper(r http.RoundTripper) http.RoundTripper {
	if r == nil {
		r = http.DefaultTransport
	}
	rt := &roundTripper{
		r:              r,
	}

	return rt
}

type roundTripper struct {
	r              http.RoundTripper
}

func (r *roundTripper) RoundTrip(req *http.Request) (*http.Response, error){
	ctx := req.Context()
	bt := ctx.Value("CavissonTx").(uint64)
	ip_handle_int := nd.IP_http_callout_begin(bt , req.Host, req.URL.Path)
        resp, err := r.r.RoundTrip(req)
       defer nd.IP_http_callout_end(bt ,  ip_handle_int)
	return resp, err
}
