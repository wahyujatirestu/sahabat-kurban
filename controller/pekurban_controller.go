package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/dto"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/service"
)

type PekurbanController struct {
	pService service.PekurbanService
}

func NewPekurbanController(pService service.PekurbanService) *PekurbanController {
	return &PekurbanController{pService: pService}
}

// Create godoc
// @Summary Create Pekurban
// @Description Daftarkan user atau panitia sebagai pekurban
// @Tags Pekurban
// @Accept json
// @Produce json
// @Param request body dto.CreatePekurbanRequest true "Pekurban request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekurban [post]
func (c *PekurbanController) Create(ctx *gin.Context) {
	userRaw, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(401, gin.H{
			"status": 401,
			"error": "Unauthorized"})
		return
	}

	currentUser := userRaw.(model.User)

	var req dto.CreatePekurbanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	if currentUser.Role == "user" {
		if req.UserID == nil {
			ctx.JSON(403, gin.H{
				"status": 403,
				"error": "You must include your own user_id"})
			return
		}
		uid, err := uuid.Parse(*req.UserID)
		if err != nil || uid != currentUser.ID {
			ctx.JSON(403, gin.H{
				"status": 403,
				"error": "You can only register yourself as pekurban"})
			return
		}
	}


	if currentUser.Role == "user" {
		existing, _ := c.pService.GetByUserId(ctx.Request.Context(), currentUser.ID)
		if existing != nil {
			ctx.JSON(409, gin.H{
				"status": 409,
				"error": "You have already registered as pekurban"})
			return
		}
	}

	res, err := c.pService.Create(ctx.Request.Context(), req)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{
		"status": 201,
		"data": res,
		"message": "Pekurban created successfully",
	})
}

// GetAll godoc
// @Summary Get all Pekurban
// @Description Ambil semua data pekurban (khusus admin & panitia)
// @Tags Pekurban
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekurban [get]
func (c *PekurbanController) GetAll(ctx *gin.Context) {
	data, err := c.pService.GetAll(ctx.Request.Context())
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": data,
		"message": "Pekurban retrieved successfully",
	})
}

// GetById godoc
// @Summary Get Pekurban by ID
// @Description Ambil detail pekurban berdasarkan ID (khusus admin & panitia)
// @Tags Pekurban
// @Produce json
// @Param id path string true "Pekurban ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekurban/{id} [get]
func (c *PekurbanController) GetById(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid ID"})
		return
	}

	data, err := c.pService.GetById(ctx.Request.Context(), id)
	if err != nil || data == nil {
		ctx.JSON(404, gin.H{
			"status": 404,
			"error": "Pekurban not found"})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": data,
		"message": "Pekurban retrieved successfully",
	})
}

// GetMe godoc
// @Summary Get My Pekurban Data
// @Description Ambil data pekurban milik user yang sedang login
// @Tags Pekurban
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekurban/me [get]
func (c *PekurbanController) GetMe(ctx *gin.Context) {
	userRaw, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(401, gin.H{
			"status": 401,
			"error": "Unauthorized"})
		return
	}

	currentUser := userRaw.(model.User)

	result, err := c.pService.GetByUserId(ctx.Request.Context(), currentUser.ID)
	if err != nil || result == nil {
		ctx.JSON(404, gin.H{
			"status": 404,
			"error": "Data pekurban tidak ditemukan"})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": result,
		"message": "Pekurban retrieved successfully",
	})
}

// Update godoc
// @Summary Update Pekurban
// @Description Ubah data pekurban berdasarkan ID (role-based access)
// @Tags Pekurban
// @Accept json
// @Produce json
// @Param id path string true "Pekurban ID"
// @Param request body dto.UpdatePekurbanRequest true "Update request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekurban/{id} [put]
func (c *PekurbanController) Update(ctx *gin.Context) {
	userRaw, exist := ctx.Get("user")
	if !exist {
		ctx.JSON(401, gin.H{
			"status": 401,
			"error": "Unauthorized"})
		return
	}

	currentUser := userRaw.(model.User)

	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid ID"})
	}

	p, err := c.pService.GetById(ctx.Request.Context(), id)
	if err != nil || p == nil {
		ctx.JSON(404, gin.H{
			"status": 404,
			"error": "Pekurban not found"})
		return
	}

	if currentUser.Role != "admin" {
		if currentUser.Role == "user" {
			if p.UserID == nil || *p.UserID != currentUser.ID.String() {
				ctx.JSON(403, gin.H{
					"status": 403,
					"error": "You can only update your own kurban data"})
				return
			}
		}
		if currentUser.Role == "panitia" {
			if !(p.UserID == nil || *p.UserID == currentUser.ID.String()) {
				ctx.JSON(403, gin.H{
					"status": 403,
					"error": "Panitia can only update offline data or their own"})
			}
		}
	}

	var req dto.UpdatePekurbanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": err.Error()})
		return
	}

	data, err := c.pService.Update(ctx.Request.Context(), id, req)
	if err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"data": data,
		"message": "Pekurban updated successfully",
	})
}

// Delete godoc
// @Summary Delete Pekurban
// @Description Hapus data pekurban berdasarkan ID (khusus admin)
// @Tags Pekurban
// @Produce json
// @Param id path string true "Pekurban ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Security BearerAuth
// @Router /pekurban/{id} [delete]
func (c *PekurbanController) Delete(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"status": 400,
			"error": "Invalid ID"})
		return
	}

	if err := c.pService.Delete(ctx.Request.Context(), id); err != nil {
		ctx.JSON(500, gin.H{
			"status": 500,
			"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"status": 200,
		"message": "Pekurban deleted successfully"})
}