package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// user business-level http error codes.
// the userNO value range is 1~100, if the same number appears, it will cause a failure to start the service.
var (
	userNO       = 64
	userName     = "user"
	userBaseCode = errcode.HCode(userNO)

	ErrCreateUser         = errcode.NewError(userBaseCode+1, "failed to create "+userName)
	ErrDeleteByIDUser     = errcode.NewError(userBaseCode+2, "failed to delete "+userName)
	ErrDeleteByIDsUser    = errcode.NewError(userBaseCode+3, "failed to delete by batch ids "+userName)
	ErrUpdateByIDUser     = errcode.NewError(userBaseCode+4, "failed to update "+userName)
	ErrGetByIDUser        = errcode.NewError(userBaseCode+5, "failed to get "+userName+" details")
	ErrGetByConditionUser = errcode.NewError(userBaseCode+6, "failed to get "+userName+" details by conditions")
	ErrListByIDsUser      = errcode.NewError(userBaseCode+7, "failed to list by batch ids "+userName)
	ErrListByLastIDUser   = errcode.NewError(userBaseCode+8, "failed to list by last id "+userName)
	ErrListUser           = errcode.NewError(userBaseCode+9, "failed to list of "+userName)

	ErrLogin = errcode.NewError(userBaseCode+10, "username or passwd error ")
	// error codes are globally unique, adding 1 to the previous error code
)
