package config

import (
	"log"
	"os"
)

type redisConfiguration struct {
	REDIS_MASTER_ADDR  string
	REDIS_REPLICA_ADDR string
	REDIS_PASSWORD     string
}

var RedisConfig *redisConfiguration

func init() {
	defer func() {
		if r := recover(); r != nil {
			log.Println("recovered from panic in reading redis configuration")
		}
	}()
	readRedisConfiguration()
}

func readRedisConfiguration() {
	RedisConfig = &redisConfiguration{
		REDIS_MASTER_ADDR:  os.Getenv("REDIS_MASTER_ADDR"),
		REDIS_REPLICA_ADDR: os.Getenv("REDIS_REPLICA_ADDR"),
		REDIS_PASSWORD:     os.Getenv("REDIS_PASSWORD"),
	}
	log.Println("read redis configuration")
	log.Println("master addr ", RedisConfig.REDIS_MASTER_ADDR)
	log.Println("redis pwd ", RedisConfig.REDIS_PASSWORD)
}
