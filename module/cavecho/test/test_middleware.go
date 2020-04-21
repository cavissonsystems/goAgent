package main

import (
        "net/http"
        "github.com/labstack/echo"
        "github.com/labstack/echo/middleware"
        md "goAgent/module/cavecho"
	nd "goAgent"
	"fmt"
)

func m1(bt uint64) {
        nd.Method_entry(bt, "a.b.m1")
        nd.Method_exit(bt, "a.b.m1")
}

func m2(bt uint64) {
        nd.Method_entry(bt, "a.b.m2")
        nd.Method_exit(bt, "a.b.m2")
}

func m3(bt uint64) {
        nd.Method_entry(bt, "a.b.m3")
        nd.Method_exit(bt, "a.b.m3")
}
func m4(bt uint64) {
        nd.Method_entry(bt, "a.b.m4")
        nd.Method_exit(bt, "a.b.m4")
}




func m5(bt uint64) {
        nd.Method_entry(bt, "a.b.m5")
        nd.Method_exit(bt, "a.b.m5")
}

func m6(bt uint64) {
        nd.Method_entry(bt, "a.b.m6")
        nd.Method_exit(bt, "a.b.m6")
}

func m7(bt uint64) {
        nd.Method_entry(bt, "a.b.m7")
        nd.Method_exit(bt, "a.b.m7")
}
func m8(bt uint64) {
        nd.Method_entry(bt, "a.b.m8")
        nd.Method_exit(bt, "a.b.m8")
}




func mainAdmin(c echo.Context)error{
	req := c.Request()
        ctx := req.Context()
        bt := ctx.Value("CavissonTx").(uint64)
        m1(bt)
        m2(bt)
	return c.String(http.StatusOK,"ID is coming")

}

func check_hero(c echo.Context)error{
        req := c.Request()
        ctx := req.Context()
        bt := ctx.Value("CavissonTx").(uint64)
        m3(bt)
        m4(bt)
        return c.String(http.StatusOK,"hey there id conding")

}


func check_root(c echo.Context)error{
        req := c.Request()
        ctx := req.Context()
        bt := ctx.Value("CavissonTx").(uint64)
        m5(bt)
        m6(bt)
        return c.String(http.StatusOK,"hey there id conding")
}
func check_cavisson(c echo.Context)error{
        req := c.Request()
        ctx := req.Context()
        bt := ctx.Value("CavissonTx").(uint64)
        m7(bt)
        m8(bt)
        return c.String(http.StatusOK,"hey there id conding")


}

func main(){
	nd.Sdk_init()
	e:=echo.New()
	e.Use(md.Middleware())
	g:=e.Group("/admin")
	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{

	Format:`[${time_rfc3339}] ${status} ${method} ${host} ${path} ${latency_human}` +"\n",

	}))
	g.Use(middleware.BasicAuth(func(username string,password string,c echo.Context)(bool,error){
	if username =="cavisson" && password =="cavisson"{
		return true,nil
	}
	return false,nil
	}))
	g.GET("/main",mainAdmin)
        g.GET("/hero",check_hero)
        g.GET("/root",check_root)
        g.GET("/cavisson",check_cavisson)
        defer nd.Sdk_free()
	e.Start(":0000")

}
