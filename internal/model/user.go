package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type User struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	Name     string `gorm:"column:name;type:char(50);NOT NULL" json:"name"`                   // username
	Password string `gorm:"column:password;type:char(100);NOT NULL" json:"password"`          // password
	Email    string `gorm:"column:email;type:char(50);NOT NULL" json:"email"`                 // email
	Phone    string `gorm:"column:phone;type:char(30);NOT NULL" json:"phone"`                 // phone number
	Avatar   string `gorm:"column:avatar;type:varchar(200)" json:"avatar"`                    // avatar
	Age      int    `gorm:"column:age;type:tinyint(4);NOT NULL" json:"age"`                   // age
	Gender   int    `gorm:"column:gender;type:tinyint(4);NOT NULL" json:"gender"`             // gender, 1:Male, 2:Female, other values:unknown
	Status   int    `gorm:"column:status;type:tinyint(4);NOT NULL" json:"status"`             // account status, 1:inactive, 2:activated, 3:blocked
	LoginAt  uint64 `gorm:"column:login_at;type:bigint(20) unsigned;NOT NULL" json:"loginAt"` // login timestamp
}

// TableName table name
func (m *User) TableName() string {
	return "user"
}
