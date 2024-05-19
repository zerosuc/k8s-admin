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

var _ RoleHandler = (*roleHandler)(nil)

// RoleHandler defining the handler interface
type RoleHandler interface {
	Create(c *gin.Context)
	DeleteByID(c *gin.Context)
	DeleteByIDs(c *gin.Context)
	UpdateByID(c *gin.Context)
	GetByID(c *gin.Context)
	GetByCondition(c *gin.Context)
	ListByIDs(c *gin.Context)
	ListByLastID(c *gin.Context)
	List(c *gin.Context)
}

type roleHandler struct {
	iDao dao.RoleDao
}

// NewRoleHandler creating the handler interface
func NewRoleHandler() RoleHandler {
	return &roleHandler{
		iDao: dao.NewRoleDao(
			model.GetDB(),
			cache.NewRoleCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create role
// @Description submit information to create role
// @Tags role
// @accept json
// @Produce json
// @Param data body types.CreateRoleRequest true "role information"
// @Success 200 {object} types.CreateRoleRespond{}
// @Router /api/v1/role [post]
// @Security BearerAuth
func (h *roleHandler) Create(c *gin.Context) {
	form := &types.CreateRoleRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	role := &model.Role{}
	err = copier.Copy(role, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateRole)
		return
	}

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, role)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": role.ID})
}

// DeleteByID delete a record by id
// @Summary delete role
// @Description delete role by id
// @Tags role
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteRoleByIDRespond{}
// @Router /api/v1/role/{id} [delete]
// @Security BearerAuth
func (h *roleHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getRoleIDFromPath(c)
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
// @Summary delete roles
// @Description delete roles by batch id
// @Tags role
// @Param data body types.DeleteRolesByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteRolesByIDsRespond{}
// @Router /api/v1/role/delete/ids [post]
// @Security BearerAuth
func (h *roleHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteRolesByIDsRequest{}
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
// @Summary update role
// @Description update role information by id
// @Tags role
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateRoleByIDRequest true "role information"
// @Success 200 {object} types.UpdateRoleByIDRespond{}
// @Router /api/v1/role/{id} [put]
// @Security BearerAuth
func (h *roleHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getRoleIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateRoleByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	role := &model.Role{}
	err = copier.Copy(role, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDRole)
		return
	}

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, role)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get role detail
// @Description get role detail by id
// @Tags role
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetRoleByIDRespond{}
// @Router /api/v1/role/{id} [get]
// @Security BearerAuth
func (h *roleHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getRoleIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	role, err := h.iDao.GetByID(ctx, id)
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

	data := &types.RoleObjDetail{}
	err = copier.Copy(data, role)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDRole)
		return
	}
	data.ID = idStr

	response.Success(c, gin.H{"role": data})
}

// GetByCondition get a record by condition
// @Summary get role by condition
// @Description get role by condition
// @Tags role
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetRoleByConditionRespond{}
// @Router /api/v1/role/condition [post]
// @Security BearerAuth
func (h *roleHandler) GetByCondition(c *gin.Context) {
	form := &types.GetRoleByConditionRequest{}
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
	role, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.RoleObjDetail{}
	err = copier.Copy(data, role)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDRole)
		return
	}
	data.ID = utils.Uint64ToStr(role.ID)

	response.Success(c, gin.H{"role": data})
}

// ListByIDs list of records by batch id
// @Summary list of roles by batch id
// @Description list of roles by batch id
// @Tags role
// @Param data body types.ListRolesByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListRolesByIDsRespond{}
// @Router /api/v1/role/list/ids [post]
// @Security BearerAuth
func (h *roleHandler) ListByIDs(c *gin.Context) {
	form := &types.ListRolesByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	roleMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	roles := []*types.RoleObjDetail{}
	for _, id := range form.IDs {
		if v, ok := roleMap[id]; ok {
			record, err := convertRole(v)
			if err != nil {
				response.Error(c, ecode.ErrListRole)
				return
			}
			roles = append(roles, record)
		}
	}

	response.Success(c, gin.H{
		"roles": roles,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of roles by last id and limit
// @Description list of roles by last id and limit
// @Tags role
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "size in each page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListRolesRespond{}
// @Router /api/v1/role/list [get]
// @Security BearerAuth
func (h *roleHandler) ListByLastID(c *gin.Context) {
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
	roles, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertRoles(roles)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDRole)
		return
	}

	response.Success(c, gin.H{
		"roles": data,
	})
}

// List of records by query parameters
// @Summary list of roles by query parameters
// @Description list of roles by paging and conditions
// @Tags role
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListRolesRespond{}
// @Router /api/v1/role/list [post]
// @Security BearerAuth
func (h *roleHandler) List(c *gin.Context) {
	form := &types.ListRolesRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	roles, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertRoles(roles)
	if err != nil {
		response.Error(c, ecode.ErrListRole)
		return
	}

	response.Success(c, gin.H{
		"roles": data,
		"total": total,
	})
}

func getRoleIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertRole(role *model.Role) (*types.RoleObjDetail, error) {
	data := &types.RoleObjDetail{}
	err := copier.Copy(data, role)
	if err != nil {
		return nil, err
	}
	data.ID = utils.Uint64ToStr(role.ID)
	return data, nil
}

func convertRoles(fromValues []*model.Role) ([]*types.RoleObjDetail, error) {
	toValues := []*types.RoleObjDetail{}
	for _, v := range fromValues {
		data, err := convertRole(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
