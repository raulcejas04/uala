package myRedis

import (
	"context"
	"fmt"

	redis "github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()

var Redisdb *redis.Client

func ConnectRedis() {

	if Redisdb == nil {
		Redisdb = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		})
	}

}

func SAdd(id string, set []string) {
	fmt.Printf("sadd %s %+v\n", id, set)
	sAdd := Redisdb.SAdd(ctx, id, set)
	if sAdd.Err() != nil {
		fmt.Println(sAdd.Err())
	}
}

func Delete(id string) {
	keys := GetAllWildcard(id)
	for _, key := range keys {
		_, err := Redisdb.Do(ctx, "DEL", key).Result()
		if err != nil {
			if err == redis.Nil {
				fmt.Println("key does not exists")
				return
			}
			panic(err)
		}

	}
}

func GetAllWildcard(id string) []string {

	var cursor uint64
	var keys []string
	var err error
	keys, cursor, err = Redisdb.Scan(ctx, cursor, id+"*", 100000).Result()
	if err != nil {
		log.Fatal("Error searching ", id)
		panic(err)
	}

	/*for _, key := range keys {
		fmt.Println("key", key)
	}*/
	return keys
}

func Smembers(key string) []string {

	s, err := Redisdb.SMembers(ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return s
}
