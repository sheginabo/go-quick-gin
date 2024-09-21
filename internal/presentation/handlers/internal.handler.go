package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type InternalHandler struct{}

func NewInternalHandler() *InternalHandler {
	return &InternalHandler{}
}

func (h *InternalHandler) PostHello(ctx *gin.Context) {
	var req PostHelloRequest
	if err := ValidatePayload(ctx, &req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, ToErrorResponse(*err))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Hello" + req.Message,
	})
}
