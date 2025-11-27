package transport

import (
	"github.com/gin-gonic/gin"

	"task/internal/service"
)

func NewRouter(mgr *service.Manager) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	registerLinksRoutes(r, mgr)
	registerReportRoutes(r, mgr)

	return r
}
