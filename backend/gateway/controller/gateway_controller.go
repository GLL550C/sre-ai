package controller

import (
	"gateway/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type GatewayController struct {
	proxyService *service.ProxyService
	logger       *zap.Logger
}

func NewGatewayController(proxyService *service.ProxyService, logger *zap.Logger) *GatewayController {
	return &GatewayController{
		proxyService: proxyService,
		logger:       logger,
	}
}

func (gc *GatewayController) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "gateway",
	})
}

func (gc *GatewayController) Metrics(c *gin.Context) {
	promhttp.Handler().ServeHTTP(c.Writer, c.Request)
}

func (gc *GatewayController) ProxyToCore(c *gin.Context) {
	gc.proxyService.ProxyRequest(c, "core")
}

func (gc *GatewayController) ProxyToRunbook(c *gin.Context) {
	gc.proxyService.ProxyRequest(c, "runbook")
}

func (gc *GatewayController) ProxyToTenant(c *gin.Context) {
	gc.proxyService.ProxyRequest(c, "tenant")
}
