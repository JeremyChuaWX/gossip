package utils

import (
	"encoding/json"
	"net/http"
)

type BaseResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}

func ReadJSON[T any](r *http.Request) (T, error) {
	var res T
	err := json.NewDecoder(r.Body).Decode(&res)
	return res, err
}

func WriteJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		json.NewEncoder(w).Encode(BaseResponse{
			Error:   true,
			Message: err.Error(),
		})
	}
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, BaseResponse{
		Error:   true,
		Message: err.Error(),
	})
}
