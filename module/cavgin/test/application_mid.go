package main

import(
    "fmt"
    "github.com/gin-gonic/gin"
    "time"
    "log"
    nd  "goAgent"
    md  "goAgent/module/cavgin"
    "strconv"
    "github.com/phayes/freeport"
)

func Logger() gin.HandlerFunc {

    return func(c *gin.Context) {

       t := time.Now()

       c.Set("example", "12345")

       c.Next()

       latency := time.Since(t)

       log.Println(latency)

       status := c.Writer.Status()

       log.Println(status)
    }
}

func DummyMiddleware(c *gin.Context) {

     fmt.Println("Im a dummy!")

     c.Next()
}

func method_check(bt uint64) {

     log.Println("m1 called")

     nd.Method_entry(bt, "m1")

     nd.Method_exit(bt, "m1")
}



func MainAdmin(c *gin.Context) {

     log.Println("MainAdmin called")

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

     fmt.Println("hey")

     r := gin.New()

     r.Use(Logger())

     r.Use(DummyMiddleware)

     r.Use(md.Middleware(r))

     r.GET("/test", func(c *gin.Context) {

          example := c.MustGet("example").(string)

          log.Println(example)
     })

     r.GET("/cat",cat_querry)

     r.GET("/",MainAdmin)

     port, err := freeport.GetFreePort()

        if err != nil {

                log.Fatal(err)

        }

     fmt.Println(strconv.Itoa(port))

     r.Run(":" + strconv.Itoa(port))

     nd.Sdk_free()

}
