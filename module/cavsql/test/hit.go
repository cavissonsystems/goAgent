package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	nd "goAgent"
	md "goAgent/module/cavecho"
	"goAgent/module/cavsql"
	"log"
	"net/http"
)

type Pet struct {
	name    string
	owner   string
	species string
	sex     string
	birth   string
}

var (
	db *sql.DB
)

func pingAndQueryDB(db *sql.DB, bt uint64, ctx context.Context) {
	nd.Method_entry(bt, "pingAndQueryDB")
	{
		status := " Connection is up"
		if err := db.PingContext(ctx); err != nil {
			status = " connection is down"
		}
		log.Println(status)
	}
	rows, err := db.QueryContext(ctx, "SELECT * FROM pet")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	bks := make([]*Pet, 0)
	for rows.Next() {
		bk := new(Pet)
		err := rows.Scan(&bk.name, &bk.owner, &bk.species, &bk.sex, &bk.birth)
		if err != nil {
			log.Fatal(err)
		}
		bks = append(bks, bk)
	}
	if err = rows.Err(); err != nil {
		log.Fatal(err)
	}

	for _, bk := range bks {
		fmt.Printf("%s, %s, %s,%s,%s \n", bk.name, bk.owner, bk.species, bk.sex, bk.birth)
	}
	nd.Method_exit(bt, "pingAndQueryDB")

}

func callExample(ctx context.Context) {
	bt := ctx.Value("CavissonTx").(uint64)
	nd.Method_entry(bt, "callExample")
	db, err := cavsql.Open("mysql", "root:cavisson@/veterinary")

	if err != nil {
		log.Println("Error opening DB")
		log.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		pingAndQueryDB(db, bt, ctx)
	}

	nd.Method_exit(bt, "callExample")
	m1(bt)
	m2(bt)
	m3(bt)
	m4(bt)
	m5(bt)
	m6(bt)
	m7(bt)
	m8(bt)
	nd.BT_end(bt)

}

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

func mainAdmin(c echo.Context) error {
	req := c.Request()
	ctx := req.Context()
	bt := ctx.Value("CavissonTx").(uint64)
	callExample(ctx)
	m1(bt)
	m2(bt)

	return c.String(http.StatusOK, "ID is coming")

}

func check_hero(c echo.Context) error {
	req := c.Request()
	ctx := req.Context()
	bt := ctx.Value("CavissonTx").(uint64)
	callExample(ctx)
	m3(bt)
	m4(bt)
	return c.String(http.StatusOK, "hey there id conding")

}
func check_root(c echo.Context) error {
	req := c.Request()
	ctx := req.Context()
	bt := ctx.Value("CavissonTx").(uint64)
	callExample(ctx)
	m5(bt)
	m6(bt)
	return c.String(http.StatusOK, "hey there id conding")
}
func check_cavisson(c echo.Context) error {
	req := c.Request()
	ctx := req.Context()
	bt := ctx.Value("CavissonTx").(uint64)
	callExample(ctx)
	m7(bt)
	m8(bt)
	return c.String(http.StatusOK, "hey there id conding")

}

func main() {
	nd.Sdk_init()
	e := echo.New()
	e.Use(md.Middleware())
	g := e.Group("/admin")
	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{

		Format: `[${time_rfc3339}] ${status} ${method} ${host} ${path} ${latency_human}` + "\n",
	}))
	g.Use(middleware.BasicAuth(func(username string, password string, c echo.Context) (bool, error) {
		if username == "cavisson" && password == "cavisson" {
			return true, nil
		}
		return false, nil
	}))
	cavsql.Register("mysql", &mysql.MySQLDriver{})
	g.GET("/main", mainAdmin)
	g.GET("/hero", check_hero)
	g.GET("/root", check_root)
	g.GET("/cavisson", check_cavisson)
	defer nd.Sdk_free()
	e.Start(":0000")
}
