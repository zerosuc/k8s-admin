package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type Role struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	RoleID    int    `gorm:"column:role_id;type:int(11);primary_key" json:"roleId"`
	RoleName  string `gorm:"column:role_name;type:text" json:"roleName"`
	Status    string `gorm:"column:status;type:text" json:"status"`
	RoleKey   string `gorm:"column:role_key;type:text" json:"roleKey"`
	RoleSort  int    `gorm:"column:role_sort;type:int(11)" json:"roleSort"`
	Flag      string `gorm:"column:flag;type:text" json:"flag"`
	Remark    string `gorm:"column:remark;type:text" json:"remark"`
	Admin     string `gorm:"column:admin;type:decimal(10)" json:"admin"`
	DataScope string `gorm:"column:data_scope;type:text" json:"dataScope"`
	CreateBy  int    `gorm:"column:create_by;type:int(11)" json:"createBy"`
	UpdateBy  int    `gorm:"column:update_by;type:int(11)" json:"updateBy"`
}

// TableName table name
func (m *Role) TableName() string {
	return "role"
}
