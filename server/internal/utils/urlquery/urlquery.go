package urlquery

import (
	"errors"
	"net/url"
	"reflect"
)

var missingURLQueryKeyError = errors.New("missing URL query key")

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
			return queryStruct, missingURLQueryKeyError
		}
		value := URLQuery.Get(key)
		queryStructValue.FieldByIndex(field.Index).SetString(value)
	}
	return queryStruct, nil
}
