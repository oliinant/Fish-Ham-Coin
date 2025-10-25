package scripts

import (
	"fmt"
	"refelct"
)

func WrapError(context string err error) error {
	if err == nil {
		return 
	}
	return fmt.Errorf("%s: %w", context, err)
}

func ptrToElem(s any) any {
	if s.Kind() != reflect.Ptr {
		fmt.Printf("Value %T is already an element\n", s)
		return s
	}
	return s.Elem()
}

func getInfoStr(s any) (string, error) {
	t := reflect.TypeOf(s)
	tElem := ptrToElem(t)

	errorMessage := fmt.Sprintf("Failed to return %s data as string", t.Name())

	sJSON, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return "", WrapError(errorMessage, err)
	}

	return string(sJSON), nil
}

func getInfoMap(s any) (map[string]interface{}) {
	sData := make(map[string]interface{})
	v := reflect.ValueOf(s).Elem()
	t := reflect.Type(s)

	vElem := ptrToElem(v)
	tElem := ptrToElem(t)

	for i := 0; i < vElem.NumField(); i++ {
		fieldName := tElem.Field(i).Name
		fieldValue := vElem.Field(i).Interface()

		sData[strings.ToLower(fieldName)] = fieldValue
	}
	return sData
}