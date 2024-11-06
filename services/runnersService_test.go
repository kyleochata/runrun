package services

import (
	"net/http"
	"testing"

	"github.com/kyleochata/runrun/models"
	"github.com/stretchr/testify/assert"
)

func TestValidateRunner(t *testing.T) {
	tests := []struct {
		name   string
		runner *models.Runner
		want   *models.ResponseError
	}{
		{
			name: "Invalid_First_Name",
			runner: &models.Runner{
				LastName: "Smith",
				Age:      30,
				Country:  "Canada",
			},
			want: &models.ResponseError{
				Message: "Invalid first name",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "Invalid_Age",
			runner: &models.Runner{
				FirstName: "John",
				LastName:  "DOE",
				Age:       5000,
				Country:   "Canada",
			},
			want: &models.ResponseError{
				Message: "Invalid age range. Must be between 5 - 125",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "Invalid_Country",
			runner: &models.Runner{
				FirstName: "John",
				LastName:  "Doe",
				Age:       30,
			},
			want: &models.ResponseError{
				Message: "Invalid Country",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "Valid_Runner",
			runner: &models.Runner{
				FirstName: "John",
				LastName:  "Doe",
				Age:       30,
				Country:   "Canada",
			},
			want: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			responseErr := validateRunner(test.runner)
			assert.Equal(t, test.want, responseErr)
		})
	}
}
