package  main

import (
	"log"
        "goAgent/module/cavredigo"
	"github.com/gomodule/redigo/redis"
       "context"
      nd "goAgent"
        "time"
)

func m3(bt uint64) {
        nd.Method_entry(bt, "a.b.m3")
        time.Sleep(2*time.Millisecond)
        nd.Method_exit(bt, "a.b.m3")
}

func m2(bt uint64) {
        nd.Method_entry(bt, "a.b.m2")
        time.Sleep(2*time.Millisecond)
        nd.Method_exit(bt, "a.b.m2")
}


func Call_redigo(ctx context.Context) {

	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}

        conn = cavredigo.Wrap(conn).WithContext(ctx)
	defer conn.Close()

	_, err = conn.Do("HMSET", "album:1", "title", "Electric Ladyland", "artist", "Jimi Hendrix", "price", 4.95, "likes", 8)
	if err != nil {
		log.Fatal(err)
	}

         _, err = conn.Do("HMSET", "album:2", "title", "software", "artist", "maxwell", "price", 7.5, "likes", 10)
        if err != nil {
                log.Fatal(err)
        }

	title, err := redis.String(conn.Do("HGET", "album:1", "title"))
	if err != nil {
		log.Fatal(err)
	}

	artist, err := redis.String(conn.Do("HGET", "album:1", "artist"))
	if err != nil {
		log.Fatal(err)
	}

	price, err := redis.Float64(conn.Do("HGET", "album:1", "price"))
	if err != nil {
		log.Fatal(err)
	}

	likes, err := redis.Int(conn.Do("HGET", "album:1", "likes"))
	if err != nil {
		log.Fatal(err)
	}
       title2, err := redis.String(conn.Do("HGET", "album:2", "title"))
        if err != nil {
                log.Fatal(err)
        }

        artist2, err := redis.String(conn.Do("HGET", "album:2", "artist"))
        if err != nil {
                log.Fatal(err)
        }

        price2, err := redis.Float64(conn.Do("HGET", "album:2", "price"))
        if err != nil {
                log.Fatal(err)
        }

        likes2, err := redis.Int(conn.Do("HGET", "album:2", "likes"))
        if err != nil {
                log.Fatal(err)
        }

        bt  :=  ctx.Value("CavissonTx").(uint64)
        m2(bt)
        m3(bt)
	log.Printf("%s by %s: £%.2f [%d likes]\n", title, artist, price, likes)
        log.Printf("%s by %s: £%.2f [%d likes]\n", title2, artist2, price2, likes2)
}
