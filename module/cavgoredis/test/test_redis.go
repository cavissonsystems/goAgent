package main

import(
        "context"
        "github.com/go-redis/redis"
        "goAgent/module/cavgoredis"
        "fmt"
        "log"
)


const (
	clientTypeBase = iota
	clientTypeCluster
	clientTypeRing
)

var (
	unitTestCases = []struct {
		clientType int
		client     redis.UniversalClient
	}{
		{
			clientTypeBase,
			redisEmptyClient(),
		},
		{
			clientTypeBase,
			cavgoredis.Wrap(redisEmptyClient()),
		},
		{
			clientTypeBase,
			cavgoredis.Wrap(redisEmptyClient()).WithContext(context.Background()),
		},
		{
			clientTypeCluster,
			redisEmptyClusterClient(),
		},
		{
			clientTypeCluster,
			cavgoredis.Wrap(redisEmptyClusterClient()),
		},
		{
			clientTypeCluster,
			cavgoredis.Wrap(redisEmptyClusterClient()).WithContext(context.Background()),
		},
		{
			clientTypeRing,
			redisEmptyRing(),
		},
		{
			clientTypeRing,
			cavgoredis.Wrap(redisEmptyRing()),
		},
		{
			clientTypeRing,
			cavgoredis.Wrap(redisEmptyRing()).WithContext(context.Background()),
		},
	}
)

func Call_redis(ctx context.Context) {

	fmt.Println("Go Redis Tutorial")

	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		Password: "",
		DB: 0,
	})

        client1 := cavgoredis.Wrap(client).WithContext(ctx)
	pong, err := client1.Ping().Result()
	fmt.Println(pong, err)

	err = client1.Set("name", "Elliot", 0).Err()
	if err != nil {
		log.Println("Error : inserting value in redis" )
	}

	val, err := client1.Get("name").Result()
	if err != nil {
		log.Println("Error : retrieving value from redis")
	}
	fmt.Println(val)

	var keys []string
	keys = append(keys, "foo")
	keys = append(keys, "bar")
	sc := client1.MGet(keys...)
        fmt.Println(sc)
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
