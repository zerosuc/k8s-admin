package routers

import (
	"github.com/gin-gonic/gin"

	"go-admin/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		roleRouter(group, handler.NewRoleHandler())
	})
}

func roleRouter(group *gin.RouterGroup, h handler.RoleHandler) {
	//group.Use(middleware.Auth()) // all of the following routes use jwt authentication
	// or group.Use(middleware.Auth(middleware.WithVerify(verify))) // token authentication

	group.POST("/role", h.Create)
	group.DELETE("/role/:id", h.DeleteByID)
	group.POST("/role/delete/ids", h.DeleteByIDs)
	group.PUT("/role/:id", h.UpdateByID)
	group.GET("/role/:id", h.GetByID)
	group.POST("/role/condition", h.GetByCondition)
	group.POST("/role/list/ids", h.ListByIDs)
	group.GET("/role/list", h.ListByLastID)
	group.POST("/role/list", h.List)
}
