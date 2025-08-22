package controller

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/service"
)

type ReportController struct {
	svc service.ReportService
}

func NewReportController(s service.ReportService) *ReportController {
	return &ReportController{svc: s}
}

// GetLaporan godoc
// @Summary Mendapatkan ringkasan laporan kurban
// @Description Ambil total pekurban, hewan, penerima, distribusi, dan pembayaran
// @Tags Laporan
// @Produce json
// @Success 200 {object} dto.LaporanResponse
// @Failure 500 {object} map[string]string
// @Security BearerAuth	
// @Router /laporan [get]
func (h *ReportController) GetLaporan(c *gin.Context) {
	filter := parseQueryToFilter(c)
	ctx := c.Request.Context()
	
	result, err := h.svc.GetConsolidatedReport(ctx, filter)
	if err != nil {
		c.JSON(500, gin.H{
			"error":   "gagal mengambil laporan",
			"details": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"status":  200,
		"laporan": result,
		"message": "Laporan retrieved successfully",
	})
}

func parseQueryToFilter(c *gin.Context) model.ReportFilter {
	// gunakan binding Gin untuk tanggal; fallback manual parsing jika perlu
	var q dto.ReportQuery
	_ = c.ShouldBindQuery(&q)

	// validasi: tukar order jika start > end
	if q.StartDate != nil && q.EndDate != nil && q.StartDate.After(*q.EndDate) {
		// swap
		start := *q.EndDate
		end := *q.StartDate
		q.StartDate = &start
		q.EndDate = &end
	}

	// normalisasi ke UTC (opsional)
	norm := func(t *time.Time) *time.Time {
		if t == nil {
			return nil
		}
		u := t.UTC()
		return &u
	}

	return model.ReportFilter{
		StartDate: norm(q.StartDate),
		EndDate:   norm(q.EndDate),
	}
}