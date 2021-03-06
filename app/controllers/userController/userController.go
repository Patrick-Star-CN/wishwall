package userController

import (
	"github.com/gin-gonic/gin"
	"log"
	"wishwall/app/apiExpection"
	"wishwall/app/models"
	"wishwall/app/services/sessionService"
	"wishwall/app/services/stuService"
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
		log.Println("request parameter error:" + errBind.Error())
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	user, err := userService.GetUser(req.Username)
	if err != nil {
		log.Println("table user error:" + err.Error())
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
		log.Println("set session error:" + err.Error())
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
		log.Println("request parameter error:" + errBind.Error())
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	user, err := userService.GetUser(req.Username)
	if err != nil {
		log.Println("table user error:" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	if user.Name == req.Username {
		utils.JsonSuccessResponse(c, "USERNAME_REGISTERED", nil)
		return
	}

	var stu *models.Stu
	stu, err = stuService.GetStu(req.Username)
	if err != nil {
		log.Println("table stu error:" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	if stu.Num != req.Username {
		utils.JsonSuccessResponse(c, "USERNAME_ERROR", nil)
		return
	}

	err = userService.CreateUser(models.User{
		Name:    req.Username,
		Pwd:     req.Pwd,
		IsAdmin: false,
	})
	if err != nil {
		log.Println("table user error:" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	user, err = userService.GetUser(req.Username)
	if err != nil {
		log.Println("table user error:" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	err = sessionService.SetUserSession(c, user)
	if err != nil {
		log.Println("set session error:" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	utils.JsonSuccessResponse(c, "SUCCESS", nil)
}

func LoginOut(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	_, errSession := sessionService.GetUserSession(c)
	if errSession != nil {
		_ = c.AbortWithError(200, apiExpection.NotLogin)
		return
	}

	sessionService.ClearUserSession(c)
	utils.JsonSuccessResponse(c, "SUCCESS", nil)
}

func CheckLogin(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	user, errSession := sessionService.GetUserSession(c)
	if errSession != nil {
		_ = c.AbortWithError(200, apiExpection.NotLogin)
		return
	}

	utils.JsonSuccessResponse(c, "SUCCESS", user.Name)
}
