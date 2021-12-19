package controller

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"mongo_study/test03/logic"
	"mongo_study/test03/resp"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func GetInfo(c echo.Context) error {

	id := c.Param("id")
	idStr, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusOK, resp.EmptyResp().WithErr(fmt.Errorf("id 错误：%s", id)))
	}

	info, err := logic.GetStudentInfo(c.Request().Context(), idStr)
	if err != nil {
		return c.JSON(http.StatusOK, resp.EmptyResp().WithErr(fmt.Errorf("读取mongoDB错误：%s", id)))
	}
	return c.JSON(http.StatusOK, resp.SuccessResp().WithData(info))

}

type CreateArgs struct {
	Data *logic.StudentModel
}

func Create(c echo.Context) error {

	args := new(logic.StudentModel)
	if err := c.Bind(&args); err != nil {
		return c.JSON(http.StatusOK, resp.EmptyResp().WithErr(err))
	}
	//cookies , err := c.Cookie("username")
	//if err != nil {
	//	return c.JSON(http.StatusOK , resp.EmptyResp().WithErr(fmt.Errorf("非法用户")))
	//}
	//
	//user := cookies.Value

	args.CreatedAt = logic.TimeStamp{time.Now()}

	err := logic.CreateStudentInfo(c.Request().Context(), args)
	if err != nil {
		return c.JSON(http.StatusOK, resp.EmptyResp().WithErr(err))
	}

	return c.JSON(http.StatusOK, resp.SuccessResp())

}

func Update(c echo.Context) error {

	args := new(logic.StudentModel)
	if err := c.Bind(&args); err != nil {
		return c.JSON(http.StatusOK, resp.EmptyResp().WithErr(err))
	}
	//args.CreatedAt = logic.TimeStamp{time.Now()}
	args.UpdatedAt = logic.TimeStamp{time.Now()}

	update := make(map[string]map[string]interface{})
	update["$set"] = map[string]interface{}{}
	fmt.Printf("NumField : %v\n", reflect.TypeOf(*args).NumField())

	lens := reflect.TypeOf(*args).NumField() // 6个字段
	value := reflect.ValueOf(*args)

	for i := 0; i < lens; i++ {
		tag := reflect.TypeOf(*args).Field(i).Tag.Get("bson")
		if tag == "" {
			continue
		}

		tag = strings.Split(tag, ",")[0]
		if tag == "id" {
			continue
		}

		if value.Field(i).IsZero() {
			continue
		}

		update["$set"][tag] = value.Field(i).Interface()
	}

	fmt.Println("update : %#v", update["$set"])
	err := logic.UpdateStudentInfo(c.Request().Context(), args.ID, update)
	if err != nil {
		return c.JSON(200, resp.EmptyResp().WithErr(fmt.Errorf("学生信息更新失败： %v", err)))
	}

	return c.JSON(200, resp.SuccessResp())

}

func Delete(c echo.Context) error {
	id := c.Param("id")
	idStr, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusOK, resp.EmptyResp().WithErr(fmt.Errorf("id 错误：%s", id)))
	}
	err = logic.DeleteStudentInfo(c.Request().Context(), idStr)
	if err != nil {
		return c.JSON(http.StatusOK, resp.EmptyResp().WithErr(fmt.Errorf("读取mongoDB错误：%s", id)))
	}
	return c.JSON(200, resp.SuccessResp())
}

func GetList(c echo.Context) error {
	type args struct {
		PageSize      int       `json:"page_size" bson:"page_size"`
		Page          int       `json:"page" bson:"page"`
		Name          string    `json:"name" bson:"name"`
		Age           int       `json:"age" bson:"age"`
		State         int       `json:"state" bson:"state"`
		StartDate     string    `json:"start_date"`
		EndDate       string    `json:"end_date"`
		StartDateTime time.Time `json:"-"`
		EndDateTime   time.Time `json:"-"`
	}

	arg := new(args)
	if err := c.Bind(arg); err != nil {
		return err
	}

	if arg.StartDate != "" && arg.EndDate != "" {
		StartDateTime, err := time.ParseInLocation("2006-01-02", arg.StartDate, time.Local)
		if err != nil {
			return c.JSON(200, resp.EmptyResp().WithErr(fmt.Errorf("查询参数错误：%v", err)))
		}

		EndDateTime, err := time.ParseInLocation("2006-01-02", arg.EndDate, time.Local)
		if err != nil {
			return c.JSON(200, resp.EmptyResp().WithErr(fmt.Errorf("查询参数错误：%v", err)))
		}

		arg.StartDateTime = StartDateTime
		arg.EndDateTime = EndDateTime
	}

	if arg.Page == 0 || arg.PageSize > 10 {
		arg.PageSize = 10
	}
	if arg.Page == 0 {
		arg.Page = 1
	}

	return nil
}
