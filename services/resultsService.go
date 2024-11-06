package services

import (
	"net/http"

	"time"

	"github.com/kyleochata/runrun/models"
	"github.com/kyleochata/runrun/repositories"
)

type ResultsService struct {
	resultsRepository *repositories.ResultsRepository
	runnersRepository *repositories.RunnersRepository
}

func NewResultsService(resultsRepository *repositories.ResultsRepository, runnersRepository *repositories.RunnersRepository) *ResultsService {
	return &ResultsService{
		resultsRepository: resultsRepository,
		runnersRepository: runnersRepository,
	}
}

func (rs ResultsService) CreateResult(result *models.Result) (*models.Result, *models.ResponseError) {
	if result.RunnerID == "" {
		return nil, &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}
	if result.RaceResult == "" {
		return nil, &models.ResponseError{
			Message: "Invalid race results",
			Status:  http.StatusBadRequest,
		}
	}
	if result.Location == "" {
		return nil, &models.ResponseError{
			Message: "Invalid Location of result",
			Status:  http.StatusBadRequest,
		}
	}
	if result.Position < 0 {
		return nil, &models.ResponseError{
			Message: "Invalid position in result.",
			Status:  http.StatusBadRequest,
		}
	}
	currYear := time.Now().Year()
	if result.Year < 0 || result.Year > currYear {
		return nil, &models.ResponseError{
			Message: "Invalid year for race result.",
			Status:  http.StatusBadRequest,
		}
	}
	raceResult, err := parseRaceResult(result.RaceResult)
	if err != nil {
		return nil, &models.ResponseError{
			Message: "Invalid race result. Unable to parse.",
			Status:  http.StatusBadRequest,
		}
	}
	res, resErr := rs.resultsRepository.CreateResult(result)
	if resErr != nil {
		return nil, resErr
	}
	runner, runErr := rs.runnersRepository.GetRunner(result.RunnerID)
	if runErr != nil {
		return nil, runErr
	}
	if runner == nil {
		return nil, &models.ResponseError{
			Message: "Invalid runner not found.",
			Status:  http.StatusNotFound,
		}
	}
	if runner.PersonalBest == "" {
		runner.PersonalBest = result.RaceResult
	} else {
		personalBest, err := parseRaceResult(runner.PersonalBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Failed to parse " + "personal best",
				Status:  http.StatusInternalServerError,
			}
		}
		if raceResult < personalBest {
			runner.PersonalBest = result.RaceResult
		}
	}

	if result.Year == currYear {
		if runner.SeasonBest == "" {
			runner.SeasonBest = result.RaceResult
		} else {
			seasonBest, err := parseRaceResult(runner.SeasonBest)
			if err != nil {
				return nil, &models.ResponseError{
					Message: "Failed to parse " + "season best",
					Status:  http.StatusBadRequest,
				}
			}
			if raceResult < seasonBest {
				runner.SeasonBest = result.RaceResult
			}
		}
	}
	updateRunnerErr := rs.runnersRepository.UpdateRunnerResults(runner)
	if updateRunnerErr != nil {
		return nil, updateRunnerErr
	}
	return res, nil
}

func (rs ResultsService) DeleteResult(resultId string) *models.ResponseError {
	if resultId == "" {
		return &models.ResponseError{
			Message: "Invalid race result. Unable to parse.",
			Status:  http.StatusBadRequest,
		}
	}
	err := repositories.BeginTransaction(rs.runnersRepository, rs.resultsRepository)
	if err != nil {
		return &models.ResponseError{
			Message: "Failed to start transaction",
			Status:  http.StatusBadRequest,
		}
	}
	result, resErr := rs.resultsRepository.DeleteResult(resultId)
	if resErr != nil {
		return resErr
	}
	runner, runErr := rs.runnersRepository.GetRunner(result.RunnerID)
	if runErr != nil {
		return runErr
	}
	if runner.PersonalBest == result.RaceResult {
		personalBest, err := rs.resultsRepository.GetPersonalBestResults(result.RunnerID)
		if err != nil {
			return err
		}
		runner.PersonalBest = personalBest
	}

	//Check if the deleted result is season best for runner
	currYear := time.Now().Year()
	if runner.SeasonBest == result.RaceResult && result.Year == currYear {
		seasonBest, err := rs.resultsRepository.GetSeasonBestResults(result.RunnerID, result.Year)
		if err != nil {
			return err
		}
		runner.SeasonBest = seasonBest
	}
	resErr = rs.runnersRepository.UpdateRunnerResults(runner)
	if resErr != nil {
		repositories.RollbackTransaction(rs.runnersRepository, rs.resultsRepository)
		return resErr
	}
	repositories.CommitTransaction(rs.runnersRepository, rs.resultsRepository)
	return nil
}

func parseRaceResult(strTime string) (time.Duration, error) {
	return time.ParseDuration(strTime[0:2] + "h" + strTime[3:5] + "m" + strTime[6:8] + "s")
}
