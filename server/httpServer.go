package server

import (
	"database/sql"
	"log"
	"runners-postgresql/controllers"
	"runners-postgresql/repositories"
	"runners-postgresql/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpServer struct {
	config            *viper.Viper
	router            *gin.Engine
	runnersController *controllers.RunnersController
	resultsController *controllers.ResultsController
}

func InitHttpServer(config *viper.Viper, dbHandler *sql.DB) HttpServer {
	runnersRepo := repositories.NewRunnersRepository(dbHandler)
	resultRepo := repositories.NewResultsRepository(dbHandler)
	runnersService := services.NewRunnersService(runnersRepo, resultRepo)
	resultsService := services.NewResultsService(resultRepo, runnersRepo)
	runnersController := controllers.NewRunnersController(runnersService)
	resultsController := controllers.NewResultsController(resultsService)
	router := setupServer(runnersController, resultsController)

	return HttpServer{
		config:            config,
		router:            router,
		runnersController: runnersController,
		resultsController: resultsController,
	}
}

func (hs HttpServer) Start() {
	err := hs.router.Run(hs.config.GetString("http.server_address"))
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}

func setupServer(runnersController *controllers.RunnersController, resultsController *controllers.ResultsController) *gin.Engine {
	router := gin.Default()

	router.POST("/runner", runnersController.CreateRunner)
	router.PUT("/runner", runnersController.UpdateRunner)
	router.GET("/runner", runnersController.GetRunnersBatch)

	router.DELETE("/runner/:id", runnersController.DeleteRunner)
	router.GET("/runner/:id", runnersController.GetRunner)

	router.POST("/result", resultsController.CreateResult)
	router.DELETE("/result/:id", resultsController.DeleteResult)

	return router
}
