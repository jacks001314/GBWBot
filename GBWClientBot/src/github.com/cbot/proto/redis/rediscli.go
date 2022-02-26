package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisClient struct {

	client *redis.Client

}


func NewRedisClient(host string,port int,pass string,timeout uint64,db int) *RedisClient {

	rdb := redis.NewClient(&redis.Options{

		Addr:               fmt.Sprintf("%s:%d",host,port),
		Password:           pass,
		DialTimeout:        time.Duration(timeout)*time.Millisecond,
		DB: db,
	})

	return &RedisClient{client:rdb}
}

func (c *RedisClient) Info() (string,error) {
	return c.client.Info(context.Background()).Result()
}


func (c *RedisClient) Set(k string,v string) (string,error) {
	return c.client.Set(context.Background(),k,v,0).Result()
}

func (c *RedisClient) ConfigSet(key,value string)(string,error) {

	return c.client.ConfigSet(context.Background(),key,value).Result()
}

func (c *RedisClient) SlaveOfNoOne()(string,error) {

	return c.client.SlaveOf(context.Background(),"NO","ONE").Result()
}

func (c *RedisClient) Save()(string,error){

	return c.client.Save(context.Background()).Result()
}

func (c *RedisClient) Do(args ... interface{}) (interface{},error) {

	return c.client.Do(context.Background(),args...).Result()

}

func (c *RedisClient) Close(){

	c.client.Close()
}


