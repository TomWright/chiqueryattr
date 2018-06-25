package chiqueryattr

import (
	"net/http"
	"strings"
	"fmt"
	"reflect"
	"errors"
	"github.com/go-chi/chi"
)

var (
	// Delimiter specifies the string that will be used to split query parameters when the target is a slice
	Delimiter = ","

	tag             = "queryattr"
	stringType      = reflect.TypeOf("")
	stringSliceType = reflect.TypeOf(make([]string, 0))

	// ErrNonPointerTarget is returned when the given interface does not represent a pointer
	ErrNonPointerTarget = errors.New("invalid Unmarshal target. must be a pointer")
	// ErrInvalidRequest is returned when the given *url.URL is nil
	ErrInvalidRequest = errors.New("invalid request provided")
	// ErrInvalidDelimiter is returned when trying to split a query param into a slice with an invalid separator
	ErrInvalidDelimiter = errors.New("invalid query attr separator")
	// ErrNilSliceField is returned when Unmarshal is given a slice target that has not been initialised
	ErrNilSliceField = errors.New("field target of slice cannot be nil")
)

// Unmarshal attempts to parse query attributes from the specified URL and store any found values
// into the given interface
func Unmarshal(r *http.Request, i interface{}) error {
	if r == nil {
		return ErrInvalidRequest
	}

	iVal := reflect.ValueOf(i)
	if iVal.Kind() != reflect.Ptr || iVal.IsNil() {
		return ErrNonPointerTarget
	}

	v := iVal.Elem()
	t := v.Type()
	var paramVal, tagVal string
	var field reflect.StructField
	var vField reflect.Value

	for i := 0; i < t.NumField(); i++ {
		field = t.Field(i)

		tagVal = field.Tag.Get(tag)
		if tagVal != "" {
			paramVal = chi.URLParam(r, tagVal)

			switch field.Type {
			case stringType:
				v.Field(i).SetString(paramVal)
			case stringSliceType:
				if len(Delimiter) == 0 {
					return ErrInvalidDelimiter
				}
				vField = v.Field(i)
				if vField.IsNil() {
					return ErrNilSliceField
				}
				vField.Set(reflect.AppendSlice(vField, reflect.ValueOf(strings.Split(paramVal, Delimiter))))
			default:
				return fmt.Errorf("invalid field type. `%s` must be `string` or `[]string`", field.Name)
			}
		}
	}
	return nil
}
