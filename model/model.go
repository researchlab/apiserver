package model

import (
	"sync"
	"time"
)

type BaseModel struct {
	Id        uint64 `gorm:"primary_key;AUTO_INCREMENT;column:id" json:"-"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type UserInfo struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	SayHello  string `json:"sayHello"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type UserList struct {
	Lock  *sync.Mutex
	IdMap map[uint]*UserInfo
}

// Token represents a JSON web token.
type Token struct {
	Token string `json:"token"`
}
