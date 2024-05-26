package httpjson

import (
	"encoding/json"
	"net/http"
)

type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Read[T any](r *http.Request) (T, error) {
	var res T
	err := json.NewDecoder(r.Body).Decode(&res)
	return res, err
}

func Write(w http.ResponseWriter, status int, value any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		json.NewEncoder(w).Encode(BaseResponse{
			Success: false,
			Message: err.Error(),
		})
	}
}

func WriteError(w http.ResponseWriter, status int, err error) {
	Write(w, status, BaseResponse{
		Success: false,
		Message: err.Error(),
	})
}
