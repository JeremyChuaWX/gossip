package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"reflect"
)

var MissingURLQueryKeyError = errors.New("missing URL query key")

type errorResponse struct {
	Error string `json:"error"`
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
		json.NewEncoder(w).Encode(errorResponse{
			Error: err.Error(),
		})
	}
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, errorResponse{
		Error: err.Error(),
	})
}

func GetURLQueryStruct[T any](URL *url.URL) (T, error) {
	var queryStruct T
	URLQuery, err := url.ParseQuery(URL.RawQuery)
	if err != nil {
		return queryStruct, err
	}
	queryStructType := reflect.TypeOf(queryStruct)
	queryStructValue := reflect.ValueOf(&queryStruct).Elem()
	for _, field := range reflect.VisibleFields(queryStructType) {
		key := field.Tag.Get("query")
		if !URLQuery.Has(key) {
			return queryStruct, MissingURLQueryKeyError
		}
		value := URLQuery.Get(key)
		queryStructValue.FieldByIndex(field.Index).SetString(value)
	}
	return queryStruct, nil
}
