package main

import (
	"net/http"
        "github.com/labstack/echo"
        "github.com/labstack/echo/middleware"
        md "goAgent/module/cavecho"
        "context"
	"fmt"
	nd "goAgent"
	"goAgent/module/cavgorm"
	"log"
	
      _ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type UserModel struct {
	ID      int    `gorm:"primary_key";"AUTO_INCREMENT"`
	Name    string `gorm:"size:255"`
	Address string `gorm:"type:varchar(100)"`
}

func callgorm(ctx context.Context) {
	 bt := ctx.Value("CavissonTx").(uint64)
         nd.Method_entry(bt, "callExample")
	db, err := cavgorm.Open("mysql", "root:cavisson@tcp(127.0.0.1:3306)/ormdemo?charset=utf8&parseTime=True")
	if err != nil {
		log.Panic(err)
	}
        db = cavgorm.WithContext(ctx, db)
	log.Println("Connection Established")
	db.DropTableIfExists(&UserModel{})
	db.AutoMigrate(&UserModel{})

	user := &UserModel{Name: "John", Address: "New York"}
	newUser := &UserModel{Name: "Martin", Address: "Los Angeles"}

	//****To insert or create the record in the database****.
	//NOTE: ##"we can insert multiple records here##".
	db.Debug().Create(user)
	fmt.Println("User Successfully Created")
	//**Also we can use save that will return primary key**
	db.Debug().Save(newUser)
	fmt.Println("New User Successfully Created")

	// "We can Update Record from here"
	db.Debug().Find(&user).Update("address", "California")
	//It will update John's address to California

	// Select, edit, and save
	db.Debug().Find(&user)
	user.Address = "Brisbane"
	db.Debug().Save(&user)

	//**Update with column names, not with attribute names**
	db.Debug().Model(&user).Update("Name", "Jack")

	db.Debug().Model(&user).Updates(
		map[string]interface{}{
			"Name":    "Amy",
			"Address": "Boston",
		})

	// UpdateColumn()
	db.Debug().Model(&user).UpdateColumn("Address", "Phoenix")
	db.Debug().Model(&user).UpdateColumns(
		map[string]interface{}{
			"Name":    "Taylor",
			"Address": "Houston",
		})
	// Using Find()
	db.Debug().Find(&user).Update("Address", "San Diego")

	// Batch Update
	db.Debug().Table("user_models").Where("address = ?", "california").Update("name", "Walker")

	// Select records and delete it
	db.Debug().Table("user_models").Where("address= ?", "San Diego").Delete(&UserModel{})

	db.Debug().Where("address = ?", "Los Angeles").First(&user)
	log.Println(user)
	db.Debug().Where("address = ?", "Los Angeles").Find(&user)
	log.Println(user)
	db.Debug().Where("address <> ?", "New York").Find(&user)
	log.Println(user)
	// IN
	db.Debug().Where("name in (?)", []string{"John", "Martin"}).Find(&user)
	log.Println(user)
	// LIKE
	db.Debug().Where("name LIKE ?", "%ti%").Find(&user)
	log.Println(user)
	// AND
	db.Debug().Where("name = ? AND address >= ?", "Martin", "Los Angeles").Find(&user)
	log.Println(user)

	//Find the record and delete it
	db.Where("address=?", "Los Angeles").Delete(&UserModel{})

	// Select all records from a model and delete all

	db.Debug().Model(&UserModel{}).Delete(&UserModel{})
        nd.Method_exit(bt, "callgorm")
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
func mainAdmin(c echo.Context)error{
        req := c.Request()
        ctx := req.Context()
        bt := ctx.Value("CavissonTx").(uint64)
        callgorm(ctx)
        m1(bt)
        m2(bt)
       
        return c.String(http.StatusOK,"ID is coming")

}

func check_hero(c echo.Context)error{
        req := c.Request()
        ctx := req.Context()
        bt := ctx.Value("CavissonTx").(uint64)
        callgorm(ctx)
        m3(bt)
        m4(bt)
        return c.String(http.StatusOK,"hey there id conding")

}
func check_root(c echo.Context)error{
        req := c.Request()
        ctx := req.Context()
        bt := ctx.Value("CavissonTx").(uint64)
        callgorm(ctx)
        m5(bt)
        m6(bt)
        return c.String(http.StatusOK,"hey there id conding")
}
func check_cavisson(c echo.Context)error{
        req := c.Request()
        ctx := req.Context()
        bt := ctx.Value("CavissonTx").(uint64)
        callgorm(ctx)
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
