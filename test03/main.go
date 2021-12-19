package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"mongo_study/test03/config"
	"mongo_study/test03/controller"
	"mongo_study/test03/models"
	"net/http"
	"strconv"
	"time"
)

const (
	ConfigFile = "./config/mongodb.yaml"
)

var (
	processQuit = make(chan interface{}, 5)
)

func main() {

	// 解析config
	cfg := config.GetInstance()
	if err := cfg.InitConfig(ConfigFile); err != nil {
		log.Fatalf("err = %s", err)
		return
	}

	models.InitMongoClient()

	echoInstance := GetEchoInstance()

	Dispatch(echoInstance)

	// 开启服务
	if err := StartService(echoInstance); err != nil {
		return
	}
	defer StopService(echoInstance)
	ClockLoop()
}

func StartService(echoInstance *echo.Echo) error {
	if echoInstance == nil {
		return nil
	}
	address := ":" + strconv.Itoa(8090)
	if err := echoInstance.Start(address); err != nil {
		return err
	}
	fmt.Println("echo server start success!")
	return nil
}

func StopService(echoInstance *echo.Echo) error {
	if echoInstance == nil {
		return nil
	}
	if err := echoInstance.Close(); err != nil {
		return err
	}
	fmt.Println("echo server close success!")
	return nil
}

func ClockLoop() {
	for true {
		select {
		case <-processQuit:
			fmt.Println("Ready to Quit Clock Loop")
			return
		default:
			time.Sleep(time.Second * 10)
		}
	}
}

func Dispatch(echoInstance *echo.Echo) {
	root := echoInstance.Group("/binshow")
	studentGroup := root.Group("/student")
	studentGroup.Match([]string{echo.POST}, "/create", controller.Create)  // url: /binshow/student/create
	studentGroup.Match([]string{echo.GET}, "/get/:id", controller.GetInfo) // url: /binshow/student/get/id
	studentGroup.Match([]string{echo.PUT}, "/update", controller.Update)   // url: /binshow/student/update
	// 一般情况下是不会删除数据库中的数据的，大部分都是改变其中的状态
	//studentGroup.Match([]string{echo.GET}, "/delete:id", controller.Delete) // url: /binshow/student/delete

}

func GetEchoInstance() *echo.Echo {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{

		AllowOrigins:     []string{"http://127.0.0.1"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut},
		AllowHeaders:     []string{},
		AllowCredentials: true,
		ExposeHeaders:    []string{},
		MaxAge:           86400,
	}))
	return e
}
