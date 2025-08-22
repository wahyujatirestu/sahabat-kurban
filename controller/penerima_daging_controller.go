package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/service"
)

type PenerimaDagingController struct {
	service service.PenerimaDagingService
}

func NewPenerimaDagingController(service service.PenerimaDagingService) *PenerimaDagingController {
	return &PenerimaDagingController{service: service}
}

// Create godoc
// @Summary Create penerima daging
// @Description Membuat data penerima daging baru
// @Tags Penerima Daging
// @Accept json
// @Produce json
// @Param request body dto.CreatePenerimaRequest true "Create Penerima Daging Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /penerima [post]
// @Security BearerAuth
func (c *PenerimaDagingController) Create(ctx *gin.Context) {
	var req dto.CreatePenerimaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	res, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"status": 201,
		"data": res,
		"message": "Penerima daging created successfully",
	})
}

// GetAll godoc
// @Summary Get all penerima daging
// @Description Mengambil semua data penerima daging
// @Tags Penerima Daging
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /penerima [get]
// @Security BearerAuth
func (c *PenerimaDagingController) GetAll(ctx *gin.Context) {
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
		"message": "Penerima daging retrieved successfully",
	})
}

// GetByID godoc
// @Summary Get penerima daging by ID
// @Description Mengambil detail penerima daging berdasarkan ID
// @Tags Penerima Daging
// @Produce json
// @Param id path string true "Penerima Daging ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /penerima/{id} [get]
// @Security BearerAuth
func (c *PenerimaDagingController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid id"})
		return
	}

	res, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": res,
		"message": "Penerima daging retrieved successfully",
	})
}

// Update godoc
// @Summary Update penerima daging
// @Description Mengupdate data penerima daging berdasarkan ID
// @Tags Penerima Daging
// @Accept json
// @Produce json
// @Param id path string true "Penerima Daging ID"
// @Param request body dto.UpdatePenerimaRequest true "Update Penerima Daging Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /penerima/{id} [put]
// @Security BearerAuth
func (c *PenerimaDagingController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid id"})
		return
	}

	var req dto.UpdatePenerimaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	data, err := c.service.Update(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": data,
		"message": "Penerima daging updated successfully",
	})
}

// Delete godoc
// @Summary Delete penerima daging
// @Description Menghapus data penerima daging berdasarkan ID
// @Tags Penerima Daging
// @Produce json
// @Param id path string true "Penerima Daging ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /penerima/{id} [delete]
// @Security BearerAuth
func (c *PenerimaDagingController) Delete(ctx *gin.Context) {
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
		"message": "Penerima daging deleted successfully",
	})
}