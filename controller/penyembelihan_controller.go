package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/service"
)

type PenyembelihanController struct {
	service service.PenyembelihanService
}

func NewPenyembelihanController(s service.PenyembelihanService) *PenyembelihanController {
	return &PenyembelihanController{service: s}
}

// Create godoc
// @Summary Create penyembelihan
// @Description Membuat data penyembelihan baru
// @Tags Penyembelihan
// @Accept json
// @Produce json
// @Param request body dto.CreatePenyembelihanRequest true "Create Penyembelihan Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /penyembelihan [post]
// @Security BearerAuth
func (c *PenyembelihanController) Create(ctx *gin.Context) {
	var req dto.CreatePenyembelihanRequest
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
		"message": "Penyembelihan created successfully",
	})
}

// GetAll godoc
// @Summary Get all penyembelihan
// @Description Mengambil semua data penyembelihan
// @Tags Penyembelihan
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /penyembelihan [get]
// @Security BearerAuth
func (c *PenyembelihanController) GetAll(ctx *gin.Context) {
	res, err := c.service.GetAll(ctx)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": res,
		"message": "Penyembelihan retrieved successfully",
	})
}

// GetById godoc
// @Summary Get penyembelihan by ID
// @Description Mengambil detail penyembelihan berdasarkan ID
// @Tags Penyembelihan
// @Produce json
// @Param id path string true "Penyembelihan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /penyembelihan/{id} [get]
// @Security BearerAuth
func (c *PenyembelihanController) GetById(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid id"})
		return
	}

	res, err := c.service.GetById(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": res,
		"message": "Penyembelihan retrieved successfully",
	})
}

// Update godoc
// @Summary Update penyembelihan
// @Description Mengupdate data penyembelihan berdasarkan ID
// @Tags Penyembelihan
// @Accept json
// @Produce json
// @Param id path string true "Penyembelihan ID"
// @Param request body dto.UpdatePenyembelihanRequest true "Update Penyembelihan Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /penyembelihan/{id} [put]
// @Security BearerAuth
func (c *PenyembelihanController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid id"})
		return
	}

	var req dto.UpdatePenyembelihanRequest
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
		"message": "Penyembelihan updated successfully",
	})
}

// Delete godoc
// @Summary Delete penyembelihan
// @Description Menghapus data penyembelihan berdasarkan ID
// @Tags Penyembelihan
// @Produce json
// @Param id path string true "Penyembelihan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /penyembelihan/{id} [delete]
// @Security BearerAuth
func (c *PenyembelihanController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid ID"})
	}

	if err := c.service.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"message": "Penyembelihan deleted successfully"})
}