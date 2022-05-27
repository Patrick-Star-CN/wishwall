package router

import (
	"github.com/gin-gonic/gin"
	"wishwall/app/controllers/wishController"
)

func wishRouter(r *gin.RouterGroup) {
	r.POST("/submitWish", wishController.CreateWish)
	r.GET("/getWish", wishController.GetWish)
}
