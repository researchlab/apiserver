package model

import (
	"apiserver/pkg/auth"
	"apiserver/pkg/constvar"
	"fmt"

	"github.com/jinzhu/gorm"
)

// User represents a registered user.

type UserModel struct {
	//	BaseModel
	gorm.Model
	Username string `json:"username" gorm:"column:username; not null" binding:"required" validate:"min=1,max=32"`
	Password string `json:"password" gorm:"column:password;not null" binding:"required" validate:"min=5, max=128"`
}

func (c *UserModel) TableName() string {
	return "tb_users"
}

//Create creates a new user account.

func (u *UserModel) Create() error {
	return DB.Self.Create(&u).Error
}

//DeleteUser deletes the user by the user identifier.
func DeleteUser(id uint) error {
	user := UserModel{}
	//	user.BaseModel.Id = id
	user.ID = id
	return DB.Self.Delete(&user).Error
}

//Update updates an user account information.
func (u *UserModel) Update() error {
	return DB.Self.Save(u).Error
}

//GetUser gets an user by the user identifier
func GetUser(username string) (*UserModel, error) {
	u := &UserModel{}
	d := DB.Self.Where("username = ?", username).First(&u)
	return u, d.Error
}

func GetUserByID(id int) (*UserModel, error) {
	u := &UserModel{}
	d := DB.Self.Where("id = ?", id).First(&u)
	return u, d.Error
}

//ListUser List all users
func ListUser(username string, offset, limit int) ([]*UserModel, uint, error) {
	if limit == 0 {
		limit = constvar.DefaultLimit
	}

	users := make([]*UserModel, 0)
	var count uint

	where := fmt.Sprintf("username like '%%%s%%'", username)
	if err := DB.Self.Model(&UserModel{}).Where(where).Count(&count).Error; err != nil {
		return users, count, err
	}

	if err := DB.Self.Where(where).Offset(offset).Limit(limit).Order("id desc").Find(&users).Error; err != nil {
		return users, count, err
	}
	return users, count, nil
}

func (u *UserModel) Compare(pwd string) error {
	return auth.Compare(u.Password, pwd)
}

//Encrypt the user password
func (u *UserModel) Encrypt() (err error) {
	u.Password, err = auth.Encrypt(u.Password)
	return
}
