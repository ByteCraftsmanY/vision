package models

import "net/http"

type Health struct {
	Status  int    `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
}

func (h Health) Get() *Health {
	return &Health{
		Status:  http.StatusOK,
		Message: "Success",
	}
}
