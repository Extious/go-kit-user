package service

//提供具体的业务实现接口
import (
	"context"
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/longjoy/micro-go-course/section08/user/dao"
	"github.com/longjoy/micro-go-course/section08/user/redis"
	"log"
	"time"
)

// UserInfoDTO 用户介绍数据
type UserInfoDTO struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// RegisterUserVO 用户注册数据
type RegisterUserVO struct {
	Username string
	Password string
	Email    string
}

// 设置错误类型
var (
	ErrUserExisted = errors.New("user is existed")
	ErrPassword    = errors.New("email and password are not match")
	ErrRegistering = errors.New("email is registering")
)

type UserService interface {
	// Login 登录接口
	Login(ctx context.Context, email, password string) (*UserInfoDTO, error)
	// Register 注册接口
	Register(ctx context.Context, vo *RegisterUserVO) (*UserInfoDTO, error)
}

// UserServiceImpl 接口实现接收器
type UserServiceImpl struct {
	userDAO dao.UserDAO
}

// MakeUserServiceImpl 接收一个实现了UserDAO接口的实例，返回用于实现UserService接口的实例
func MakeUserServiceImpl(userDAO dao.UserDAO) UserService {
	return &UserServiceImpl{
		userDAO: userDAO,
	}
}

// Login 接收器实现Login接口
func (userService *UserServiceImpl) Login(ctx context.Context, email, password string) (*UserInfoDTO, error) {

	user, err := userService.userDAO.SelectByEmail(email)
	if err == nil {
		if user.Password == password {
			return &UserInfoDTO{
				ID:       user.ID,
				Username: user.Username,
				Email:    user.Email,
			}, nil
		} else {
			return nil, ErrPassword
		}
	} else {
		log.Printf("err : %s", err)
	}
	return nil, err
}

// Register 接收器实现Register接口
func (userService UserServiceImpl) Register(ctx context.Context, vo *RegisterUserVO) (*UserInfoDTO, error) {

	lock := redis.GetRedisLock(vo.Email, time.Duration(5)*time.Second)
	err := lock.Lock()
	if err != nil {
		log.Printf("err : %s", err)
		return nil, ErrRegistering
	}
	defer lock.Unlock()

	existUser, err := userService.userDAO.SelectByEmail(vo.Email)

	if (err == nil && existUser == nil) || err == gorm.ErrRecordNotFound {
		newUser := &dao.UserEntity{
			Username: vo.Username,
			Password: vo.Password,
			Email:    vo.Email,
		}
		err = userService.userDAO.Save(newUser)
		if err == nil {
			return &UserInfoDTO{
				ID:       newUser.ID,
				Username: newUser.Username,
				Email:    newUser.Email,
			}, nil
		}
	}
	if err == nil {
		err = ErrUserExisted
	}
	return nil, err

}
