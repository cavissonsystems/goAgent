package nd

import (
	"github.com/labstack/echo"
        "github.com/labstack/echo/middleware"
        md "goAgent/module/cavecho"
)

func Echointegration(e *echo.Echo){
	middle := md.Middleware()
	e.Use(middle)
}
