package main

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v3"
	"github.com/casbin/redis-adapter/v3"
	"github.com/go-redis/redis/v8"
	"gorm.io/driver/sqlite"
)

func main() {
	f_redis()
}

// 初始化Redis客户端
func initRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return rdb
}

// 发送消息通知
func notify(rdb *redis.Client, channel string, message string) {
	err := rdb.Publish(context.Background(), channel, message).Err()
	if err != nil {
		log.Fatalf("Could not publish message: %v", err)
	}
}

func f1() {
	// 初始化GORM数据库
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}

	// 初始化Casbin适配器和Enforcer
	adapter, err := gormadapter.NewAdapterByDBUseTableName(db, "casbin", "rules")
	if err != nil {
		log.Fatalf("Could not create GORM adapter: %v", err)
	}

	enforcer, err := casbin.NewEnforcer("path/to/model.conf", adapter)
	if err != nil {
		log.Fatalf("Could not create enforcer: %v", err)
	}

	// 从DB加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		log.Fatalf("Could not load policy: %v", err)
	}

	// 添加一些策略
	enforcer.AddPolicy("alice", "data1", "read")
	enforcer.AddPolicy("bob", "data2", "write")

	// 保存策略到DB
	err = enforcer.SavePolicy()
	if err != nil {
		log.Fatalf("Could not save policy: %v", err)
	}

	// 初始化Redis客户端
	rdb := initRedisClient()

	// 检查权限并发送通知
	sub, obj, act := "alice", "data1", "read"
	ok, err := enforcer.Enforce(sub, obj, act)
	if err != nil {
		log.Fatalf("Could not enforce policy: %v", err)
	}

	if ok {
		fmt.Println("Permission granted")
		notify(rdb, "permissions", fmt.Sprintf("%s is allowed to %s %s", sub, act, obj))
	} else {
		fmt.Println("Permission denied")
		notify(rdb, "permissions", fmt.Sprintf("%s is not allowed to %s %s", sub, act, obj))
	}
}

// redis
func f_redis() {
	// 直接初始化：初始化 Redis 适配器并在 Casbin 执行器中使用它：
	// a, _ := redisadapter.NewAdapter("tcp", "127.0.0.1:6379") // Your Redis network and address.

	// 如果 Redis 的密码为“123”，请使用以下命令
	a, err := redisadapter.NewAdapterWithPassword("tcp", "192.168.6.41:30379", "36UdcuHtpZSfRj")
	if err != nil {
		log.Fatal(err)
	}

	// 如果您将 Redis 与特定用户一起使用，请使用以下命令
	// a, err := redisadapter.NewAdapterWithUser("tcp", "127.0.0.1:6379", "username", "password")

	// 如果使用 Redis 连接池，请使用以下命令
	// pool := &redis.Pool{}
	// a, err := redisadapter.NewAdapterWithPool(pool)

	//使用不同的用户选项进行初始化：
	//如果您将 Redis 与 passowrd 一起使用，例如“123”，请使用以下命令：
	// a, err := redisadapter.NewAdapterWithOption(redisadapter.WithNetwork("tcp"), redisadapter.WithAddress("127.0.0.1:6379"), redisadapter.WithPassword("123"))

	// 如果您将 Redis 与用户名、密码和 TLS 选项一起使用，请使用以下命令：
	// var clientTLSConfig tls.Config
	// ...
	// a, err := redisadapter.NewAdapterWithOption(redisadapter.WithNetwork("tcp"), redisadapter.WithAddress("127.0.0.1:6379"), redisadapter.WithUsername("testAccount"), redisadapter.WithPassword("123456"), redisadapter.WithTls(&clientTLSConfig))

	e, err := casbin.NewEnforcer("D:\\code\\goland\\casbin\\model.conf", a)
	if err != nil {
		log.Fatal(err)
	}

	// 从数据库加载策略。
	err = e.LoadPolicy()
	if err != nil {
		log.Fatal(err)
	}

	// 检查权限。
	ok, err := e.Enforce("alice", "data1", "read")
	if err != nil {
		log.Fatal(err)
	}
	if ok {
		log.Println("ok")
	} else {
		log.Println("no ok")
	}
	// 修改策略。
	// e.AddPolicy(...)
	// e.RemovePolicy(...)

	// 将策略保存回数据库。
	err = e.SavePolicy()
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Hour)
}

// https://learnku.com/articles/72226
