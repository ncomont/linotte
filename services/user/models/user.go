package models

import "github.com/jinzhu/gorm"

type User struct {
	ID       uint   `gorm:"primary_key"`
	Username string `gorm:"type:varchar(1024);unique" json:"username"`
	Password string `gorm:"type:varchar(1024)" json:"password"`
}

type UserAccessor struct {
	*gorm.DB
}

func NewUserAccessor(db *gorm.DB) *UserAccessor {
	return &UserAccessor{db}
}

func (user User) Validate() bool {
	return len(user.Username) > 0 && len(user.Password) > 0
}

func (accessor *UserAccessor) GetByUsername(username string) User {
	user := User{}
	accessor.DB.Where("username LIKE ?", username).First(&user)
	return user
}

func (accessor *UserAccessor) Create(username string, hash string) (User, error) {
	user := User{
		Username: username,
		Password: hash,
	}
	err := accessor.DB.Create(&user).Error

	return user, err
}
