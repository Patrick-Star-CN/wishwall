package wishController

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"wishwall/app/apiExpection"
	"wishwall/app/models"
	"wishwall/app/services/sessionService"
	"wishwall/app/services/wishService"
	"wishwall/app/utils"
)

type WishRes struct {
	Name     string
	Content  string
	IsClaim  bool
	ClaimUID int
}

type WishReq struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type WishReqUser struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	IsClaim  bool   `json:"isClaim"`
	ClaimUID int    `json:"claimUID"`
}

func GetWish(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
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
			Name:     wishes[num].Name,
			Content:  wishes[num].Content,
			IsClaim:  wishes[num].IsClaim,
			ClaimUID: wishes[num].ClaimUID,
		})
	}

	utils.JsonSuccessResponse(c, "SUCCESS", gin.H{
		"length": len(data),
		"data":   data,
	})
}

func GetWishUser(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	user, errSession := sessionService.GetUserSession(c)
	if errSession != nil {
		log.Println("session error:" + errSession.Error())
		_ = c.AbortWithError(200, apiExpection.NotLogin)
		return
	}

	wishes, err := wishService.GetWishUser(user.ID)
	if err != nil {
		log.Println("table wish error:" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	var data []WishReqUser
	for _, wish := range wishes {
		data = append(data, WishReqUser{
			ID:       wish.ID,
			Name:     wish.Name,
			Content:  wish.Content,
			IsClaim:  wish.IsClaim,
			ClaimUID: wish.ClaimUID,
		})
	}

	utils.JsonSuccessResponse(c, "SUCCESS", gin.H{
		"length": len(data),
		"data":   data,
	})
}

func CreateWish(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
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

func DelWish(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	id_ := c.Query("wishID")
	id, _ := strconv.Atoi(id_)

	user, errSession := sessionService.GetUserSession(c)
	if errSession != nil {
		log.Println("session error:" + errSession.Error())
		_ = c.AbortWithError(200, apiExpection.NotLogin)
		return
	}

	wish, err := wishService.GetWishID(id)
	if err != nil {
		log.Println("table wish error" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	if wish.ID != id {
		utils.JsonSuccessResponse(c, "WISH_ID_ERROR", nil)
		return
	} else if wish.UID != user.ID && !user.IsAdmin {
		utils.JsonSuccessResponse(c, "UID_ERROR", nil)
		return
	}

	err = wishService.DeleteWish(id)
	if err != nil {
		log.Println("table wish error" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, "SUCCESS", nil)
}

func ChangeWish(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	var req WishReqUser

	errBind := c.ShouldBindJSON(&req)
	if errBind != nil {
		log.Println("request parameter error:" + errBind.Error())
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}
	if req.IsClaim {
		utils.JsonSuccessResponse(c, "CLAIMED", nil)
		return
	}

	user, err := sessionService.GetUserSession(c)
	if err != nil {
		log.Println("session error:" + err.Error())
		_ = c.AbortWithError(200, apiExpection.NotLogin)
		return
	}

	wish, err := wishService.GetWishID(req.ID)
	if err != nil {
		log.Println("table wish error" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}
	if wish.ID != req.ID {
		utils.JsonSuccessResponse(c, "WISH_ID_ERROR", nil)
		return
	} else if wish.UID != user.ID && !user.IsAdmin {
		utils.JsonSuccessResponse(c, "UID_ERROR", nil)
		return
	}

	err = wishService.UpdateWish(models.Wish{
		ID:       req.ID,
		UID:      user.ID,
		Content:  req.Content,
		Name:     req.Name,
		IsClaim:  req.IsClaim,
		ClaimUID: req.ClaimUID,
	})
	if err != nil {
		log.Println("table wish error" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, "SUCCESS", nil)
}
