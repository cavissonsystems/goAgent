package main

import (
	"io"
	"os"
	"log"
        "net/http"
        "github.com/labstack/echo"
        "github.com/labstack/echo/middleware"
        md "goAgent/module/cavecho"
	nd "goAgent"
	ht "goAgent/module/cavhttp"
	"context"
        "time"
)


func save(c echo.Context) error {
	// Get name and email
	name := c.FormValue("name")
	email := c.FormValue("email")
	return c.String(http.StatusOK, "name:" + name + ", email:" + email)
}

func save1(c echo.Context) error {
        // Get name and email
        name := c.FormValue("name")
        email := c.FormValue("email")
        return c.String(http.StatusOK, "name:" + name + ", email:" + email)
}

func save2(c echo.Context) error {
        // Get name and email
        name := c.FormValue("name")
        email := c.FormValue("email")
        return c.String(http.StatusOK, "name:" + name + ", email:" + email)
}

func m1(bt uint64) {
        nd.Method_entry(bt, "a.b.m1")
       time.Sleep(2*time.Millisecond)
        nd.Method_exit(bt, "a.b.m1")
}

// cavhttp
func call_wrapclient(ctx context.Context){

	client := ht.WrapClient(http.DefaultClient)

        req, err := http.NewRequest("GET", "https://www.geeksforgeeks.org/find-triplets-array-whose-sum-equal-zero", nil)
	if err != nil {
		log.Println("Error : creating on new request")
	}
        req = req.WithContext(ctx)

	resp, err := client.Do(req)

	if err != nil {
		log.Println("Error : reading response. ")
	}
	defer resp.Body.Close()

	/* body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading body. ", err)
	}
	*/

	// writing the output to a file
	out, err := os.Create("ResponseBody.txt")
	if err != nil {
		log.Println("Error : creating responsebody txt file. ")
	}
	defer out.Close()
	io.Copy(out,resp.Body)
}

func mainAdmin(c echo.Context)error{
	req := c.Request()
        ctx := req.Context()

        call_wrapclient(ctx)
     //   Call_redis(ctx)
        bt := ctx.Value("CavissonTx").(uint64)
        Call_redigo(ctx)
        m1(bt)

	return c.String(http.StatusOK,"ID is coming")

}

func check1(c echo.Context)error{
      return c.String(http.StatusOK,"hey there id conding")

}

func ServerHeader(next echo.HandlerFunc)echo.HandlerFunc{
      return func(c echo.Context)error{
         c.Response().Header().Set(echo.HeaderServer,"BlueBot/1.0")
         return next(c)
	}
}


func main(){
	nd.Sdk_init()
	e:=echo.New()
	e.Use(ServerHeader)
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
	e.POST("/cats",save)
        e.POST("/dog",save1)
        e.POST("/rat",save2)
        g.GET("/hero",check1)
	e.Start(":0000")
	nd.Sdk_free()
}