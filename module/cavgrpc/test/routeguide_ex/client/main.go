package main

import (
	"fmt"
     pb "google.golang.org/grpc/examples/routeguide_ex/routeguide_ex"
     nd "goAgent"
	"log"
	"net/http"
	"strconv"
        "goAgent/module/cavgrpc"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
       logger "goAgent/logger"
)

/*func m1(bt uint64) {
        nd.Method_entry(bt, "m1")
        logger.TracePrint("m1 called")    
        fmt.Println("m1 called") 
        nd.Method_exit(bt, "m1")
}

func m2(bt uint64) {
        nd.Method_entry(bt, "m2")
        logger.TracePrint("m2 called")    
        fmt.Println("m2 called") 
        nd.Method_exit(bt, "m2")
}
*/

                     
func main() {
        nd.Sdk_init()
         conn, err := grpc.Dial("localhost:4040", grpc.WithInsecure(),grpc.WithUnaryInterceptor(cavgrpc.NewUnaryClientInterceptor()))
        if err != nil {
                panic(err)
        }

        client := pb.NewAddServiceClient(conn)

        g := gin.Default()
        g.GET("/add/:a/:b", func(ctx *gin.Context) {
                a, err := strconv.ParseUint(ctx.Param("a"), 10, 64)
                if err != nil {
                        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter A"})
                        return
                }

                b, err := strconv.ParseUint(ctx.Param("b"), 10, 64)
                if err != nil {
                        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter B"})
                        return
                }

                req := &pb.Request{A: int64(a), B: int64(b)}
                if response, err := client.Add(ctx, req); err == nil {
                        ctx.JSON(http.StatusOK, gin.H{
                                "result": fmt.Sprint(response.Result),
                        })
                } else {
                        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                
                     


}
                    
                

        })
        g.GET("/mult/:a/:b", func(ctx *gin.Context) {
                a, err := strconv.ParseUint(ctx.Param("a"), 10, 64)
                if err != nil {
                        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter A"})
                        return
                }
                b, err := strconv.ParseUint(ctx.Param("b"), 10, 64)
                if err != nil {
                        ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Parameter B"})
                        return
                }
                req := &pb.Request{A: int64(a), B: int64(b)}

                if response, err := client.Multiply(ctx, req); err == nil {
                        ctx.JSON(http.StatusOK, gin.H{
                                "result": fmt.Sprint(response.Result),
                        })
                } else {
                        ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                       

                }
         
      
        })
        log.Println("server run at 4041")
        if err := g.Run(":4041"); err != nil {
                log.Fatalf("Failed to run server: %v", err)
        }


        nd.Sdk_free() 
}
