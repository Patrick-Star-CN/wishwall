package router

import (
	"github.com/gin-gonic/gin"
	"wishwall/app/controllers/userController"
)

func userRouter(r *gin.RouterGroup) {
	r.POST("/login", userController.Login)
	r.POST("/register", userController.Register)
	r.GET("/loginOut", userController.LoginOut)
	r.GET("/checkLogin", userController.CheckLogin)
}
