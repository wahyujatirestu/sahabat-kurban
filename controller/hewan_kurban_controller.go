package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/service"
)

type HewanKurbanController struct {
	service service.HewanKurbanService
}

func NewHewanKurbanController(service service.HewanKurbanService) *HewanKurbanController {
	return &HewanKurbanController{service: service}
}

// Create godoc
// @Summary Create Hewan Kurban
// @Description Tambahkan data hewan kurban baru
// @Tags HewanKurban
// @Accept json
// @Produce json
// @Param request body dto.CreateHewanKurbanRequest true "Hewan Kurban request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hewan-kurban [post]
func (c *HewanKurbanController) Create(ctx *gin.Context) {
	var req dto.CreateHewanKurbanRequest
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

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": res,
		"message": "Hewan kurban created successfully",
	})
}

// GetAll godoc
// @Summary Get all Hewan Kurban
// @Description Ambil semua data hewan kurban
// @Tags HewanKurban
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hewan-kurban [get]
func (c *HewanKurbanController) GetAll(ctx *gin.Context) {
	list, err := c.service.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"data": list,
		"message": "Hewan kurban retrieved successfully",
	})
}

// GetByID godoc
// @Summary Get Hewan Kurban by ID
// @Description Ambil detail hewan kurban berdasarkan ID
// @Tags HewanKurban
// @Produce json
// @Param id path string true "Hewan Kurban ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hewan-kurban/{id} [get]
func (c *HewanKurbanController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid id"})
		return
	}

	data, err := c.service.GetByID(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(404, gin.H{
			"status": 404,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": data,
		"message": "Hewan kuban retrieved successfully",
	})
}

// Update godoc
// @Summary Update Hewan Kurban
// @Description Ubah data hewan kurban berdasarkan ID
// @Tags HewanKurban
// @Accept json
// @Produce json
// @Param id path string true "Hewan Kurban ID"
// @Param request body dto.UpdateHewanKurbanRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hewan-kurban/{id} [put]
func (c *HewanKurbanController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid id"})
		return
	}

	var req dto.UpdateHewanKurbanRequest
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
		"message": "Hewan kuban updated successfully",
	})
}


// Delete godoc
// @Summary Delete Hewan Kurban
// @Description Hapus hewan kurban berdasarkan ID
// @Tags HewanKurban
// @Produce json
// @Param id path string true "Hewan Kurban ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /hewan-kurban/{id} [delete]
func (c *HewanKurbanController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid id"})
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
		"message": "Hewan kurban deleted successfully",
	})
}