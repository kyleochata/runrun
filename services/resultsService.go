package services

import (
	"net/http"
	"runners-postgresql/models"
	"runners-postgresql/repositories"
	"time"
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
	res, err := rs.resultsRepository.CreateResult(result)
	if err != nil {
		return nil, err
	}
	runner, err := rs.runnersRepository.GetRunner(result.RunnerID)
	if err != nil {
		return nil, err
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
			runner.PersonalBest = result.raceResult
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
	err = rs.runnersRepository.UpdateRunnerResults(runner)
	if err != nil {
		return nil, err
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
	result, err := rs.resultsRepository.DeleteResult(resultId)
	if err != nil {
		return err
	}
	runner, err := rs.runnersRepository.GetRunner(result.RunnerId)
	if err != nil {
		return err
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
	if runner.SeasonBest == currYear && result.Year == currYear {
		seasonBest, err := rs.resultsRepository.GetSeasonBestResults(result.RunnerID, result.Year)
		if err != nil {
			return err
		}
		runner.SeasonBest = seasonBest
	}
	err = rs.runnersRepository.UpdateRunnerResults(runner)
	if err != nil {
		return err
	}
	return nil
}
