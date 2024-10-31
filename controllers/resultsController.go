package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"runners-postgresql/models"
	"runners-postgresql/services"

	"github.com/gin-gonic/gin"
)

type ResultsController struct {
	resultsService *services.ResultsService
}

func NewResultsController(resultsService *services.ResultsService) *ResultsController {
	return &ResultsController{
		resultsService: resultsService,
	}
}

func (rc ResultsController) CreateResult(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.Println("Error while reading"+" create result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var result models.Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println("Error while unmarshaling "+"creates result request body", err)
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response, resErr := rc.resultsService.CreateResult(&result)
	if resErr != nil {
		ctx.JSON(resErr.Status, resErr)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (rc ResultsController) DeleteResult(ctx *gin.Context) {
	resultId := ctx.Param("id")
	err := rc.resultsService.DeleteResult(resultId)
	if err != nil {
		ctx.JSON(err.Status, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
