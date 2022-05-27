package wishController

import (
	"github.com/gin-gonic/gin"
	"log"
	"wishwall/app/apiExpection"
	"wishwall/app/services/wishService"
	"wishwall/app/utils"
)

type WishRes struct {
	Name    string
	Content string
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

	utils.JsonSuccessResponse(c, "SUCCESS", data)
}
