package server

import (
	"database/sql"
	"log"

	"github.com/kyleochata/runrun/controllers"
	"github.com/kyleochata/runrun/repositories"
	"github.com/kyleochata/runrun/services"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpServer struct {
	config            *viper.Viper
	router            *gin.Engine
	runnersController *controllers.RunnersController
	resultsController *controllers.ResultsController
	usersController   *controllers.UsersController
}

func InitHttpServer(config *viper.Viper, dbHandler *sql.DB) HttpServer {
	runnersRepo := repositories.NewRunnersRepository(dbHandler)
	resultRepo := repositories.NewResultsRepository(dbHandler)
	usersRepo := repositories.NewUsersRepository(dbHandler)
	runnersService := services.NewRunnersService(runnersRepo, resultRepo)
	resultsService := services.NewResultsService(resultRepo, runnersRepo)
	usersService := services.NewUsersService(usersRepo)
	runnersController := controllers.NewRunnersController(runnersService, usersService)
	resultsController := controllers.NewResultsController(resultsService, usersService)
	usersController := controllers.NewUsersController(usersService)
	router := setupServer(runnersController, resultsController, usersController)

	return HttpServer{
		config:            config,
		router:            router,
		runnersController: runnersController,
		resultsController: resultsController,
		usersController:   usersController,
	}
}

func (hs HttpServer) Start() {
	err := hs.router.Run(hs.config.GetString("http.server_address"))
	if err != nil {
		log.Fatalf("Error while starting HTTP server: %v", err)
	}
}

func setupServer(runnersController *controllers.RunnersController, resultsController *controllers.ResultsController, usersController *controllers.UsersController) *gin.Engine {
	router := gin.Default()

	router.POST("/runner", runnersController.CreateRunner)
	router.PUT("/runner", runnersController.UpdateRunner)
	router.GET("/runner", runnersController.GetRunnersBatch)

	router.DELETE("/runner/:id", runnersController.DeleteRunner)
	router.GET("/runner/:id", runnersController.GetRunner)

	router.POST("/result", resultsController.CreateResult)
	router.DELETE("/result/:id", resultsController.DeleteResult)

	router.POST("/login", usersController.Login)
	router.POST("/logout", usersController.Logout)
	return router
}
