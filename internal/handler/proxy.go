package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/zhufuyi/sponge/pkg/gin/middleware"
	"github.com/zhufuyi/sponge/pkg/gin/response"
	"github.com/zhufuyi/sponge/pkg/logger"
	"go-admin/internal/dao"
	"go-admin/internal/ecode"
	"go-admin/internal/utils"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"

	"k8s.io/apimachinery/pkg/util/proxy"
	"k8s.io/client-go/rest"
)

var _ ApiHandler = (*apiHandler)(nil)

// ProxyHandler defining the handler interface
type ProxyHandler interface {
	Proxy(c *gin.Context)
}

type proxyHandler struct {
	iDao dao.ApiDao
}

// NewProxyHandler creating the handler interface
func NewProxyHandler() ProxyHandler {
	return &proxyHandler{
		iDao: nil,
	}
}

// Proxy 代理K8s的所有接口
// @Summary 代理K8s的所有接口
// @Description 代理K8s的所有接口
// @Tags 代理K8s的所有接口
// @Accept application/json
// @Produce application/json
// @Router /api/v1/proxy/api [get]
func (p *proxyHandler) Proxy(c *gin.Context) {
	//r := utils.NewResponse()
	config, err := utils.GetKubeConfig()
	if err != nil {
		//logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	transport, err := rest.TransportFor(config)
	if err != nil {
		//logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}
	target, err := parseTarget(*c.Request.URL, config.Host)
	if err != nil {
		logger.Warn("ShouldBindJSON error: ", logger.Err(err), middleware.GCtxRequestIDField(c))
		response.Error(c, ecode.InvalidParams)
		return
	}

	httpProxy := proxy.NewUpgradeAwareHandler(target, transport, false, false, nil)
	httpProxy.UpgradeTransport = proxy.NewUpgradeRequestRoundTripper(transport, transport)
	//httpProxy.ServeHTTP(c.Writer, c.Request)
	// 创建自定义的 ResponseWriter
	crw := NewCustomResponseWriter()
	// 将自定义的 ResponseWriter 传递给 httpProxy.ServeHTTP
	klog.Infoln(c.Request, target, transport)
	httpProxy.ServeHTTP(crw, c.Request)
	// 将缓冲区的内容赋值给 Response 的 Result 字段
	result := crw.buf.String()
	data := json.RawMessage(result)

	response.Success(c, data)
}

// CustomResponseWriter 自定义的 ResponseWriter
type CustomResponseWriter struct {
	buf      bytes.Buffer
	response http.Response
	header   http.Header
}

func NewCustomResponseWriter() *CustomResponseWriter {
	return &CustomResponseWriter{
		header: make(http.Header),
	}
}

// Header 实现 http.ResponseWriter 的 Header 方法
func (crw *CustomResponseWriter) Header() http.Header {
	return crw.header
}

// Write 实现 http.ResponseWriter 的 Write 方法
func (crw *CustomResponseWriter) Write(b []byte) (int, error) {
	return crw.buf.Write(b)
}

// WriteHeader 实现 http.ResponseWriter 的 WriteHeader 方法
func (crw *CustomResponseWriter) WriteHeader(statusCode int) {
	crw.response.StatusCode = statusCode
}
func parseTarget(target url.URL, host string) (*url.URL, error) {
	kubeURL, err := url.Parse(host)
	if err != nil {
		return nil, err
	}
	// TODO: 检查 URL 是否规范
	target.Path = target.Path[len("/api/v1/proxy"+"/"):]
	target.Host = kubeURL.Host
	target.Scheme = kubeURL.Scheme
	logrus.Infoln(target.Path, target.Host, target.Scheme)
	return &target, nil
}
