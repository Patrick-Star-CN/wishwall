package router

import (
	"github.com/gin-gonic/gin"
	"wishwall/app/controllers/claimController"
)

func claimRouter(r *gin.RouterGroup) {
	r.POST("/submitClaim", claimController.SubmitClaim)
	r.DELETE("/deleteClaim", claimController.CancelClaim)
}
