package logic

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
	"math/rand"
	"mongo_study/test03/models"
	"time"
)

// 学生信息设置
type StudentModel struct {
	ID        int       `json:"id" bson:"id"`
	Name      string    `json:"name" bson:"name"`
	Age       int       `json:"age" bson:"age"`
	Address   string    `json:"address" bson:"address"`
	CreatedAt TimeStamp `json:"created_at" bson:"created_at"`
	UpdatedAt TimeStamp `json:"updated_at" bson:"updated_at"`
	StartTime TimeStamp `json:"start_time" bson:"start_time"`
	EndTime   TimeStamp `json:"end_time" bson:"end_time"`
	AuditTime TimeStamp `json:"audit_time" bson:"audit_time"`
	State     int       `json:"state" bson:"state"`
}

type TimeStamp struct {
	time.Time
}

func (t *TimeStamp) UnmarshalParam(tStr string) error {
	ts, err := time.Parse("2006-01-02 15:04:05", tStr)
	*t = TimeStamp{ts}
	return err
}

func (t *TimeStamp) UnmarshalJson(data []byte) error {
	if len(data) < 2 {
		return errors.New("时间格式非法")
	}

	now, err := time.ParseInLocation("2006-01-02 15:04:05", string(data[1:len(data)-1]), time.Local)
	*t = TimeStamp{now}
	return err
}

func (t *TimeStamp) MarshalJson() ([]byte, error) {
	b := make([]byte, 0, len("2006-01-02 15:04:05")+2)
	b = append(b, '"')
	b = t.Local().AppendFormat(b, "2006-01-02 15:04:05")
	b = append(b, '"')
	return b, nil
}

func GetStudentInfo(ctx context.Context, id int) (*StudentModel, error) {
	cli, err := models.GetMongoCli("default")
	if err != nil {
		return nil, err
	}
	receiver := new(StudentModel)
	res := cli.Database("binshow").Collection("student", options.Collection().SetReadConcern(readconcern.Majority())).FindOne(ctx, map[string]interface{}{"id": id})
	err = res.Decode(&receiver)
	return receiver, err
}

func CreateStudentInfo(ctx context.Context, req *StudentModel) error {
	req.ID = int(rand.Uint32())
	cli, err := models.GetMongoCli("default")
	if err != nil {
		return err
	}

	res, err := cli.Database("binshow").Collection("student", options.Collection().SetWriteConcern(writeconcern.New(writeconcern.WMajority()))).InsertOne(ctx, req)
	fmt.Printf("插入成功：%v", res)
	return err

}

func UpdateStudentInfo(ctx context.Context, id int, update interface{}) error {

	cli, err := models.GetMongoCli("default")
	if err != nil {
		return err
	}

	res, err := cli.Database("binshow").Collection("student", options.Collection().SetWriteConcern(writeconcern.New(writeconcern.WMajority()))).UpdateOne(ctx, map[string]interface{}{"id": id}, update)
	fmt.Printf("UpdateStudentInfo : %v", res)
	return err
}

func DeleteStudentInfo(ctx context.Context, id int) error {

	cli, err := models.GetMongoCli("default")
	if err != nil {
		return err
	}
	res, err := cli.Database("binshow").Collection("student", options.Collection().SetReadConcern(readconcern.Majority())).DeleteOne(ctx, map[string]interface{}{"id": id})
	fmt.Printf("delete: %v", res)
	return err

}
