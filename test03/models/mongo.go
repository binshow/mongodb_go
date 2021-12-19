package models

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"mongo_study/test03/config"
	"sync"
	"time"
)

var (
	mongoClient sync.Map
)

func InitMongoClient() {
	for name, conf := range config.GetInstance().MongoConf {
		//fmt.Println(name)
		//fmt.Println(conf)
		cli := func() *mongo.Client {
			mongoConfig := conf
			opt := options.Client().ApplyURI(mongoConfig.Addr)
			opt.SetMaxPoolSize(uint64(mongoConfig.PoolMaxSize))
			opt.SetMinPoolSize(uint64(mongoConfig.PoolMinSize))
			opt.SetMaxConnIdleTime(time.Duration(mongoConfig.MaxConnIdleTime) * time.Second)

			client, err := mongo.NewClient(opt)
			if err != nil {
				panic(err)
			}
			if err := client.Connect(context.TODO()); err != nil {
				fmt.Println("connect failed")
			}
			return client
		}()
		mongoClient.Store(name, cli)
	}

}

func GetMongoCli(name string) (*mongo.Client, error) {
	cli, ok := mongoClient.Load(name)
	if !ok {
		return nil, fmt.Errorf("can't find connect named: %s", name)
	}
	return cli.(*mongo.Client), nil // cli.(*mongo.Client) 做一个类型判断
}
