package cavgin

import (
	"github.com/gin-gonic/gin"
        nd "goAgent"
        logger "goAgent/logger"
        "strings"
)

func Middleware(engine *gin.Engine) gin.HandlerFunc {

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

        var Name []string = strings.Split(handlerName, ".")

        handlerName=Name[1]

	logger.TracePrint(handlerName)

	req := nd.Start_transacation(handlerName,c.Request)

        c.Request = req

        bt := nd.Current_Transaction(req.Context())

        defer nd.BT_end(bt)

	nd.BT_store(bt, unique_id)

	c.Next()
}


