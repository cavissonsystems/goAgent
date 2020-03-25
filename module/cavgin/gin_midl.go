package cavgin

import (
	"github.com/gin-gonic/gin"
        "log"
        nd "goAgent"
)

func Middleware(engine *gin.Engine) gin.HandlerFunc {

     log.Println("hey")

      m := &middleware{

		engine:         engine,

	}

	return m.handle
}

type middleware struct {

    engine         *gin.Engine

}



func (m *middleware) handle(c *gin.Context) {

        unique_id:="1"

	handlerName := c.HandlerName()

	log.Println(handlerName)

	req := nd.Start_transacation(handlerName,c.Request)

	log.Println("value of request\n")

        c.Request = req

        log.Println("value of update request")

	bt := nd.Current_Transaction(req.Context())

        log.Println("got bt value")

        log.Println(bt)

        defer nd.BT_end(bt)

        log.Println("bt_end")

        log.Println(bt)

        log.Println(unique_id)

	nd.BT_store(bt, unique_id)

        log.Println("bt_store")

	c.Next()
}


