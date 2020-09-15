package main

import (
	"github.com/gin-gonic/gin"
	"github.com/phayes/freeport"
	nd "goAgent"
	logger "goAgent/logger"
	md "goAgent/module/cavgin"
	"log"
	"strconv"
	"time"
)

func method_check1(bt uint64) {

	nd.Method_entry(bt, "a.b.method_check1")

	nd.Method_exit(bt, "a.b.method_check1")
}

func method_check2(bt uint64) {

	nd.Method_entry(bt, "a.b.method_check2")

	nd.Method_exit(bt, "a.b.method_check2")
}

func method_check3(bt uint64) {

	nd.Method_entry(bt, "a.b.method_check3")

	nd.Method_exit(bt, "a.b.method_check3")
}

func method_check4(bt uint64) {

	nd.Method_entry(bt, "a.b.method_check4")

	nd.Method_exit(bt, "a.b.method_check4")
}

func method_check5(bt uint64) {

	nd.Method_entry(bt, "a.b.method_check5")

	nd.Method_exit(bt, "a.b.method_check5")
}

func method_check6(bt uint64) {

	nd.Method_entry(bt, "a.b.method_check6")

	nd.Method_exit(bt, "a.b.method_check6")
}

func Logger() gin.HandlerFunc {

	return func(c *gin.Context) {

		t := time.Now()

		c.Set("example", "12345")

		c.Next()

		latency := time.Since(t)

		log.Println(latency)

		status := c.Writer.Status()

		logger.TracePrint(strconv.Itoa(status))
	}
}

func DummyMiddleware(c *gin.Context) {

	c.Next()
}

func method_cavisson(c *gin.Context) {

	ctx := c.Request.Context()

	bt := ctx.Value("CavissonTx").(uint64)

	method_check1(bt)

	method_check2(bt)

	c.JSON(200, gin.H{

		"message": "pong",
	})

}

func method_root(c *gin.Context) {

	name := c.Query("name")

	age := c.Query("age")

	c.JSON(200, gin.H{

		"name": name,
		"age":  age,
	})

	ctx := c.Request.Context()

	bt := ctx.Value("CavissonTx").(uint64)

	method_check3(bt)

	method_check4(bt)

}

func method_test(c *gin.Context) {

	ctx := c.Request.Context()

	bt := ctx.Value("CavissonTx").(uint64)

	method_check5(bt)

	method_check6(bt)

	example := c.MustGet("example").(string)

	logger.TracePrint(example)
}

func main() {

	nd.Sdk_init()

	r := gin.New()

	r.Use(Logger())

	r.Use(DummyMiddleware)

	r.Use(md.Middleware(r))

	r.GET("/test", method_test)

	r.GET("/root", method_root)

	r.GET("/cavisson", method_cavisson)

	port, err := freeport.GetFreePort()

	if err != nil {

	}

	logger.TracePrint(strconv.Itoa(port))

	defer nd.Sdk_free()

	r.Run(":" + strconv.Itoa(port))

}
