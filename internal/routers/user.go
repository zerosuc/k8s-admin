package routers

import (
	"github.com/gin-gonic/gin"

	"go-admin/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		userRouter(group, handler.NewUserHandler())
	})
}

func userRouter(group *gin.RouterGroup, h handler.UserHandler) {
	//group.Use(middleware.Auth()) // all of the following routes use jwt authentication
	// or group.Use(middleware.Auth(middleware.WithVerify(verify))) // token authentication

	group.POST("/user/login", h.Login)
	group.POST("/user/reg", h.Register)
	group.DELETE("/user/:id", h.DeleteByID)
	group.POST("/user/delete/ids", h.DeleteByIDs)
	group.PUT("/user/:id", h.UpdateByID)
	group.GET("/user/:id", h.GetByID)
	group.POST("/user/condition", h.GetByCondition)
	group.POST("/user/list/ids", h.ListByIDs)
	group.GET("/user/list", h.ListByLastID)
	group.POST("/user/list", h.List)

}
