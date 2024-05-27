package routers

import (
	"github.com/gin-gonic/gin"
	"go-admin/internal/handler"
)

func init() {
	apiV1RouterFns = append(apiV1RouterFns, func(group *gin.RouterGroup) {
		proxyRouter(group, handler.NewProxyHandler())
	})
}

func proxyRouter(group *gin.RouterGroup, h handler.ProxyHandler) {
	//group.Use(middleware.Auth()) // all of the following routes use jwt authentication
	// or group.Use(middleware.Auth(middleware.WithVerify(verify))) // token authentication
	group = group.Group("/proxy")
	group.Any("/*path", h.Proxy)

}
