package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/longjoy/micro-go-course/section08/user/dao"
	"github.com/longjoy/micro-go-course/section08/user/endpoint"
	"github.com/longjoy/micro-go-course/section08/user/redis"
	"github.com/longjoy/micro-go-course/section08/user/service"
	"github.com/longjoy/micro-go-course/section08/user/transport"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

func main() {

	var (
		/*
			增加整型命令行参数可以设置服务端口
			三个参数分别为：标志名称，默认值，说明文档
			使用方式：go run main.go --service.port=10086
		*/
		servicePort = flag.Int("service.port", 8888, "service port")
	)
	//解析命令行参数
	flag.Parse()
	//创建空上下文对象
	ctx := context.Background()
	//创建err类型的通道
	errChan := make(chan error)
	//连接MySQL
	err := dao.InitMysql("127.0.0.1", "3306", "root", "1234", "user")
	if err != nil {
		//将错误信息输出到stderr中并终止程序
		log.Fatal(err)
	}
	//连接redis
	err = redis.InitRedis("127.0.0.1", "6379", "")
	if err != nil {
		log.Fatal(err)
	}
	//创建一个实现UserService接口的实例
	userService := service.MakeUserServiceImpl(&dao.UserDAOImpl{})

	userEndpoints := &endpoint.UserEndpoints{
		endpoint.MakeRegisterEndpoint(userService),
		endpoint.MakeLoginEndpoint(userService),
	}

	r := transport.MakeHttpHandler(ctx, userEndpoints)

	go func() {
		errChan <- http.ListenAndServe(":"+strconv.Itoa(*servicePort), r)
	}()

	go func() {
		// 监控系统信号，等待 ctrl + c 系统信号通知服务关闭
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	error := <-errChan
	log.Println(error)

}
