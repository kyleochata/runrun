package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/kyleochata/runrun/models"
	"github.com/kyleochata/runrun/services"

	"github.com/gin-gonic/gin"
)

const (
	ROLE_ADMIN  string = "admin"
	ROLE_RUNNER string = "runner"
)

type RunnersController struct {
	runnersService *services.RunnersService
	usersService   *services.UsersService
}

func NewRunnersController(runnersService *services.RunnersService, usersService *services.UsersService) *RunnersController {
	return &RunnersController{
		runnersService: runnersService,
		usersService:   usersService,
	}
}

func (rc RunnersController) CreateRunner(ctx *gin.Context) {

	accessToken := ctx.Request.Header.Get("Token")
	auth, authErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	if authErr != nil {
		ctx.JSON(authErr.Status, authErr)
		return
	}
	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading create runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("Error during unmarshaling"+" update runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	responseErr := rc.runnersService.UpdateRunner(&runner)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (rc RunnersController) UpdateRunner(ctx *gin.Context) {

	accessToken := ctx.Request.Header.Get("Token")
	auth, authErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	if authErr != nil {
		ctx.JSON(authErr.Status, authErr)
		return
	}
	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading: update runner request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println("Error while unmarshaling: update runner req body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	resErr := rc.runnersService.UpdateRunner(&runner)
	if resErr != nil {
		ctx.AbortWithStatusJSON(resErr.Status, resErr)
		return
	}
	ctx.Status(http.StatusNoContent)
}

func (rc RunnersController) DeleteRunner(ctx *gin.Context) {
	accessToken := ctx.Request.Header.Get("Token")
	auth, authErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN})
	if authErr != nil {
		ctx.JSON(authErr.Status, authErr)
		return
	}
	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	runnerId := ctx.Param("id")
	responseErr := rc.runnersService.DeleteRunner(runnerId)
	if responseErr != nil {
		ctx.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
}

func (rc RunnersController) GetRunner(ctx *gin.Context) {

	accessToken := ctx.Request.Header.Get("Token")
	auth, authErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN, ROLE_RUNNER})
	if authErr != nil {
		ctx.JSON(authErr.Status, authErr)
		return
	}
	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	runnerId := ctx.Param("id")
	response, err := rc.runnersService.GetRunner(runnerId)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (rc RunnersController) GetRunnersBatch(ctx *gin.Context) {
	accessToken := ctx.Request.Header.Get("Token")
	auth, authErr := rc.usersService.AuthorizeUser(accessToken, []string{ROLE_ADMIN, ROLE_RUNNER})
	if authErr != nil {
		ctx.JSON(authErr.Status, authErr)
		return
	}
	if !auth {
		ctx.Status(http.StatusUnauthorized)
		return
	}

	params := ctx.Request.URL.Query()
	country := params.Get("country")
	year := params.Get("year")
	response, err := rc.runnersService.GetRunnersBatch(country, year)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
