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

	ctx.JSON(201, gin.H{"message": "Patungan berhasil ditambahkan"})
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






