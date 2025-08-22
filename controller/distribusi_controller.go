package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/service"
)

type DistribusiDagingController struct {
	service service.DistribusiDagingService
}

func NewDistribusiDagingController(service service.DistribusiDagingService) *DistribusiDagingController {
	return &DistribusiDagingController{service: service}
}

// Create godoc
// @Summary Create distribusi daging
// @Description Membuat distribusi daging baru
// @Tags DistribusiDaging
// @Accept json
// @Produce json
// @Param request body dto.CreateDistribusiRequest true "Distribusi daging request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 422 {object} map[string]interface{}
// @Security BearerAuth
// @Router /distribusi [post]
func (c *DistribusiDagingController) Create(ctx *gin.Context) {
	var req dto.CreateDistribusiRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	res, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(422, gin.H{
			"status": 422,
			"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"status": 201,
		"data": res,
		"message": "Distribusi daging created successfully",
	})
}

// GetAll godoc
// @Summary Get all distribusi daging
// @Description Mengambil semua distribusi daging
// @Tags DistribusiDaging
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /distribusi [get]
func (c *DistribusiDagingController) GetAll(ctx *gin.Context) {
	res, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": res,
		"message": "Distribusi daging retrieved successfully",
	})
}

// GetByID godoc
// @Summary Get distribusi daging by ID
// @Description Mengambil distribusi daging berdasarkan ID
// @Tags DistribusiDaging
// @Produce json
// @Param id path string true "Distribusi ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /distribusi/{id} [get]
func (c *DistribusiDagingController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid ID"})
		return
	}

	res, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(404, gin.H{
			"status": 404,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": res,
		"message": "Distribusi daging retrieved successfully",
	})
}


// GetTotalPaket godoc
// @Summary Get total paket distribusi
// @Description Mengambil total jumlah paket distribusi daging
// @Tags DistribusiDaging
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /distribusi/total-paket [get]
func (c *DistribusiDagingController) GetTotalPaket(ctx *gin.Context) {
	total, err := c.service.GetTotalDistribusiPaket(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"total_distribusi_paket": total,
		"message": "Distribusi daging retrieved successfully",
	})
}

// GetPenerimaBelumDistribusi godoc
// @Summary Get penerima yang belum menerima distribusi
// @Description Mengambil daftar penerima daging yang belum terdistribusi
// @Tags DistribusiDaging
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /distribusi/belum-terdistribusi [get]
func (c *DistribusiDagingController) GetPenerimaBelumDistribusi(ctx *gin.Context) {
	data, err := c.service.GetPenerimaBelumTerdistribusi(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"data": data,
		"message": "Distribusi daging retrieved successfully",
	})
}

// Delete godoc
// @Summary Delete distribusi daging
// @Description Menghapus distribusi daging berdasarkan ID
// @Tags DistribusiDaging
// @Produce json
// @Param id path string true "Distribusi ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /distribusi/{id} [delete]
func (c *DistribusiDagingController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid ID"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"message": "Distribusi daging deleted successfully",
	})
}