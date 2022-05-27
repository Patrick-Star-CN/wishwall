package wishController

import (
	"github.com/gin-gonic/gin"
	"log"
	"wishwall/app/apiExpection"
	"wishwall/app/models"
	"wishwall/app/services/sessionService"
	"wishwall/app/services/wishService"
	"wishwall/app/utils"
)

type WishRes struct {
	Name    string
	Content string
}

type WishReq struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func GetWish(c *gin.Context) {
	wishes, err := wishService.GetWishAll()
	if err != nil {
		log.Println("table wish error:" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	var data []WishRes
	nums := utils.GenerateRandomNumber(0, len(wishes), utils.Min(len(wishes), 9))
	for _, num := range nums {
		data = append(data, WishRes{
			Name:    wishes[num].Name,
			Content: wishes[num].Content,
		})
	}

	utils.JsonSuccessResponse(c, "SUCCESS", gin.H{
		"length": len(data),
		"data":   data,
	})
}

func CreateWish(c *gin.Context) {
	var req WishReq

	errBind := c.ShouldBindJSON(&req)
	if errBind != nil {
		log.Println("request parameter error:" + errBind.Error())
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	user, err := sessionService.GetUserSession(c)
	if err != nil {
		log.Println("session error:" + err.Error())
		_ = c.AbortWithError(200, apiExpection.NotLogin)
		return
	}

	err = wishService.CreateWish(models.Wish{
		Content: req.Content,
		Name:    req.Name,
		UID:     user.ID,
	})
	if err != nil {
		log.Println("table wish error" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, "SUCCESS", nil)
}
