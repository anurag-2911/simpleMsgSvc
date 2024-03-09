package pkg

import (
	"context"
	"log"
	"simpleMsgSvc/config"
	"time"
	"github.com/go-redis/redis/v8"
)

var redisMasterClient *redis.Client
var redisReplicaClient *redis.Client

func init() {
	defer func(){
		if r:=recover();r!=nil{
			log.Println("recovered from panic in redis client init ")
		}
	}()
	redisMasterClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.REDIS_MASTER_ADDR,
		Password: config.RedisConfig.REDIS_PASSWORD,
		DB:       0, // default DB
	})

	// Initialize Redis replica client for read-only operations
	redisReplicaClient = redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.REDIS_REPLICA_ADDR,
		Password: config.RedisConfig.REDIS_PASSWORD,
		DB:       0, // default DB
	})
}

// getValue retrieves a value by key from Redis
func GetValue(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	value, err := redisReplicaClient.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Println("error in getting the key ", err)
		return "", err
	} else if err != nil {

		return "", err
	}
	return value, nil
}

// setValue sets a key-value pair in Redis
func SetValue(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	err := redisMasterClient.Set(ctx, key, value, 0).Err()
	if err != nil {
		log.Printf("error in setting the value %s for the key %s ,error %v\n ", key, value, err)
		return err
	}
	return nil
}
