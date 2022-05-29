package claimController

import (
	"github.com/gin-gonic/gin"
	"log"
	"wishwall/app/apiExpection"
	"wishwall/app/controllers/wishController"
	"wishwall/app/models"
	"wishwall/app/services/sessionService"
	"wishwall/app/services/wishService"
	"wishwall/app/utils"
)

func SubmitClaim(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	var req wishController.Req

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
	if !user.IsAdmin {
		utils.JsonSuccessResponse(c, "IS_NOT_ADMIN", nil)
		return
	}

	var wish *models.Wish
	wish, err = wishService.GetWishID(req.ID)
	if wish.IsClaim {
		utils.JsonSuccessResponse(c, "CLAIMED", nil)
		return
	} else if wish.ID != req.ID {
		utils.JsonSuccessResponse(c, "ID_ERROR", nil)
		return
	}
	err = wishService.UpdateWish(models.Wish{
		Name:     wish.Name,
		Content:  wish.Content,
		ID:       req.ID,
		UID:      wish.UID,
		IsClaim:  true,
		ClaimUID: user.ID,
	})
	if err != nil {
		log.Println("table wish error:" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, "SUCCESS", nil)
}

func CancelClaim(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	var req wishController.Req

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
	if !user.IsAdmin {
		utils.JsonSuccessResponse(c, "IS_NOT_ADMIN", nil)
		return
	}

	var wish *models.Wish
	wish, err = wishService.GetWishID(req.ID)
	if wish.ID != req.ID || !wish.IsClaim {
		utils.JsonSuccessResponse(c, "ID_ERROR", nil)
		return
	} else if wish.ClaimUID != user.ID {
		utils.JsonSuccessResponse(c, "UID_ERROR", nil)
		return
	}
	err = wishService.UpdateWish(models.Wish{
		Name:     wish.Name,
		Content:  wish.Content,
		ID:       wish.ID,
		UID:      wish.UID,
		IsClaim:  false,
		ClaimUID: 0,
	})
	if err != nil {
		log.Println("table wish error:" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, "SUCCESS", nil)
}
