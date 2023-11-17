package endpoint

//用于接收请求并返回响应，调用service方法实现业务逻辑，同时生成响应体
import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/longjoy/micro-go-course/section08/user/service"
)

// UserEndpoints 对于每一个服务接口，都使用一个抽象的endpoint来表示
type UserEndpoints struct {
	RegisterEndpoint endpoint.Endpoint
	LoginEndpoint    endpoint.Endpoint
}

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	UserInfo *service.UserInfoDTO `json:"user_info"`
}

func MakeLoginEndpoint(userService service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*LoginRequest)
		userInfo, err := userService.Login(ctx, req.Email, req.Password)
		return &LoginResponse{UserInfo: userInfo}, err

	}
}

type RegisterRequest struct {
	Username string
	Email    string
	Password string
}

type RegisterResponse struct {
	UserInfo *service.UserInfoDTO `json:"user_info"`
}

func MakeRegisterEndpoint(userService service.UserService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*RegisterRequest)
		userInfo, err := userService.Register(ctx, &service.RegisterUserVO{
			Username: req.Username,
			Password: req.Password,
			Email:    req.Email,
		})
		return &RegisterResponse{UserInfo: userInfo}, err

	}
}
