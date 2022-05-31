package wishController

import (
	"github.com/gin-gonic/gin"
	"log"
	"wishwall/app/apiExpection"
	"wishwall/app/models"
	"wishwall/app/services/sessionService"
	"wishwall/app/services/userService"
	"wishwall/app/services/wishService"
	"wishwall/app/utils"
)

type WishReq struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type Req struct {
	ID int `json:"id"`
}

type WishReqUser struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Content   string `json:"content"`
	IsClaim   bool   `json:"isClaim"`
	ClaimName string `json:"claimName"`
}

func GetWish(c *gin.Context) {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
	wishes, err := wishService.GetWishAll()
	if err != nil {
		log.Println("table wish error:" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	var data []WishReqUser
	nums := utils.GenerateRandomNumber(0, len(wishes), utils.Min(len(wishes), 9))
	for _, num := range nums {
		data = append(data, WishReqUser{
			Name:      wishes[num].Name,
			Content:   wishes[num].Content,
			IsClaim:   wishes[num].IsClaim,
			ClaimName: wishes[num].ClaimName,
			ID:        wishes[num].ID,
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
		_ = c.AbortWithError(200, apiExpection.NotLogin)
		return
	}

	var wishes []models.Wish
	var err error

	if user.IsAdmin {
		wishes, err = wishService.GetWishAll()
	} else {
		wishes, err = wishService.GetWishUser(user.ID)
	}
	if err != nil {
		log.Println("table wish error:" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	var data []WishReqUser
	for _, wish := range wishes {
		data = append(data, WishReqUser{
			ID:        wish.ID,
			Name:      wish.Name,
			Content:   wish.Content,
			IsClaim:   wish.IsClaim,
			ClaimName: wish.ClaimName,
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
		_ = c.AbortWithError(200, apiExpection.NotLogin)
		return
	}

	err = wishService.CreateWish(models.Wish{
		Content:   req.Content,
		Name:      req.Name,
		UID:       user.ID,
		IsClaim:   false,
		ClaimName: "",
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
	var req Req

	errBind := c.ShouldBindJSON(&req)
	if errBind != nil {
		log.Println("request parameter error:" + errBind.Error())
		_ = c.AbortWithError(200, apiExpection.ParamError)
		return
	}

	user, errSession := sessionService.GetUserSession(c)
	if errSession != nil {
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

	err = wishService.DeleteWish(req.ID)
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

	user, err = userService.GetUser(req.Name)
	if err != nil {
		log.Println("table user error" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	err = wishService.UpdateWish(models.Wish{
		ID:        req.ID,
		UID:       user.ID,
		Content:   req.Content,
		Name:      req.Name,
		IsClaim:   req.IsClaim,
		ClaimName: req.ClaimName,
	})
	if err != nil {
		log.Println("table wish error" + err.Error())
		_ = c.AbortWithError(200, apiExpection.ServerError)
		return
	}

	utils.JsonSuccessResponse(c, "SUCCESS", nil)
}
