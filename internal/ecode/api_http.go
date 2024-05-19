package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// api business-level http error codes.
// the apiNO value range is 1~100, if the same number appears, it will cause a failure to start the service.
var (
	apiNO       = 31
	apiName     = "api"
	apiBaseCode = errcode.HCode(apiNO)

	ErrCreateApi         = errcode.NewError(apiBaseCode+1, "failed to create "+apiName)
	ErrDeleteByIDApi     = errcode.NewError(apiBaseCode+2, "failed to delete "+apiName)
	ErrDeleteByIDsApi    = errcode.NewError(apiBaseCode+3, "failed to delete by batch ids "+apiName)
	ErrUpdateByIDApi     = errcode.NewError(apiBaseCode+4, "failed to update "+apiName)
	ErrGetByIDApi        = errcode.NewError(apiBaseCode+5, "failed to get "+apiName+" details")
	ErrGetByConditionApi = errcode.NewError(apiBaseCode+6, "failed to get "+apiName+" details by conditions")
	ErrListByIDsApi      = errcode.NewError(apiBaseCode+7, "failed to list by batch ids "+apiName)
	ErrListByLastIDApi   = errcode.NewError(apiBaseCode+8, "failed to list by last id "+apiName)
	ErrListApi           = errcode.NewError(apiBaseCode+9, "failed to list of "+apiName)
	// error codes are globally unique, adding 1 to the previous error code
)
