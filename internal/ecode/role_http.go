package ecode

import (
	"github.com/zhufuyi/sponge/pkg/errcode"
)

// role business-level http error codes.
// the roleNO value range is 1~100, if the same number appears, it will cause a failure to start the service.
var (
	roleNO       = 90
	roleName     = "role"
	roleBaseCode = errcode.HCode(roleNO)

	ErrCreateRole         = errcode.NewError(roleBaseCode+1, "failed to create "+roleName)
	ErrDeleteByIDRole     = errcode.NewError(roleBaseCode+2, "failed to delete "+roleName)
	ErrDeleteByIDsRole    = errcode.NewError(roleBaseCode+3, "failed to delete by batch ids "+roleName)
	ErrUpdateByIDRole     = errcode.NewError(roleBaseCode+4, "failed to update "+roleName)
	ErrGetByIDRole        = errcode.NewError(roleBaseCode+5, "failed to get "+roleName+" details")
	ErrGetByConditionRole = errcode.NewError(roleBaseCode+6, "failed to get "+roleName+" details by conditions")
	ErrListByIDsRole      = errcode.NewError(roleBaseCode+7, "failed to list by batch ids "+roleName)
	ErrListByLastIDRole   = errcode.NewError(roleBaseCode+8, "failed to list by last id "+roleName)
	ErrListRole           = errcode.NewError(roleBaseCode+9, "failed to list of "+roleName)
	// error codes are globally unique, adding 1 to the previous error code
)
