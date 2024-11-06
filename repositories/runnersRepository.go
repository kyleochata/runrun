package repositories

import (
	"database/sql"
	"net/http"
	"runners-postgresql/models"

	"github.com/kyleochata/runrun/models"
)

type RunnersRepository struct {
	dbHandler   *sql.DB
	transaction *sql.Tx
}

func NewRunnersRepository(dbHandler *sql.DB) *RunnersRepository {
	return &RunnersRepository{
		dbHandler: dbHandler,
	}
}

// CreateRunner returns a filled runner model and and ID returned from the database layer.
func (rr RunnersRepository) CreateRunner(runner *models.Runner) (*models.Runner, *models.ResponseError) {
	query := `
		INSERT INTO runners(first_name, last_name, age, country)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	// query from db layer returns rows. SELECT cmd and accepts the query and query args.
	rows, err := rr.dbHandler.Query(query, runner.FirstName, runner.LastName, runner.Age, runner.Country)

	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	defer rows.Close()

	var runnerId string

	// Next will turn false and loop will stop when rows are all read
	for rows.Next() {
		// Scan copies cols from the current row into values pointed by the id.
		err := rows.Scan(&runnerId)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	return &models.Runner{
		ID:        runnerId,
		FirstName: runner.FirstName,
		LastName:  runner.LastName,
		Age:       runner.Age,
		IsActive:  true,
		Country:   runner.Country,
	}, nil
}

// UpdateRunner will execute the UPDATE cmd within its query and use Exec. This will execute the query without returning any rows.
func (rr RunnersRepository) UpdateRunner(runner *models.Runner) *models.ResponseError {
	query := `
		UPDATE runners
		SET
			first_name = $1,
			last_name = $2,
			age = $3,
			country = $4,
			WHERE id = $5,
	`
	res, err := rr.dbHandler.Exec(query, runner.FirstName, runner.LastName, runner.Age, runner.Country, runner.ID)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}

	//RowsAffected to check how many rows were affected by the executed query.
	// If zero, then provided runner ID doesn't exist.
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if rowsAffected == 0 {
		return &models.ResponseError{
			Message: "Runner not found",
			Status:  http.StatusNotFound,
		}
	}
	return nil
}

// UpdateRunnerResults executes in transaction because it is connected with the creation and deletion of results.
func (rr RunnersRepository) UpdateRunnerResults(runner *models.Runner) *models.ResponseError {
	query := `
		UPDATE runners
		SET
			personal_best = $1,
			season_best = $2
		WHERE id = $3
	`
	_, err := rr.transaction.Exec(query, runner.PersonalBest, runner.SeasonBest, runner.ID)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}

// DeleteRunner sets the runner's is_active status to false in the database, effectively "deleting" the runner.
// This function returns a *models.ResponseError if an error occurs, or nil if successful.
func (rr RunnersRepository) DeleteRunner(runnerId string) *models.ResponseError {
	query := `
		UPDATE runners
		SET is_active = 'false'
		WHERE id = $1
	`
	// execute the query. is_active updated to false.
	res, err := rr.dbHandler.Exec(query, runnerId)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if rowsAffected == 0 {
		return &models.ResponseError{
			Message: "Runner not found",
			Status:  http.StatusNotFound,
		}
	}
	return nil
}

// GetRunner will return all columns from the runners table.
func (rr RunnersRepository) GetRunner(runnerId string) (*models.Runner, *models.ResponseError) {
	query := `
		SELECT *
		FROM runners
		WHERE id = $1
	`

	rows, err := rr.dbHandler.Query(query, runnerId)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	var id, firstName, lastName, country string
	var personalBest, seasonBest sql.NullString
	var age int
	var isActive bool

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &isActive, &country, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}
	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return &models.Runner{
		ID:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Age:          age,
		IsActive:     isActive,
		Country:      country,
		PersonalBest: personalBest.String,
		SeasonBest:   seasonBest.String,
	}, nil
}

func (rr RunnersRepository) GetAllRunners() ([]*models.Runner, *models.ResponseError) {
	query := `
		SELECT *
		FROM runners
	`
	rows, err := rr.dbHandler.Query(query)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	runners := make([]*models.Runner, 0)
	var id, firstName, lastName, country string
	var personalBest, seasonBest sql.NullString
	var age int
	var isActive bool

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &isActive, &country, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		runner := &models.Runner{
			ID:           id,
			FirstName:    firstName,
			LastName:     lastName,
			Age:          age,
			IsActive:     isActive,
			Country:      country,
			PersonalBest: personalBest.String,
			SeasonBest:   seasonBest.String,
		}
		runners = append(runners, runner)
	}
	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return runners, nil
}

// GetRunnersByCountry takes in a country and returns the fastest 10 active athletes (by personal best)
func (rr RunnersRepository) GetRunnersByCountry(country string) ([]*models.Runner, *models.ResponseError) {
	query := `
		SELECT *
		FROM runners
		WHERE country = $1 AND is_active = 'true'
		ORDER BY personal_best
		LIMIT 10
	`
	rows, err := rr.dbHandler.Query(query, country)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	var id, firstName, lastName string
	var personalBest, seasonBest sql.NullString
	var age int
	var isActive bool

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &isActive, &country, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}
	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return &models.Runner{
		ID:           id,
		FirstName:    firstName,
		LastName:     lastName,
		Age:          age,
		IsActive:     isActive,
		Country:      country,
		PersonalBest: personalBest.String,
		SeasonBest:   seasonBest.String,
	}, nil
}

func (rr RunnersRepository) GetRunnersByYear(year int) ([]*models.Runner, *models.ResponseError) {
	query := `
		SELECT runners.id, runners.first_name, runners.last_name, runners.age, runners.is_active, runners.country, runners.personal_best, results.race_result
		FROM runners
		INNER JOIN (
			SELECT runner_id,
				MIN(race_result) as race_result
			FROM results
			WHERE year = $1
			GROUP BY runner_id) results
		ON runners.id = results.runner_id
		ORDER BY results.race_result
		LIMIT 10
	`
	rows, err := rr.dbHandler.Query(query, year)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	runners := make([]*models.Runner, 0)
	var id, firstName, lastName, country string
	var personalBest, seasonBest sql.NullString
	var age int
	var isActive bool

	for rows.Next() {
		err := rows.Scan(&id, &firstName, &lastName, &age, &isActive, &country, &personalBest, &seasonBest)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
		runner := &models.Runner{
			ID:           id,
			FirstName:    firstName,
			LastName:     lastName,
			Age:          age,
			IsActive:     isActive,
			Country:      country,
			PersonalBest: personalBest.String,
			SeasonBest:   seasonBest.String,
		}
		runners = append(runners, runner)
	}
	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return runners, nil
}
