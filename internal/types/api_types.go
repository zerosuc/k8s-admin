package types

import (
	"time"

	"github.com/zhufuyi/sponge/pkg/ggorm/query"
)

var _ time.Time

// Tip: suggested filling in the binding rules https://github.com/go-playground/validator in request struct fields tag.

// CreateApiRequest request params
type CreateApiRequest struct {
	Handle   string `json:"handle" binding:""`
	Title    string `json:"title" binding:""`
	Path     string `json:"path" binding:""`
	Type     string `json:"type" binding:""`
	Action   string `json:"action" binding:""`
	CreateBy int    `json:"createBy" binding:""`
	UpdateBy int    `json:"updateBy" binding:""`
}

// UpdateApiByIDRequest request params
type UpdateApiByIDRequest struct {
	ID uint64 `json:"id" binding:""` // uint64 id

	Handle   string `json:"handle" binding:""`
	Title    string `json:"title" binding:""`
	Path     string `json:"path" binding:""`
	Type     string `json:"type" binding:""`
	Action   string `json:"action" binding:""`
	CreateBy int    `json:"createBy" binding:""`
	UpdateBy int    `json:"updateBy" binding:""`
}

// ApiObjDetail detail
type ApiObjDetail struct {
	ID string `json:"id"` // convert to string id

	Handle    string    `json:"handle"`
	Title     string    `json:"title"`
	Path      string    `json:"path"`
	Type      string    `json:"type"`
	Action    string    `json:"action"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreateBy  int       `json:"createBy"`
	UpdateBy  int       `json:"updateBy"`
}

// CreateApiRespond only for api docs
type CreateApiRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		ID uint64 `json:"id"` // id
	} `json:"data"` // return data
}

// UpdateApiByIDRespond only for api docs
type UpdateApiByIDRespond struct {
	Result
}

// GetApiByIDRespond only for api docs
type GetApiByIDRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Api ApiObjDetail `json:"api"`
	} `json:"data"` // return data
}

// DeleteApiByIDRespond only for api docs
type DeleteApiByIDRespond struct {
	Result
}

// DeleteApisByIDsRequest request params
type DeleteApisByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// DeleteApisByIDsRespond only for api docs
type DeleteApisByIDsRespond struct {
	Result
}

// GetApiByConditionRequest request params
type GetApiByConditionRequest struct {
	query.Conditions
}

// GetApiByConditionRespond only for api docs
type GetApiByConditionRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Api ApiObjDetail `json:"api"`
	} `json:"data"` // return data
}

// ListApisByIDsRequest request params
type ListApisByIDsRequest struct {
	IDs []uint64 `json:"ids" binding:"min=1"` // id list
}

// ListApisByIDsRespond only for api docs
type ListApisByIDsRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Apis []ApiObjDetail `json:"apis"`
	} `json:"data"` // return data
}

// ListApisRequest request params
type ListApisRequest struct {
	query.Params
}

// ListApisRespond only for api docs
type ListApisRespond struct {
	Code int    `json:"code"` // return code
	Msg  string `json:"msg"`  // return information description
	Data struct {
		Apis []ApiObjDetail `json:"apis"`
	} `json:"data"` // return data
}
