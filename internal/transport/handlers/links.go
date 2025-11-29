package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"task/internal/service"
)

type LinksHandler struct {
	manager *service.Manager
}

func RegisterLinksRoutes(r *gin.Engine, mgr *service.Manager) {
	h := &LinksHandler{manager: mgr}
	r.POST("/links", h.CreateLinks)
}

// CreateLinks godoc
// @Summary      Add links batch
// @Tags         links
// @Accept       json
// @Produce      json
// @Param        request  body      LinksRequest   true  "Links batch"
// @Success      200      {object}  LinksResponse
// @Failure      400      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /links [post]
func (h *LinksHandler) CreateLinks(c *gin.Context) {
	var req LinksRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	if len(req.Links) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "links must not be empty"})
		return
	}

	batch, err := h.manager.CreateBatch(req.Links)
	if err != nil {
		log.Println("CreateBatch:", err)
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
}
