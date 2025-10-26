package scripts

import (
	"fmt"
	"reflect"
	"testing"
	"encoding/json"
	"strings"
)

func WrapError(context string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", context, err)
}

func PtrToElem(v reflect.Value) reflect.Value {
	if v.Kind() != reflect.Ptr {
		fmt.Printf("Value %T is already an element\n", v)
		return v
	}
	return v.Elem()
}

func GetInfoStr(s any) (string, error) {
	t := reflect.TypeOf(s)
	errorMessage := fmt.Sprintf("Failed to return %s data as string", t.Name())

	sJSON, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", WrapError(errorMessage, err)
	}

	return string(sJSON), nil
}

func GetInfoMap(s any) (map[string]interface{}) {
	sData := make(map[string]interface{})
	v := reflect.ValueOf(s).Elem()

	vElem := PtrToElem(v)
	tElem := vElem.Type()

	for i := 0; i < vElem.NumField(); i++ {
		fieldName := tElem.Field(i).Name
		fieldValue := vElem.Field(i).Interface()

		sData[strings.ToLower(fieldName)] = fieldValue
	}
	return sData
}

type TestCase[T any, U any] struct { // T = expected type; U = got type
	Name string
	Input T
	Want U
	WantErr bool
}

func BoilerTestFunc[T any, U any](
	t *testing.T,
	fn func(T) (U, error),
	testCases []TestCase[T, U],
) {
	for _, test := range testCases {
		t.Run(test.Name, func(t *testing.T) {
			got, err := fn(test.Input)

			if test.WantErr {
				if err == nil {
					t.Errorf("Expected error for input \"%v\", got nil", test.Input)
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error for input \"%v\": %v", test.Input, err)
				return
			}

			if !reflect.DeepEqual(got, test.Want) {
				t.Errorf("For input \"%v\" expected %v, got %v", test.Input, test.Want, got)
			}
		})
	}
}