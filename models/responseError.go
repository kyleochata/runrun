package models

type ResponseError struct {
	Message string `json:"message"`
	Status  int    `json:"-"` //tells json parser to completely ignore field
}
