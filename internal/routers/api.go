package routers

import (
	"github.com/gin-gonic/gin"

	"go-admin/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		apiRouter(group, handler.NewApiHandler())
	})
}

func apiRouter(group *gin.RouterGroup, h handler.ApiHandler) {
	//group.Use(middleware.Auth()) // all of the following routes use jwt authentication
	// or group.Use(middleware.Auth(middleware.WithVerify(verify))) // token authentication

	group.POST("/api", h.Create)
	group.DELETE("/api/:id", h.DeleteByID)
	group.POST("/api/delete/ids", h.DeleteByIDs)
	group.PUT("/api/:id", h.UpdateByID)
	group.GET("/api/:id", h.GetByID)
	group.POST("/api/condition", h.GetByCondition)
	group.POST("/api/list/ids", h.ListByIDs)
	group.GET("/api/list", h.ListByLastID)
	group.POST("/api/list", h.List)
}
