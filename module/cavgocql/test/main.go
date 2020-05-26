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
        "fmt"

       "goAgent/module/cavgocql"

      "github.com/gocql/gocql"
      logger "goAgent/logger"

)

func init() {

        var err error

        observer := cavgocql.NewObserver()

        cluster := gocql.NewCluster("127.0.0.1")
        cluster.QueryObserver = observer
	cluster.BatchObserver = observer
        cluster.Keyspace = "code2succeed"

        Session, err = cluster.CreateSession()

        if err != nil {

                panic(err)

        }

        fmt.Println("cassandra init done")

}

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
        nd.Method_exit(bt, "a.b.m1")
}
                                                                                                                   
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
        callGOCql(ctx)
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
        defer nd.Sdk_free()
        e.Start(":0000")

}
var Session *gocql.Session

type Emp struct {

	id        string

	firstName string

	lastName  string

	age       int

}


func createEmp(emp Emp,ctx context.Context) {


	if err := Session.Query("INSERT INTO emps(empid, first_name, last_name, age) VALUES(?, ?, ?, ?)",

		emp.id, emp.firstName, emp.lastName, emp.age).WithContext(ctx).Exec(); err != nil {

		logger.ErrorPrint("Error while inserting Emp")


	}

}



func updateEmp(emp Emp,ctx context.Context) {


	if err := Session.Query("UPDATE emps SET first_name = ?, last_name = ?, age = ? WHERE empid = ?",

		emp.firstName, emp.lastName, emp.age, emp.id).WithContext(ctx).Exec(); err != nil {

		 logger.ErrorPrint("Error while updating Emp")


	}

}



func deleteEmp(id string,ctx context.Context){


	if err := Session.Query("DELETE FROM emps WHERE empid = ?", id).WithContext(ctx).Exec(); err != nil {

	         logger.ErrorPrint("Error while deleting Emp")

		fmt.Println(err)

	}

}



func getEmps(ctx context.Context) []Emp {

	var emps []Emp

	m := map[string]interface{}{}



	iter := Session.Query("SELECT * FROM emps").WithContext(ctx).Iter()

	for iter.MapScan(m) {

		emps = append(emps, Emp{

			id:        m["empid"].(string),

			firstName: m["first_name"].(string),

			lastName:  m["last_name"].(string),

			age:       m["age"].(int),

		})

		m = map[string]interface{}{}

	}



	return emps

}



func callGOCql(ctx context.Context) {

	emp1 := Emp{"E-1", "Anupam", "Raj", 20}

	emp2 := Emp{"E-2", "Rahul", "Anand", 30}

	createEmp(emp1,ctx)

	createEmp(emp2,ctx)

	emp3 := Emp{"E-1", "Rahul", "Anand", 30}

	updateEmp(emp3,ctx)

	deleteEmp("E-2",ctx)

}
