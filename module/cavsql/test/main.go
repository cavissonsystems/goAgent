package main

import (
	"context"
	"database/sql"
	"fmt"
	nd "goAgent"
	"goAgent/module/cavsql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Pet struct {
	name    string
	owner   string
	species string
	sex     string
	birth   string
}

var (
	ctx context.Context
	db  *sql.DB
)

func pingAndQueryDB(db *sql.DB, bt uint64,ctx context.Context) {
	nd.Method_entry(bt, "pingAndQueryDB")
	fmt.Println("Value of bt from main.go is", bt)
	{
		ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		//ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		status := " Connection is up"
		if err := db.PingContext(ctx); err != nil {
			status = " connection is down"
		}
		log.Println(status)
	}
	rows, err := db.Query("SELECT * FROM pet")
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
	//time.Sleep(5 * time.Second)
}

func callExample() {
	bt := nd.BT_begin("admin/main", "")

	// method entry
	nd.BT_store(bt, "122")
	nd.Method_entry(bt, "callExample")
	// call bt be0gin here of sdk
	ctx := context.Background()
	// create new context -- using background method
	ctx = context.WithValue(ctx, "CavissonTx", bt)
	// add bt id into context
	cavsql.Register("mysql", &mysql.MySQLDriver{})
	// call register
	db, err := cavsql.Open("mysql", "root:cavisson@/veterinary") // correct

	if err != nil {
		log.Println("Error opening DB")
		log.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		pingAndQueryDB(db, bt,ctx)
	}

	time.Sleep(5 * time.Second)

	nd.Method_exit(bt, "callExample")
	time.Sleep(5 * time.Second)

	//fmt.Println("Value of bt before bt end is ", bt)
	nd.BT_end(bt)
	//fmt.Println("Value of bt end ", bt_end)
}

func main() {
	nd.Sdk_init()
	callExample()
	nd.Sdk_free()
}
