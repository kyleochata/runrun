package services

import (
	"encoding/base64"
	"net/http"
	"time"

	"github.com/kyleochata/runrun/models"
	"github.com/kyleochata/runrun/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UsersService struct {
	usersRepository *repositories.UsersRepository
}

func NewUsersService(usersRepository *repositories.UsersRepository) *UsersService {
	return &UsersService{
		usersRepository: usersRepository,
	}
}

func (us UsersService) Login(username string, password string) (string, *models.ResponseError) {
	if username == "" || password == "" {
		return "", &models.ResponseError{
			Message: "Invalid username or password",
			Status:  http.StatusBadRequest,
		}
	}
	id, idErr := us.usersRepository.LoginUser(username, password)
	if idErr != nil {
		return "", idErr
	}
	if id == "" {
		return "", &models.ResponseError{
			Message: "Login failed",
			Status:  http.StatusUnauthorized,
		}
	}
	accessToken, err := generateAccessToken(username)
	if err != nil {
		return "", err
	}
	us.usersRepository.SetAccessToken(accessToken, id)
	return accessToken, nil
}

func (us UsersService) Logout(accessToken string) *models.ResponseError {
	if accessToken == "" {
		return &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}
	return us.usersRepository.RemoveAccessToken(accessToken)
}

func (us UsersService) AuthorizeUser(accessToken string, expectedRoles []string) (bool, *models.ResponseError) {
	if accessToken == "" {
		return false, &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}
	role, err := us.usersRepository.GetUserRole(accessToken)
	if err != nil {
		return false, err
	}
	if role == "" {
		return false, &models.ResponseError{
			Message: "Failed to authorize user",
			Status:  http.StatusUnauthorized,
		}
	}
	for _, expectedRole := range expectedRoles {
		if expectedRole == role {
			return true, nil
		}
	}
	return false, nil
}

func (us UsersService) IsAccessTokenValid(accessToken string) (bool, *models.ResponseError) {
	expTime, err := us.usersRepository.GetAccessTokenExpiry(accessToken)
	if err != nil {
		return false, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	if expTime.Before(time.Now()) {
		return false, &models.ResponseError{
			Message: "Access token expired",
			Status:  http.StatusUnauthorized,
		}
	}
	return true, nil
}

func generateAccessToken(username string) (string, *models.ResponseError) {
	hash, err := bcrypt.GenerateFromPassword([]byte(username), bcrypt.DefaultCost)
	if err != nil {
		return "", &models.ResponseError{
			Message: "Failed to generate token",
			Status:  http.StatusInternalServerError,
		}
	}
	return base64.StdEncoding.EncodeToString(hash), nil
}
