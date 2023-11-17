package dao

import "time"

// UserEntity user实体
type UserEntity struct {
	ID        int64
	Username  string
	Password  string
	Email     string
	CreatedAt time.Time
}

// TableName 返回TableName方法
func (UserEntity) TableName() string {
	return "user"
}

type UserDAO interface {
	// SelectByEmail 根据email查找user
	SelectByEmail(email string) (*UserEntity, error)
	// Save 保存到数据库中
	Save(user *UserEntity) error
}

// UserDAOImpl 接口实现接收器
type UserDAOImpl struct {
}

func (userDAO *UserDAOImpl) SelectByEmail(email string) (*UserEntity, error) {
	user := &UserEntity{}
	err := db.Where("email = ?", email).First(user).Error
	return user, err
}

func (userDAO *UserDAOImpl) Save(user *UserEntity) error {
	return db.Create(user).Error
}
