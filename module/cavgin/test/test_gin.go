package main

import(
    "github.com/gin-gonic/gin"
    "time"
    "log"
    nd  "goAgent"
    md  "goAgent/module/cavgin"
    "strconv"
    "github.com/phayes/freeport"
    logger "goAgent/logger"
)

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

func method_check(bt uint64) {

     nd.Method_entry(bt, "method_check")

     nd.Method_exit(bt, "method_check")
}



func MainAdmin(c *gin.Context) {

     ctx:=c.Request.Context()

     bt := ctx.Value("CavissonTx").(uint64)

     method_check(bt)

     c.JSON(200, gin.H{

        "message": "pong",

      })

}

func cat_querry(c *gin.Context){

     name:=c.Query("name")

     age:=c.Query("age")

     c.JSON(200, gin.H{

      "name": name,
      "age": age,

     })
}


func main() {

     nd.Sdk_init()

     r := gin.New()

     r.Use(Logger())

     r.Use(DummyMiddleware)

     r.Use(md.Middleware(r))

     r.GET("/test", func(c *gin.Context) {

          example := c.MustGet("example").(string)

             logger.TracePrint(example)
})

     r.GET("/cat",cat_querry)

     r.GET("/",MainAdmin)

     port, err := freeport.GetFreePort()

        if err != nil {


        }

     logger.TracePrint(strconv.Itoa(port))

     r.Run(":" + strconv.Itoa(port))

     nd.Sdk_free()

}
