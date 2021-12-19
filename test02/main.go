package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Student struct {
	Name string
	Age  int
}

func main() {

	s1 := Student{"小红", 12}
	s2 := Student{"小兰", 10}
	s3 := Student{"小黄", 11}

	//1. 连接客户端
	db, err := ConnectToMongo("mongodb://localhost:27017", "binshow", 1000, 500)
	if err != nil {
		log.Fatal(err)
	}

	//2. 创建 collection
	err = db.CreateCollection(context.Background(), "student")
	if err != nil {
		log.Fatal(err)
	}

	//3. 插入一个文档
	err = insertOne(db, s1)
	if err != nil {
		fmt.Printf("insert one document err : %v \n", err)
		log.Fatal(err)
	}

	//4. 插入多个文档
	stu := []interface{}{s2, s3}
	err = insertMany(db, stu)
	if err != nil {
		fmt.Printf("insert many document err : %v \n", err)
		log.Fatal(err)
	}

	//5. 更新文档
	filter := bson.D{{"name", "小兰"}} // 筛选文档
	fmt.Println("filter: ", filter)
	update := bson.D{ // 更新，给他增加了一岁---> 更新文档
		{"$inc", bson.D{
			{"age", 1},
		}},
	}
	fmt.Println("update: ", update)

	err = updateBson(db, filter, update)
	if err != nil {
		fmt.Printf("update document err : %v \n", err)
		log.Fatal(err)
	}

	//6. 查找 文档
	// 创建一个Student变量用来接收查询的结果
	var result Student
	err = findOne(db, filter, result)
	if err != nil {
		fmt.Printf("findOne document err : %v \n", err)
		log.Fatal(err)
	}

	//7. 查找多个文档

	//返回一个游标，游标提供了一个文档流，你可以通过它一次迭代和解码一个文档。当游标用完之后，应该关闭游标
	findOptions := options.Find()
	findOptions.SetLimit(2)

	// 定义一个切片用来存储查询结果
	var results []*Student

	// 把bson.D{{}}作为一个filter来匹配所有文档
	cur, err := db.Collection("binshow").Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 查找多个文档返回一个光标
	// 遍历游标允许我们一次解码一个文档
	for cur.Next(context.TODO()) {
		// 创建一个值，将单个文档解码为该值
		var elem Student
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}

	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// 完成后关闭游标
	cur.Close(context.TODO())
	fmt.Printf("Found multiple documents (array of pointers): %#v\n", results)

	// 8. 删除单个文档
	// 删除名字是小黄的那个
	deleteResult1, err := db.Collection("binshow").DeleteOne(context.TODO(), bson.D{{"name", "小黄"}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult1.DeletedCount)

	// 9. 删除所有文档
	deleteResult2, err := db.Collection("binshow").DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Deleted %v documents in the trainers collection\n", deleteResult2.DeletedCount)

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

func insertOne(db *mongo.Database, s1 Student) error {
	res, err := db.Collection("student").InsertOne(context.TODO(), s1)
	if err != nil {
		fmt.Printf("insert error=%v\n", err)
		log.Fatal(err)
	}
	fmt.Printf("insert a single document:%v ", res.InsertedID)
	return nil
}

func insertMany(db *mongo.Database, stu []interface{}) error {
	manyResult, err := db.Collection("student").InsertMany(context.TODO(), stu)
	if err != nil {
		fmt.Printf("insert error=%v\n", err)
		log.Fatal(err)
	}
	fmt.Printf("insert multiple document:%v\n ", manyResult.InsertedIDs)
	return nil
}

func updateBson(db *mongo.Database, filter bson.D, update bson.D) error {
	updateResult, err := db.Collection("student").UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return err
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)
	return nil
}

func findOne(db *mongo.Database, filter bson.D, result Student) error {
	err := db.Collection("student").FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return err
	}
	fmt.Printf("Found a single document: %v\n", result)
	return nil
}
