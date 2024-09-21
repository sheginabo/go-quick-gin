package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "OK"})
}

func ToErrorResponse(err ResponseError) map[string]any {
	return gin.H{"error": err}
}
