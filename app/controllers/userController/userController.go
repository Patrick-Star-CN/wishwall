package userController

import (
	"github.com/gin-gonic/gin"
	"log"
	"wishwall/app/apiExpection"
	"wishwall/app/models"
	"wishwall/app/services/sessionService"
	"wishwall/app/services/userService"
	"wishwall/app/utils"
)

type userForm struct {
	Username string `json:"username"`
	Pwd      string `json:"pwd"`
}

func Login(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	var req userForm

	errBind := c.ShouldBindJSON(&req)
	if errBind != nil {
		log.Println("request parameter error")
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	user, err := userService.GetUser(req.Username)
	if err != nil {
		log.Println("table user error")
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	if user.Name != req.Username {
		utils.JsonSuccessResponse(c, "USERNAME_ERROR", nil)
		return
	} else if user.Pwd != req.Pwd {
		utils.JsonSuccessResponse(c, "PWD_ERROR", nil)
		return
	}

	err = sessionService.SetUserSession(c, user)
	if err != nil {
		log.Println("set session error")
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, "SUCCESS", nil)
}

func Register(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	var req userForm

	errBind := c.ShouldBindJSON(&req)
	if errBind != nil {
		log.Println("request parameter error")
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	user, err := userService.GetUser(req.Username)
	if err != nil {
		log.Println("table user error")
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	if user.Name == req.Username {
		utils.JsonSuccessResponse(c, "USERNAME_ERROR", nil)
		return
	}

	err = userService.CreateUser(models.User{
		Name: req.Username,
		Pwd:  req.Pwd,
	})
	if err != nil {
		log.Println("table user error")
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	err = sessionService.SetUserSession(c, user)
	if err != nil {
		log.Println("set session error")
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, "SUCCESS", nil)
}
