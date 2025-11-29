package transport

import (
	"github.com/gin-gonic/gin"

	"task/internal/service"
	"task/internal/transport/handlers"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "task/internal/docs"
)

func NewRouter(mgr *service.Manager) *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	handlers.RegisterLinksRoutes(r, mgr)
	handlers.RegisterReportRoutes(r, mgr)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
