package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kyleochata/runrun/services"
)

type UsersController struct {
	usersService *services.UsersService
}

func NewUsersController(usersService *services.UsersService) *UsersController {
	return &UsersController{usersService: usersService}
}

func (uc UsersController) Login(ctx *gin.Context) {
	username, password, ok := ctx.Request.BasicAuth()
	if !ok {
		log.Println("Error while reading credentials")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	accessToken, err := uc.usersService.Login(username, password)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Status, err)
		return
	}
	ctx.JSON(http.StatusOK, accessToken)
}

func (uc UsersController) Logout(ctx *gin.Context) {
	accessToken := ctx.Request.Header.Get("Token")
	err := uc.usersService.Logout(accessToken)
	if err != nil {
		ctx.AbortWithStatusJSON(err.Status, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
