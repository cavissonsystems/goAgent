package main

import (
	"context"
	"github.com/go-redis/redis"
	nd "goAgent"
	logger "goAgent/logger"
	"goAgent/module/cavgoredis"
	"log"
	"time"
)

func m3(bt uint64) {
	nd.Method_entry(bt, "a.b.m3")
	time.Sleep(2 * time.Millisecond)
	nd.Method_exit(bt, "a.b.m3")
}

func m4(bt uint64) {
	nd.Method_entry(bt, "a.b.m4")
	time.Sleep(2 * time.Millisecond)
	nd.Method_exit(bt, "a.b.m4")
}

const (
	clientTypeBase = iota
	clientTypeCluster
	clientTypeRing
)

func Call_redis(ctx context.Context) {

	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	client1 := cavgoredis.Wrap(client).WithContext(ctx)
	pong, err := client1.Ping().Result()
	logger.TracePrint(pong)
	err = client1.Set("name", "Elliot", 0).Err()
	if err != nil {
		logger.ErrorPrint("Error : inserting value in redis")
	}

	val, err := client1.Get("name").Result()
	if err != nil {
		logger.ErrorPrint("Error : retrieving value from redis")
	}
	logger.TracePrint(val)
	var keys []string
	keys = append(keys, "foo")
	keys = append(keys, "bar")
	sc := client1.MGet(keys...)
	log.Println(sc)
	bt := ctx.Value("CavissonTx").(uint64)
	m3(bt)
	m4(bt)
}

func redisEmptyClient() *redis.Client {
	return redis.NewClient(&redis.Options{})
}

func redisEmptyClusterClient() *redis.ClusterClient {
	return redis.NewClusterClient(&redis.ClusterOptions{})
}

func redisEmptyRing() *redis.Ring {
	return redis.NewRing(&redis.RingOptions{})
}
