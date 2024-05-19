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

var _ ApiHandler = (*apiHandler)(nil)

// ApiHandler defining the handler interface
type ApiHandler interface {
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

type apiHandler struct {
	iDao dao.ApiDao
}

// NewApiHandler creating the handler interface
func NewApiHandler() ApiHandler {
	return &apiHandler{
		iDao: dao.NewApiDao(
			model.GetDB(),
			cache.NewApiCache(model.GetCacheType()),
		),
	}
}

// Create a record
// @Summary create api
// @Description submit information to create api
// @Tags api
// @accept json
// @Produce json
// @Param data body types.CreateApiRequest true "api information"
// @Success 200 {object} types.CreateApiRespond{}
// @Router /api/v1/api [post]
// @Security BearerAuth
func (h *apiHandler) Create(c *gin.Context) {
	form := &types.CreateApiRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	api := &model.Api{}
	err = copier.Copy(api, form)
	if err != nil {
		response.Error(c, ecode.ErrCreateApi)
		return
	}

	ctx := middleware.WrapCtx(c)
	err = h.iDao.Create(ctx, api)
	if err != nil {
		logger.Error("Create error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c, gin.H{"id": api.ID})
}

// DeleteByID delete a record by id
// @Summary delete api
// @Description delete api by id
// @Tags api
// @accept json
// @Produce json
// @Param id path string true "id"
// @Success 200 {object} types.DeleteApiByIDRespond{}
// @Router /api/v1/api/{id} [delete]
// @Security BearerAuth
func (h *apiHandler) DeleteByID(c *gin.Context) {
	_, id, isAbort := getApiIDFromPath(c)
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
// @Summary delete apis
// @Description delete apis by batch id
// @Tags api
// @Param data body types.DeleteApisByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.DeleteApisByIDsRespond{}
// @Router /api/v1/api/delete/ids [post]
// @Security BearerAuth
func (h *apiHandler) DeleteByIDs(c *gin.Context) {
	form := &types.DeleteApisByIDsRequest{}
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
// @Summary update api
// @Description update api information by id
// @Tags api
// @accept json
// @Produce json
// @Param id path string true "id"
// @Param data body types.UpdateApiByIDRequest true "api information"
// @Success 200 {object} types.UpdateApiByIDRespond{}
// @Router /api/v1/api/{id} [put]
// @Security BearerAuth
func (h *apiHandler) UpdateByID(c *gin.Context) {
	_, id, isAbort := getApiIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	form := &types.UpdateApiByIDRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	form.ID = id

	api := &model.Api{}
	err = copier.Copy(api, form)
	if err != nil {
		response.Error(c, ecode.ErrUpdateByIDApi)
		return
	}

	ctx := middleware.WrapCtx(c)
	err = h.iDao.UpdateByID(ctx, api)
	if err != nil {
		logger.Error("UpdateByID error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	response.Success(c)
}

// GetByID get a record by id
// @Summary get api detail
// @Description get api detail by id
// @Tags api
// @Param id path string true "id"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetApiByIDRespond{}
// @Router /api/v1/api/{id} [get]
// @Security BearerAuth
func (h *apiHandler) GetByID(c *gin.Context) {
	idStr, id, isAbort := getApiIDFromPath(c)
	if isAbort {
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	api, err := h.iDao.GetByID(ctx, id)
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

	data := &types.ApiObjDetail{}
	err = copier.Copy(data, api)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDApi)
		return
	}
	data.ID = idStr

	response.Success(c, gin.H{"api": data})
}

// GetByCondition get a record by condition
// @Summary get api by condition
// @Description get api by condition
// @Tags api
// @Param data body types.Conditions true "query condition"
// @Accept json
// @Produce json
// @Success 200 {object} types.GetApiByConditionRespond{}
// @Router /api/v1/api/condition [post]
// @Security BearerAuth
func (h *apiHandler) GetByCondition(c *gin.Context) {
	form := &types.GetApiByConditionRequest{}
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
	api, err := h.iDao.GetByCondition(ctx, &form.Conditions)
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

	data := &types.ApiObjDetail{}
	err = copier.Copy(data, api)
	if err != nil {
		response.Error(c, ecode.ErrGetByIDApi)
		return
	}
	data.ID = utils.Uint64ToStr(api.ID)

	response.Success(c, gin.H{"api": data})
}

// ListByIDs list of records by batch id
// @Summary list of apis by batch id
// @Description list of apis by batch id
// @Tags api
// @Param data body types.ListApisByIDsRequest true "id array"
// @Accept json
// @Produce json
// @Success 200 {object} types.ListApisByIDsRespond{}
// @Router /api/v1/api/list/ids [post]
// @Security BearerAuth
func (h *apiHandler) ListByIDs(c *gin.Context) {
	form := &types.ListApisByIDsRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	apiMap, err := h.iDao.GetByIDs(ctx, form.IDs)
	if err != nil {
		logger.Error("GetByIDs error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	apis := []*types.ApiObjDetail{}
	for _, id := range form.IDs {
		if v, ok := apiMap[id]; ok {
			record, err := convertApi(v)
			if err != nil {
				response.Error(c, ecode.ErrListApi)
				return
			}
			apis = append(apis, record)
		}
	}

	response.Success(c, gin.H{
		"apis": apis,
	})
}

// ListByLastID get records by last id and limit
// @Summary list of apis by last id and limit
// @Description list of apis by last id and limit
// @Tags api
// @accept json
// @Produce json
// @Param lastID query int true "last id, default is MaxInt32" default(0)
// @Param limit query int false "size in each page" default(10)
// @Param sort query string false "sort by column name of table, and the "-" sign before column name indicates reverse order" default(-id)
// @Success 200 {object} types.ListApisRespond{}
// @Router /api/v1/api/list [get]
// @Security BearerAuth
func (h *apiHandler) ListByLastID(c *gin.Context) {
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
	apis, err := h.iDao.GetByLastID(ctx, lastID, limit, sort)
	if err != nil {
		logger.Error("GetByLastID error", logger.Err(err), logger.Uint64("latsID", lastID), logger.Int("limit", limit), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertApis(apis)
	if err != nil {
		response.Error(c, ecode.ErrListByLastIDApi)
		return
	}

	response.Success(c, gin.H{
		"apis": data,
	})
}

// List of records by query parameters
// @Summary list of apis by query parameters
// @Description list of apis by paging and conditions
// @Tags api
// @accept json
// @Produce json
// @Param data body types.Params true "query parameters"
// @Success 200 {object} types.ListApisRespond{}
// @Router /api/v1/api/list [post]
// @Security BearerAuth
func (h *apiHandler) List(c *gin.Context) {
	form := &types.ListApisRequest{}
	err := c.ShouldBindJSON(form)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	ctx := middleware.WrapCtx(c)
	apis, total, err := h.iDao.GetByColumns(ctx, &form.Params)
	if err != nil {
		logger.Error("GetByColumns error", logger.Err(err), logger.Any("form", form), middleware.GCtxRequestIDField(c))
		response.Output(c, ecode.InternalServerError.ToHTTPCode())
		return
	}

	data, err := convertApis(apis)
	if err != nil {
		response.Error(c, ecode.ErrListApi)
		return
	}

	response.Success(c, gin.H{
		"apis":  data,
		"total": total,
	})
}

func getApiIDFromPath(c *gin.Context) (string, uint64, bool) {
	idStr := c.Param("id")
	id, err := utils.StrToUint64E(idStr)
	if err != nil || id == 0 {
		logger.Warn("StrToUint64E error: ", logger.String("idStr", idStr), middleware.GCtxRequestIDField(c))
		return "", 0, true
	}

	return idStr, id, false
}

func convertApi(api *model.Api) (*types.ApiObjDetail, error) {
	data := &types.ApiObjDetail{}
	err := copier.Copy(data, api)
	if err != nil {
		return nil, err
	}
	data.ID = utils.Uint64ToStr(api.ID)
	return data, nil
}

func convertApis(fromValues []*model.Api) ([]*types.ApiObjDetail, error) {
	toValues := []*types.ApiObjDetail{}
	for _, v := range fromValues {
		data, err := convertApi(v)
		if err != nil {
			return nil, err
		}
		toValues = append(toValues, data)
	}

	return toValues, nil
}
