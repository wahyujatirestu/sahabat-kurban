package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/service"
)

type PekurbanHewanController struct {
	service service.PekurbanHewanService
	serv service.PekurbanService
}

func NewPekurbanHewanController(s service.PekurbanHewanService, serv service.PekurbanService) *PekurbanHewanController {
	return &PekurbanHewanController{service: s, serv: serv}
}

// @Summary Create new patungan hewan
// @Description User menambahkan dirinya ke dalam patungan hewan kurban
// @Tags Patungan
// @Accept json
// @Produce json
// @Param request body dto.CreatePekurbanHewanRequest true "Request Body"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patungan [post]
// @Security BearerAuth
func (c *PekurbanHewanController) Create(ctx *gin.Context) {
	userRaw, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(401, gin.H{
			"status": 401,
			"error": "Unauthorized"})
		return
	}
	currentUser := userRaw.(model.User)

	var req dto.CreatePekurbanHewanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	if currentUser.Role == "user" {
		p, err := c.serv.GetByUserId(ctx.Request.Context(), currentUser.ID)
		if err != nil || p == nil {
			ctx.JSON(403, gin.H{
				"status": 403,
				"error": "You have no registered pekurban data"})
			return
		}

		if 	p.ID != req.PekurbanID {
			ctx.JSON(403, gin.H{
				"status": 403,
				"error": "You can only register patungan for your own data"})
			return
		}
	}	

	data, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"status": 201,
		"data": data,
		"message": "Patungan added successfully", 
	})
}

// GetAll godoc
// @Summary Get all patungan
// @Description Ambil semua data patungan (admin/panitia), atau hanya milik user sendiri
// @Tags Patungan
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patungan [get]
// @Security BearerAuth
func (c *PekurbanHewanController) GetAll(ctx *gin.Context) {
	userRaw, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(401, gin.H{
			"status": 401,
			"error": "unauthorized"})
		return
	}
	currentUser := userRaw.(model.User)

	if currentUser.Role == "admin" || currentUser.Role == "panitia" {
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
			"message": "Pekurban retrieved successfully",
		})
		return
	}

	// Kalau user biasa, hanya ambil berdasarkan pekurban miliknya
	pekurban, err := c.serv.GetByUserId(ctx.Request.Context(), currentUser.ID)
	if err != nil || pekurban == nil {
		ctx.JSON(403, gin.H{
			"status": 403,
			"error": "Pekurban not found for this user"})
		return
	}

	pekurbanID, err := uuid.Parse(pekurban.ID)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": "Invalid UUID format for pekurban ID"})
		return
	}

	list, err := c.service.GetByPekurbanId(ctx.Request.Context(), pekurbanID)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{
		"status": 200,
		"data": list,
		"message": "Pekurban retrieved successfully",
	})
}

// GetByHewanID godoc
// @Summary Get patungan by hewan_id
// @Description Ambil semua pekurban berdasarkan ID hewan
// @Tags Patungan
// @Produce json
// @Param hewan_id path string true "Hewan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patungan/hewan/{hewan_id} [get]
// @Security BearerAuth
func (c *PekurbanHewanController) GetByHewanID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("hewan_id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "invalid hewan_id"})
		return
	}

	list, err := c.service.GetByHewanId(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": list,
		"message": "Pekurban retrieved successfully",
	})
}

// GetByPekurbanID godoc
// @Summary Get patungan by pekurban_id
// @Description Ambil semua data patungan berdasarkan ID pekurban
// @Tags Patungan
// @Produce json
// @Param pekurban_id path string true "Pekurban ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patungan/pekurban/{pekurban_id} [get]
// @Security BearerAuth
func (c *PekurbanHewanController) GetByPekurbanID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("pekurban_id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "invalid pekurban_id"})
		return
	}

	list, err := c.service.GetByPekurbanId(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": list,
		"message": "Pekurban retrieved successfully",
	})
}

// Update godoc
// @Summary Update patungan
// @Description Update porsi kontribusi user dalam patungan
// @Tags Patungan
// @Accept json
// @Produce json
// @Param pekurban_id path string true "Pekurban ID"
// @Param hewan_id path string true "Hewan ID"
// @Param request body dto.UpdatePekurbanHewanRequest true "Request Body"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patungan/{pekurban_id}/{hewan_id} [put]
// @Security BearerAuth
func (c *PekurbanHewanController) Update(ctx *gin.Context) {
	var req dto.UpdatePekurbanHewanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	pekurbanIDStr := ctx.Param("pekurban_id")
	hewanIDStr := ctx.Param("hewan_id")

	pekurbanID, err := uuid.Parse(pekurbanIDStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "invalid pekurban_id"})
		return
	}
	hewanID, err := uuid.Parse(hewanIDStr)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "invalid hewan_id"})
		return
	}

	data, err := c.service.Update(ctx.Request.Context(), pekurbanID, hewanID, req)
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data":    data,
		"message": "porsi updated successfully",
	})
}

// Delete godoc
// @Summary Delete patungan
// @Description Hapus relasi patungan pekurban dengan hewan
// @Tags Patungan
// @Produce json
// @Param pekurban_id path string true "Pekurban ID"
// @Param hewan_id path string true "Hewan ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /patungan/{pekurban_id}/{hewan_id} [delete]
// @Security BearerAuth
func (c *PekurbanHewanController) Delete(ctx *gin.Context) {
	pekurbanID, err := uuid.Parse(ctx.Param("pekurban_id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "invalid pekurban_id"})
		return
	}
	hewanID, err := uuid.Parse(ctx.Param("hewan_id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "invalid hewan_id"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), pekurbanID, hewanID); err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"message": "Joint contribution relation has been deleted successfully",
	})
}