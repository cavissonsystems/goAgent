package main

import (
	"context"
	"fmt"
	nd "goAgent"
	"goAgent/module/cavgorm"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type UserModel struct {
	ID      int    `gorm:"primary_key";"AUTO_INCREMENT"`
	Name    string `gorm:"size:255"`
	Address string `gorm:"type:varchar(100)"`
}

func callExample() {
	bt := nd.BT_begin("admin/main", "")
	fmt.Println(bt)

	// method entry
	nd.BT_store(bt, "122")
	nd.Method_entry(bt, "callExample")
	// call bt be0gin here of sdk
	ctx := context.Background()
	// create new context -- using background method
	ctx = context.WithValue(ctx, "CavissonTx", bt)
	db, err := cavgorm.Open("mysql", "root:cavisson@tcp(127.0.0.1:3306)/ormdemo?charset=utf8&parseTime=True")
	if err != nil {
		log.Panic(err)
	}

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
	time.Sleep(5 * time.Second)
	nd.Method_exit(bt, "callExample")
	time.Sleep(5 * time.Second)
	fmt.Println("Value of bt before bt end is ", bt)
	bt_end := nd.BT_end(bt)
	fmt.Println("Value of bt end ", bt_end)
}

func main() {
	nd.Sdk_init()
	callExample()
	nd.Sdk_free()
}
