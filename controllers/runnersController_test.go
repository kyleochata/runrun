package controllers_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/kyleochata/runrun/controllers"
	"github.com/kyleochata/runrun/models"
	"github.com/kyleochata/runrun/repositories"
	"github.com/kyleochata/runrun/services"
	"github.com/stretchr/testify/assert"
)

func initTestRouter(dbHandler *sql.DB) *gin.Engine {
	runnersRepository := repositories.NewRunnersRepository(dbHandler)
	runnersService := services.NewRunnersService(runnersRepository, nil)
	runnersController := controllers.NewRunnersController(runnersService)
	router := gin.Default()
	router.GET("/runner", runnersController.GetRunnersBatch)
	return router
}

func TestGetRunnersResponse(t *testing.T) {
	dbHandler, mock, _ := sqlmock.New()
	defer dbHandler.Close()

	cols := []string{"id", "first_name", "last_name", "age", "is_active", "country", "personal_best", "season_best"}
	mock.ExpectQuery("SELECT *").WillReturnRows(sqlmock.NewRows(cols).AddRow("1", "John", "Doe", 30, true, "Canada", "02:00:41", "02:13:13").AddRow("2", "Maria", "Dove", 30, true, "Serbia", "01:17:28", "01:19:28"))

	router := initTestRouter(dbHandler)
	req, _ := http.NewRequest("GET", "/runner", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)

	var runners []*models.Runner
	json.Unmarshal(recorder.Body.Bytes(), &runners)
	assert.NotEmpty(t, runners)
	assert.Equal(t, 2, len(runners))
}
