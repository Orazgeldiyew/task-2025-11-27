package transport

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"

	"task/internal/service"
)

func registerReportRoutes(r *gin.Engine, mgr *service.Manager) {
	r.POST("/links/report", func(c *gin.Context) {
		var req ReportRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
			return
		}

		if len(req.LinksList) == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "links_list must not be empty"})
			return
		}

		batches := mgr.GetBatches(req.LinksList)
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
				status := b.Status[link]

				pdf.Cell(20, 8, strconv.FormatInt(b.ID, 10))
				pdf.Cell(100, 8, truncate(link, 60))
				pdf.Cell(40, 8, string(status))
				pdf.Ln(8)
			}
		}

		c.Header("Content-Type", "application/pdf")
		c.Header("Content-Disposition", `attachment; filename="report.pdf"`)

		if err := pdf.Output(c.Writer); err != nil {
			log.Println("pdf output error:", err)
		}
	})
}

func truncate(s string, n int) string {
	rs := []rune(s)
	if len(rs) <= n {
		return s
	}
	return string(rs[:n-3]) + "..."
}
