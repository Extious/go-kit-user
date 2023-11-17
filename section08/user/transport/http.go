package transport

//指定项目提供服务的方式，比如HTTP或gRPC等
import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/longjoy/micro-go-course/section08/user/endpoint"
	"net/http"
	"os"
)

var (
	ErrorBadRequest = errors.New("invalid request parameter") //定义新错误
)

// MakeHttpHandler make http handler use mux
func MakeHttpHandler(ctx context.Context, endpoints *endpoint.UserEndpoints) http.Handler {
	r := mux.NewRouter()                                     //创建路由
	kitLog := log.NewLogfmtLogger(os.Stderr)                 //设置日志记录
	kitLog = log.With(kitLog, "ts", log.DefaultTimestampUTC) //增加日志字段：时间戳
	kitLog = log.With(kitLog, "caller", log.DefaultCaller)   //增加日志字段：调用者

	options := []kithttp.ServerOption{ //配置http服务器
		kithttp.ServerErrorHandler(transport.NewLogErrorHandler(kitLog)), //创建日志处理器
		kithttp.ServerErrorEncoder(encodeError),                          //错误编码函数，将错误信息以指定的格式返回给客户端
	}
	r.Methods("POST").Path("/register").Handler(kithttp.NewServer(
		endpoints.RegisterEndpoint,
		decodeRegisterRequest,
		encodeJSONResponse,
		options...,
	))
	r.Methods("POST").Path("/login").Handler(kithttp.NewServer(
		endpoints.LoginEndpoint,
		decodeLoginRequest,
		encodeJSONResponse,
		options...,
	))
	return r
}
func decodeRegisterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")

	if username == "" || password == "" || email == "" {
		return nil, ErrorBadRequest
	}
	return &endpoint.RegisterRequest{
		Username: username,
		Password: password,
		Email:    email,
	}, nil
}

func decodeLoginRequest(_ context.Context, r *http.Request) (interface{}, error) {
	email := r.FormValue("email")
	password := r.FormValue("password")

	if email == "" || password == "" {
		return nil, ErrorBadRequest
	}
	return &endpoint.LoginRequest{
		Email:    email,
		Password: password,
	}, nil
}

func encodeJSONResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}
