package cavecho

import (
	"net/http"
        nd "goAgent"
        "github.com/labstack/echo"
)

func Middleware()echo.MiddlewareFunc {
     return func(h echo.HandlerFunc)echo.HandlerFunc{
		 m := &middleware{
		 handler:        h,
		}
		return m.handle
	}
}

type middleware struct {
	handler        echo.HandlerFunc

}

//func context_redis()echo.Context {
  //     return echo.;

func (m *middleware) handle(c echo.Context) error {
        unique_id:="1"
        req := c.Request()
	name := c.Path()
        req = nd.Start_transacation(name,req)

        c.SetRequest(req)
        bt:= nd.Current_Transaction(req.Context())
        
        defer nd.BT_end(bt)
        nd.BT_store(bt,unique_id)
        resp := c.Response()
	var handlerErr error
	handlerErr = m.handler(c)
        if handlerErr != nil {
	 // will need later
	} else if !resp.Committed {
		resp.WriteHeader(http.StatusOK)
	}
	return handlerErr
}
