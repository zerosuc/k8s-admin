package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateUserRequest request params
type CreateUserRequest struct {
	Name     string `json:"name" binding:""`     // username
	Password string `json:"password" binding:""` // password
	Email    string `json:"email" binding:""`    // email
	Phone    string `json:"phone" binding:""`    // phone number
	Avatar   string `json:"avatar" binding:""`   // avatar
	Age      int    `json:"age" binding:""`      // age
	Gender   int    `json:"gender" binding:""`   // gender, 1:Male, 2:Female, other values:unknown
	Status   int    `json:"status" binding:""`   // account status, 1:inactive, 2:activated, 3:blocked
	LoginAt  uint64 `json:"loginAt" binding:""`  // login timestamp
}

// UpdateUserByIDRequest request params
type UpdateUserByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Name     string `json:"name" binding:""`     // username
	Password string `json:"password" binding:""` // password
	Email    string `json:"email" binding:""`    // email
	Phone    string `json:"phone" binding:""`    // phone number
	Avatar   string `json:"avatar" binding:""`   // avatar
	Age      int    `json:"age" binding:""`      // age
	Gender   int    `json:"gender" binding:""`   // gender, 1:Male, 2:Female, other values:unknown
	Status   int    `json:"status" binding:""`   // account status, 1:inactive, 2:activated, 3:blocked
	LoginAt  uint64 `json:"loginAt" binding:""`  // login timestamp
}

// UserObjDetail detail
type UserObjDetail struct {
	ID string `json:"id"` // convert to string id

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Name      string    `json:"name"`     // username
	Password  string    `json:"password"` // password
	Email     string    `json:"email"`    // email
	Phone     string    `json:"phone"`    // phone number
	Avatar    string    `json:"avatar"`   // avatar
	Age       int       `json:"age"`      // age
	Gender    int       `json:"gender"`   // gender, 1:Male, 2:Female, other values:unknown
	Status    int       `json:"status"`   // account status, 1:inactive, 2:activated, 3:blocked
	LoginAt   uint64    `json:"loginAt"`  // login timestamp
}

// CreateUserRespond only for api docs
type CreateUserRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateUserByIDRespond only for api docs
type UpdateUserByIDRespond struct {
	Result
}

// GetUserByIDRespond only for api docs
type GetUserByIDRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		User UserObjDetail `json:"user"`
	} `json:"data"` // return data
}

// DeleteUserByIDRespond only for api docs
type DeleteUserByIDRespond struct {
	Result
}

// DeleteUsersByIDsRequest request params
type DeleteUsersByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// DeleteUsersByIDsRespond only for api docs
type DeleteUsersByIDsRespond struct {
	Result
}

// GetUserByConditionRequest request params
type GetUserByConditionRequest struct {
	query.Conditions
}

// GetUserByConditionRespond only for api docs
type GetUserByConditionRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		User UserObjDetail `json:"user"`
	} `json:"data"` // return data
}

// ListUsersByIDsRequest request params
type ListUsersByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListUsersByIDsRespond only for api docs
type ListUsersByIDsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Users []UserObjDetail `json:"users"`
	} `json:"data"` // return data
}

// ListUsersRequest request params
type ListUsersRequest struct {
	query.Params
}

// ListUsersRespond only for api docs
type ListUsersRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Users []UserObjDetail `json:"users"`
	} `json:"data"` // return data
}
