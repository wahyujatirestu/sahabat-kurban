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

func (c *PekurbanHewanController) Create(ctx *gin.Context) {
	userRaw, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}
	currentUser := userRaw.(model.User)

	var req dto.CreatePekurbanHewanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if currentUser.Role == "user" {
		p, err := c.serv.GetByUserId(ctx.Request.Context(), currentUser.ID)
		if err != nil || p == nil {
			ctx.JSON(403, gin.H{"error": "You have no registered pekurban data"})
			return
		}

		if 	p.ID != req.PekurbanID {
			ctx.JSON(403, gin.H{"error": "You can only register patungan for your own data"})
			return
		}
	}	

	if err := c.service.Create(ctx.Request.Context(), req); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(201, gin.H{"message": "Patungan added successfully"})
}

func (c *PekurbanHewanController) GetAll(ctx *gin.Context) {
	userRaw, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}
	currentUser := userRaw.(model.User)

	if currentUser.Role == "admin" || currentUser.Role == "panitia" {
		list, err := c.service.GetAll(ctx.Request.Context())
		if err != nil {
			ctx.JSON(500, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, gin.H{"data": list})
		return
	}

	// Kalau user biasa, hanya ambil berdasarkan pekurban miliknya
	pekurban, err := c.serv.GetByUserId(ctx.Request.Context(), currentUser.ID)
	if err != nil || pekurban == nil {
		ctx.JSON(403, gin.H{"error": "Pekurban not found for this user"})
		return
	}

	pekurbanID, err := uuid.Parse(pekurban.ID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": "Invalid UUID format for pekurban ID"})
		return
	}

	list, err := c.service.GetByPekurbanId(ctx.Request.Context(), pekurbanID)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"data": list})
}


func (c *PekurbanHewanController) GetByHewanID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("hewan_id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid hewan_id"})
		return
	}

	list, err := c.service.GetByHewanId(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": list})
}

func (c *PekurbanHewanController) GetByPekurbanID(ctx *gin.Context) {
	id, err := uuid.Parse(ctx.Param("pekurban_id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid pekurban_id"})
		return
	}

	list, err := c.service.GetByPekurbanId(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"data": list})
}

func (c *PekurbanHewanController) Update(ctx *gin.Context) {
	var req dto.UpdatePekurbanHewanRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	pekurbanIDStr := ctx.Param("pekurban_id")
	hewanIDStr := ctx.Param("hewan_id")

	pekurbanID, err := uuid.Parse(pekurbanIDStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid pekurban_id"})
		return
	}
	hewanID, err := uuid.Parse(hewanIDStr)
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid hewan_id"})
		return
	}

	if err := c.service.Update(ctx.Request.Context(), pekurbanID, hewanID, req); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "porsi updated successfully"})
}



func (c *PekurbanHewanController) Delete(ctx *gin.Context) {
	pekurbanID, err := uuid.Parse(ctx.Param("pekurban_id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid pekurban_id"})
		return
	}
	hewanID, err := uuid.Parse(ctx.Param("hewan_id"))
	if err != nil {
		ctx.JSON(400, gin.H{"error": "invalid hewan_id"})
		return
	}

	if err := c.service.Delete(ctx.Request.Context(), pekurbanID, hewanID); err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"message": "Joint contribution relation has been deleted successfully"})
}