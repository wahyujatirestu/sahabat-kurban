package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/wahyujatirestu/sahabat-kurban/model"
	"github.com/wahyujatirestu/sahabat-kurban/utils/service"
)

type AuthMiddleware interface {
	RequireToken(roles ...string) gin.HandlerFunc
}

type auhtMiddleware struct {
	jwtService service.JWTService
}

type authHeader struct {
	Authorization string `header:"Authorization" binding:"required"`
}

func NewAuthMiddleware(jwtService service.JWTService) AuthMiddleware {
	return &auhtMiddleware{jwtService: jwtService}
}

func (a *auhtMiddleware) RequireToken(roles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var authHeader authHeader
		if err := ctx.ShouldBindHeader(&authHeader); err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return 
		}

		token :=  strings.Replace(authHeader.Authorization, "Bearer ", "", 1)
		if token == "" {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return 
		}

		tokenClaim, err := a.jwtService.ValidateAccessToken(token)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return 
		}

		userId, err := uuid.Parse(tokenClaim.UserId)
		if err != nil {
			ctx.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
			return 
		}

		ctx.Set("user", model.User{
			ID: userId,
			Role: tokenClaim.Role,
		})

		validRole := false
		if len(roles) == 0 {
			validRole = true
		} else {
			for _, role := range roles{
				if role == tokenClaim.Role {
					validRole = true
					break
				}
			}
		}

		if !validRole {
			ctx.AbortWithStatusJSON(403, gin.H{"error": "Forbidden: insufficient permissions"})
			return 
		}

		ctx.Next()
	}
}