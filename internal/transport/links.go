package transport

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"task/internal/service"
)

func registerLinksRoutes(r *gin.Engine, mgr *service.Manager) {
	r.POST("/links", func(c *gin.Context) {
		var req LinksRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		if len(req.Links) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "links must not be empty"})
			return
		}

		batch, err := mgr.CreateBatch(req.Links)
		if err != nil {
			log.Println("CreateBatch error:", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		resp := LinksResponse{
			Links:    make(map[string]string, len(batch.Status)),
			LinksNum: batch.ID,
		}

		for link, st := range batch.Status {
			resp.Links[link] = string(st)
		}

		c.JSON(http.StatusOK, resp)
	})
}
