package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateRoleRequest request params
type CreateRoleRequest struct {
	RoleID    int    `json:"roleId" binding:""`
	RoleName  string `json:"roleName" binding:""`
	Status    string `json:"status" binding:""`
	RoleKey   string `json:"roleKey" binding:""`
	RoleSort  int    `json:"roleSort" binding:""`
	Flag      string `json:"flag" binding:""`
	Remark    string `json:"remark" binding:""`
	Admin     string `json:"admin" binding:""`
	DataScope string `json:"dataScope" binding:""`
	CreateBy  int    `json:"createBy" binding:""`
	UpdateBy  int    `json:"updateBy" binding:""`
}

// UpdateRoleByIDRequest request params
type UpdateRoleByIDRequest struct {
	ID        uint64 `json:"roleId" binding:""`
	RoleName  string `json:"roleName" binding:""`
	Status    string `json:"status" binding:""`
	RoleKey   string `json:"roleKey" binding:""`
	RoleSort  int    `json:"roleSort" binding:""`
	Flag      string `json:"flag" binding:""`
	Remark    string `json:"remark" binding:""`
	Admin     string `json:"admin" binding:""`
	DataScope string `json:"dataScope" binding:""`
	CreateBy  int    `json:"createBy" binding:""`
	UpdateBy  int    `json:"updateBy" binding:""`
}

// RoleObjDetail detail
type RoleObjDetail struct {
	ID        string    `json:"roleId"`
	RoleName  string    `json:"roleName"`
	Status    string    `json:"status"`
	RoleKey   string    `json:"roleKey"`
	RoleSort  int       `json:"roleSort"`
	Flag      string    `json:"flag"`
	Remark    string    `json:"remark"`
	Admin     string    `json:"admin"`
	DataScope string    `json:"dataScope"`
	CreateBy  int       `json:"createBy"`
	UpdateBy  int       `json:"updateBy"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// CreateRoleRespond only for api docs
type CreateRoleRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateRoleByIDRespond only for api docs
type UpdateRoleByIDRespond struct {
	Result
}

// GetRoleByIDRespond only for api docs
type GetRoleByIDRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Role RoleObjDetail `json:"role"`
	} `json:"data"` // return data
}

// DeleteRoleByIDRespond only for api docs
type DeleteRoleByIDRespond struct {
	Result
}

// DeleteRolesByIDsRequest request params
type DeleteRolesByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// DeleteRolesByIDsRespond only for api docs
type DeleteRolesByIDsRespond struct {
	Result
}

// GetRoleByConditionRequest request params
type GetRoleByConditionRequest struct {
	query.Conditions
}

// GetRoleByConditionRespond only for api docs
type GetRoleByConditionRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Role RoleObjDetail `json:"role"`
	} `json:"data"` // return data
}

// ListRolesByIDsRequest request params
type ListRolesByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListRolesByIDsRespond only for api docs
type ListRolesByIDsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Roles []RoleObjDetail `json:"roles"`
	} `json:"data"` // return data
}

// ListRolesRequest request params
type ListRolesRequest struct {
	query.Params
}

// ListRolesRespond only for api docs
type ListRolesRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Roles []RoleObjDetail `json:"roles"`
	} `json:"data"` // return data
}
