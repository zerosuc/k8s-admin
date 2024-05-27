package handler

import (
	"errors"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"

	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"github.com/zhufuyi/sponge/pkg/utils"

	"go-admin/internal/cache"
	"go-admin/internal/dao"
	"go-admin/internal/ecode"
	"go-admin/internal/model"
	"go-admin/internal/types"
)

var _ UserHandler = (*userHandler)(nil)

// UserHandler defining the handler interface
type UserHandler interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	DeleteByID(c *gin.Context)
	DeleteByIDs(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	ListByLastID(c *gin.Context)
	List(c *gin.Context)
}

type userHandler struct {
	iDao dao.UserDao
}

// NewUserHandler creating the handler interface
func NewUserHandler() UserHandler {
	return &userHandler{
		iDao: dao.NewUserDao(
			model.GetDB(),
			cache.NewUserCache(model.GetCacheType()),
		),
	}
}

// Login
// @Summary Login api
// @Description Login information
// @Tags user
// @accept json
// @Produce json
// @Param data body types.LoginRequest true "user information"
// @Success 200 {object} types.CreateUserRespond{}
// @Router /api/v1/login [post]
// @Security BearerAuth
func (h *userHandler) Login(c *gin.Context) {
	form := &types.LoginRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	user := &model.User{}
	err = copier.Copy(user, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateUser)
		return
	}
	ctx := middleware.WrapCtx(c)
	userInfo, err := h.iDao.GetByName(ctx, user.Name)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}
	if userInfo.Password != user.Password {
		logger.Error(ecode.ErrLogin.Msg())
		response.Error(c, ecode.ErrLogin)
		return
	}

	response.Success(c, userInfo)
}

// Register
// @Summary Register api
// @Description submit information to create user
// @Tags user
// @accept json
// @Produce json
// @Param data body types.CreateUserRequest true "user information"
// @Success 200 {object} types.CreateUserRespond{}
// @Router /api/v1/reg [post]
// @Security BearerAuth
func (h *userHandler) Register(c *gin.Context) {
	form := &types.CreateUserRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	user := &model.User{}
	err = copier.Copy(user, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateUser)
		return
	}

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, user)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": user.ID})
}

// DeleteByID delete a record by id
// @Summary delete user
// @Description delete user by id
// @Tags user
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteUserByIDRespond{}
// @Router /api/v1/user/{id} [delete]
// @Security BearerAuth
func (h *userHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getUserIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err := h.iDao.DeleteByID(ctx, id)
	if err != nil {
		logger.Error("DeleteByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// DeleteByIDs delete records by batch id
// @Summary delete users
// @Description delete users by batch id
// @Tags user
// @Param data body types.DeleteUsersByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteUsersByIDsRespond{}
// @Router /api/v1/user/delete/ids [post]
// @Security BearerAuth
func (h *userHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteUsersByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	err = h.iDao.DeleteByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// UpdateByID update information by id
// @Summary update user
// @Description update user information by id
// @Tags user
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateUserByIDRequest true "user information"
// @Success 200 {object} types.UpdateUserByIDRespond{}
// @Router /api/v1/user/{id} [put]
// @Security BearerAuth
func (h *userHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getUserIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateUserByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	user := &model.User{}
	err = copier.Copy(user, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDUser)
		return
	}

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, user)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get user detail
// @Description get user detail by id
// @Tags user
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUserByIDRespond{}
// @Router /api/v1/user/{id} [get]
// @Security BearerAuth
func (h *userHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getUserIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	user, err := h.iDao.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			logger.Warn("GetByID not found", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByID error", logger.Err(err), logger.Any("id", id), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.UserObjDetail{}
	err = copier.Copy(data, user)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDUser)
		return
	}
	data.ID = idStr

	response.Success(c, gin.H{"user": data})
}

// GetByCondition get a record by condition
// @Summary get user by condition
// @Description get user by condition
// @Tags user
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetUserByConditionRespond{}
// @Router /api/v1/user/condition [post]
// @Security BearerAuth
func (h *userHandler) GetByCondition(c *gin.Context) {
	form := &types.GetUserByConditionRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	err = form.Conditions.CheckValid()
	if err != nil {
		logger.Warn("Parameters error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	user, err := h.iDao.GetByCondition(ctx, &form.Conditions)
	if err != nil {
		if errors.Is(err, model.ErrRecordNotFound) {
			logger.Warn("GetByCondition not found", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Error(c, ecode.NotFound)
		} else {
			logger.Error("GetByCondition error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
			response.Output(c, ecode.InternalServerError.ToHTTPCode())
		}
		return
	}

	data := &types.UserObjDetail{}
	err = copier.Copy(data, user)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDUser)
		return
	}
	data.ID = utils.Uint64ToStr(user.ID)

	response.Success(c, gin.H{"user": data})
}

// ListByIDs list of records by batch id
// @Summary list of users by batch id
// @Description list of users by batch id
// @Tags user
// @Param data body types.ListUsersByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListUsersByIDsRespond{}
// @Router /api/v1/user/list/ids [post]
// @Security BearerAuth
func (h *userHandler) ListByIDs(c *gin.Context) {
	form := &types.ListUsersByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	userMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	users := []*types.UserObjDetail{}
	for _, id := range form.IDs {
		if v, ok := userMap[id]; ok {
			record, err := convertUser(v)
			if err != nil {
				response.Error(c, ecode.ErrListUser)
				return
			}
			users = append(users, record)
		}
	}

	response.Success(c, gin.H{
		"users": users,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of users by last id and limit
// @Description list of users by last id and limit
// @Tags user
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "size in each page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListUsersRespond{}
// @Router /api/v1/user/list [get]
// @Security BearerAuth
func (h *userHandler) ListByLastID(c *gin.Context) {
	lastID := utils.StrToUint64(c.Query("lastID"))
	if lastID == 0 {
		lastID = math.MaxInt32
	}
	limit := utils.StrToInt(c.Query("limit"))
	if limit == 0 {
		limit = 10
	}
	sort := c.Query("sort")

	ctx := middleware.WrapCtx(c)
	users, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUsers(users)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDUser)
		return
	}

	response.Success(c, gin.H{
		"users": data,
	})
}

// List of records by query parameters
// @Summary list of users by query parameters
// @Description list of users by paging and conditions
// @Tags user
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListUsersRespond{}
// @Router /api/v1/user/list [post]
// @Security BearerAuth
func (h *userHandler) List(c *gin.Context) {
	form := &types.ListUsersRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	users, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertUsers(users)
	if err != nil {
		response.Error(c, ecode.ErrListUser)
		return
	}

	response.Success(c, gin.H{
		"users": data,
		"total": total,
	})
}

func getUserIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertUser(user *model.User) (*types.UserObjDetail, error) {
	data := &types.UserObjDetail{}
	err := copier.Copy(data, user)
	if err != nil {
		return nil, err
	}
	data.ID = utils.Uint64ToStr(user.ID)
	return data, nil
}

func convertUsers(fromValues []*model.User) ([]*types.UserObjDetail, error) {
	toValues := []*types.UserObjDetail{}
	for _, v := range fromValues {
		data, err := convertUser(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
