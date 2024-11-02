package services

import (
	"net/http"
	"runners-postgresql/models"
	"runners-postgresql/repositories"
	"strconv"
	"time"
)

type RunnersService struct {
	runnersRepository *repositories.RunnersRepository
	resultsRepository *repositories.ResultsRepository
}

func NewRunnersService(runnersRepository *repositories.RunnersRepository, resultsRepository *repositories.ResultsRepository) *RunnersService {
	return &RunnersService{
		runnersRepository: runnersRepository,
		resultsRepository: resultsRepository,
	}
}

func (rs RunnersService) CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError) {
	err := validateRunner(runner)
	if err != nil {
		return nil, err
	}
	return rs.runnersRepository.CreateRunner(runner)
}

func validateRunner(runner *models.Runner) *models.ResponseError {
	if runner.FirstName == "" {
		return &models.ResponseError{
			Message: "Invalid first name",
			Status:  http.StatusBadRequest,
		}
	}
	if runner.LastName == "" {
		return &models.ResponseError{
			Message: "Invalid last name",
			Status:  http.StatusBadRequest,
		}
	}
	if runner.Age <= 5 || runner.Age > 125 {
		return &models.ResponseError{
			Message: "Invalid age range. Must be between 5 - 125",
			Status:  http.StatusBadRequest,
		}
	}
	if runner.Country == "" {
		return &models.ResponseError{
			Message: "Invalid Country",
			Status:  http.StatusBadRequest,
		}
	}
	return nil
}

func (rs RunnersService) UpdateRunner(runner *models.Runner) *models.ResponseError {
	err := validateRunnerId(runner.ID)
	if err != nil {
		return err
	}
	err = validateRunner(runner)
	if err != nil {
		return err
	}
	return rs.runnersRepository.UpdateRunner(runner)
}

func validateRunnerId(runnerId string) *models.ResponseError {
	if runnerId == "" {
		return &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}
	return nil
}

func (rs RunnersService) DeleteRunner(runnerId string) *models.ResponseError {
	err := validateRunnerId(runnerId)
	if err != nil {
		return err
	}
	return rs.runnersRepository.DeleteRunner(runnerId)
}

func (rs RunnersService) GetRunner(runnerId string) (*models.Runner, *models.ResponseError) {
	err := validateRunnerId(runnerId)
	if err != nil {
		return nil, err
	}

	runner, err := rs.runnersRepository.GetRunner(runnerId)
	if err != nil {
		return nil, err
	}

	results, err := rs.resultsRepository.GetRunner(runnerId)
	if err != nil {
		return nil, err
	}

	runner.Results = results
	return runner, nil
}

func (rs RunnersService) GetRunnersBatch(country, year string) ([]*models.Runner, *models.ResponseError) {
	if country != "" && year != "" {
		return nil, &models.ResponseError{
			Message: "Only one parameter can be passed. Country or Year.",
			Status:  http.StatusBadRequest,
		}
	}
	if country != "" {
		return rs.runnersRepository.GetRunnersByCountry(country)
	}
	if year != "" {
		intYear, err := strconv.Atoi(year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Invalid year",
				Status:  http.StatusBadRequest,
			}
		}
		currentYear := time.Now().Year()
		if intYear < 0 || intYear > currentYear {
			return nil, &models.ResponseError{
				Message: "Invalid year.",
				Status:  http.StatusBadRequest,
			}
		}
		return rs.runnersRepository.GetRunnersByYear(intYear)
	}
	return rs.runnersRepository.GetAllRunners()
}
