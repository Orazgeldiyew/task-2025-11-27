package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"

	"task/internal/service"
)

type ReportHandler struct {
	manager *service.Manager
}

func RegisterReportRoutes(r *gin.Engine, mgr *service.Manager) {
	h := &ReportHandler{manager: mgr}
	r.POST("/links/report", h.Create)
}

// @Summary      Generate report
// @Tags         report
// @Accept       json
// @Produce      application/pdf
// @Param        request  body      ReportRequest  true  "Report request"
// @Success      200      "PDF"
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Failure      500      {object}  map[string]string
// @Router       /links/report [post]
func (h *ReportHandler) Create(c *gin.Context) {
	var req ReportRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}
	if len(req.LinksList) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "links_list must not be empty"})
		return
	}

	batches := h.manager.GetBatches(req.LinksList)
	if len(batches) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no batches found"})
		return
	}

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(40, 10, "Links status report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(20, 8, "Batch")
	pdf.Cell(100, 8, "Link")
	pdf.Cell(40, 8, "Status")
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 11)

	for _, b := range batches {
		for _, link := range b.Links {
			pdf.Cell(20, 8, strconv.FormatInt(b.ID, 10))
			pdf.Cell(100, 8, cut(link, 60))
			pdf.Cell(40, 8, string(b.Status[link]))
			pdf.Ln(8)
		}
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", `attachment; filename="report.pdf"`)
	_ = pdf.Output(c.Writer)
}

func cut(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n-3]) + "..."
}
