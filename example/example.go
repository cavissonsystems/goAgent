package main

import (
	"fmt"
	nd "goAgent"
	logger "goAgent/logger"
	"time"
)

func m1(bt uint64) {
	nd.Method_entry(bt, "m1")
	logger.TracePrint("m1 called")
	nd.Method_exit(bt, "m1")
}

func m2(bt uint64) {
	nd.Method_entry(bt, "m2")
	logger.TracePrint("m2 called")
	nd.Method_exit(bt, "m2")
}

func m3(bt uint64) {
	nd.Method_entry(bt, "m3")
	logger.TracePrint("m3 called")
	nd.Method_exit(bt, "m3")
}

func m4(bt uint64) {
	nd.Method_entry(bt, "m4")
	logger.TracePrint("m4 called")
	nd.Method_exit(bt, "m4")
}

func sample() {
	bt := nd.BT_begin("postgres", "")
	nd.BT_store(bt, "1")
	nd.Method_entry(bt, "sample")
	m1(bt)
	fmt.Println("I'm inside the sample!")
	m2(bt)
	nd.Method_exit(bt, "sample")
	nd.BT_end(bt)
}

func sample1() {
	bt := nd.BT_begin("mongodb", "")
	nd.BT_store(bt, "2")
	nd.Method_entry(bt, "sample1")
	m3(bt)
	fmt.Println("I'm inside the sample1!")
	m4(bt)
	nd.Method_exit(bt, "sample1")
	nd.BT_end(bt)
}

func main() {
	nd.Sdk_init()
	for i := 0; i < 50; i++ {
		sample()
		sample1()
		time.Sleep(5000 * time.Millisecond)
	}
	nd.Sdk_free()
}
