package main

import(
        "context"
        "github.com/go-redis/redis"
        "goAgent/module/cavgoredis"
      logger  "goAgent/logger"
       "log"
)


const (
	clientTypeBase = iota
	clientTypeCluster
	clientTypeRing
)


func Call_redis(ctx context.Context) {


	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

        client1 := cavgoredis.Wrap(client).WithContext(ctx)
	pong, err := client1.Ping().Result()
        logger.TracePrint(pong)
	err = client1.Set("name", "Elliot", 0).Err()
	if err != nil {
		logger.ErrorPrint("Error : inserting value in redis" )
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
