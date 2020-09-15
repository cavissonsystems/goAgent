package cavecho

import (
	"github.com/labstack/echo"
	nd "goAgent"
	logger "goAgent/logger"
	"net/http"
)

func Middleware() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		m := &middleware{
			handler: h,
		}
		return m.handle
	}
}

type middleware struct {
	handler echo.HandlerFunc
}

func (m *middleware) handle(c echo.Context) error {
	unique_id := "1"

	req := c.Request()

	if req == nil {
		logger.ErrorPrint("Error : request not found")
	}

	name := c.Path()

	if name == "" {
		logger.ErrorPrint("Error : name not found")
	}

	req = nd.Start_transacation(name, req)

	c.SetRequest(req)

	bt := nd.Current_Transaction(req.Context())

	defer nd.BT_end(bt)

	nd.BT_store(bt, unique_id)

	resp := c.Response()

	var handlerErr error

	handlerErr = m.handler(c)

	if handlerErr != nil {
		// soon enough
	} else if !resp.Committed {
		resp.WriteHeader(http.StatusOK)
	}
	return handlerErr
}
