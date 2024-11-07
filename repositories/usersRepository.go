package repositories

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/kyleochata/runrun/models"
)

type UsersRepository struct {
	dbHandler *sql.DB
}

func NewUsersRepository(dbHandler *sql.DB) *UsersRepository {
	return &UsersRepository{
		dbHandler: dbHandler,
	}
}

func (ur UsersRepository) LoginUser(username, password string) (string, *models.ResponseError) {
	query := `
		SELECT id
		FROM users
		WHERE username = $1 AND
		user_password = crypt($2, user_password)`

	rows, err := ur.dbHandler.Query(query, username, password)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	var id string
	for rows.Next() {
		err := rows.Scan(&id)
		if err != nil {
			return "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}

	if err := rows.Err(); err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return id, nil
}

func (ur UsersRepository) GetUserRole(accessToken string) (string, *models.ResponseError) {
	query := `
	SELECT user_role
	FROM users
	WHERE access_token = $1 AND
	access_token_expiry > NOW()`

	rows, err := ur.dbHandler.Query(query, accessToken)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()

	var role string
	for rows.Next() {
		err := rows.Scan(&role)
		if err != nil {
			return "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}
	if err := rows.Err(); err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return role, nil
}

func (ur UsersRepository) SetAccessToken(accessToken, id string) *models.ResponseError {
	expiry := time.Now().Add(15 * time.Minute)
	query := `
	UPDATE users
	SET access_token = $1, access_token_expiry = $2
	WHERE id = $3`
	_, err := ur.dbHandler.Exec(query, accessToken, expiry, id)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}

func (ur UsersRepository) RemoveAccessToken(accessToken string) *models.ResponseError {
	query := `
	UPDATE users
	SET access_token = '', access_token_expiry = NULL
	WHERE access_token = $1`
	_, err := ur.dbHandler.Exec(query, accessToken)
	if err != nil {
		return &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return nil
}

func (ur UsersRepository) GetAccessTokenExpiry(accessToken string) (time.Time, error) {
	query := `
	SELECT access_token_expiry
	FROM users
	WHERE access_token = $1`

	var expiry time.Time
	err := ur.dbHandler.QueryRow(query, accessToken).Scan(&expiry)
	if err != nil {
		return time.Time{}, err
	}
	return expiry, nil
}
