package main

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

// 连接本地MongoDB
func main() {
	////1. 设置客户端连接的配置
	//clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	//
	////2. 开启一个context来连接 MongoDB
	//client, err := mongo.Connect(context.TODO(), clientOptions)
	//if err != nil {
	//	fmt.Printf("connect mongodb err: %v \n" , err)
	//	log.Fatal(err)
	//}
	//
	////3. 检查连接:简单的发送一个ping命令
	//err = client.Ping(context.TODO(), nil)
	//if err != nil {
	//	fmt.Printf("ping mongodb err: %v \n" , err)
	//	log.Fatal(err)
	//}
	//
	//fmt.Println("connected to mongodb done!")
	//
	////4. 开始操作数据集
	//collection := client.Database("binshow").Collection("user")
	//println("当前数据集的数据库为：", collection.Database().Name())
	//println("当前数据集的名称为：", collection.Name())
	//
	////5. 记得要断开连接
	//
	//err = client.Disconnect(context.TODO())
	//if err != nil {
	//	fmt.Printf("disconnect mongodb err: %v \n" , err)
	//	log.Fatal(err)
	//}
	//fmt.Println("disconnect mongodb done!")

	mongo, err := ConnectToMongo("mongodb://localhost:27017", "binshow", 1000, 500)
	if err != nil {
		log.Fatal(err)
	}

	println(mongo.Name())

}

// 连接池模式
func ConnectToMongo(uri, name string, timeout time.Duration, num uint64) (*mongo.Database, error) {
	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	defer cancelFunc()

	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetMaxPoolSize(num)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	return client.Database(name), nil

}
