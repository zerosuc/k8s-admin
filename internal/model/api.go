package model

import (
	"github.com/zhufuyi/sponge/pkg/ggorm"
)

type Api struct {
	ggorm.Model `gorm:"embedded"` // embed id and time

	Handle   string `gorm:"column:handle;type:text" json:"handle"`
	Title    string `gorm:"column:title;type:text" json:"title"`
	Path     string `gorm:"column:path;type:text" json:"path"`
	Type     string `gorm:"column:type;type:text" json:"type"`
	Action   string `gorm:"column:action;type:text" json:"action"`
	CreateBy int    `gorm:"column:create_by;type:int(11)" json:"createBy"`
	UpdateBy int    `gorm:"column:update_by;type:int(11)" json:"updateBy"`
}

// TableName table name
func (m *Api) TableName() string {
	return "api"
}
