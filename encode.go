package logit

import (
	"errors"
	"reflect"
	"strings"
)

func Encode(input any) (any, error) {
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Map:
		return encodeMap(v)
	case reflect.Array, reflect.Slice:
		return encodeArray(v)
	case reflect.Struct:
		return encodeStruct(v)
	default:
	}
	return nil, errors.New("input type not supported")
}

func encodeMap(input reflect.Value) (map[string]any, error) {
	if input.Kind() != reflect.Map {
		return nil, errors.New("input should be a map")
	}

	if input.Type().Key().Kind() != reflect.String {
		return nil, errors.New("map keys must be strings")
	}

	resultMap := map[string]any{}

	iter := input.MapRange()
	for iter.Next() {
		mapKey := iter.Key()

		var key string
		if mapKey.Kind() == reflect.String {
			key = mapKey.String()
		} else {
			return nil, errors.New("map keys must be strings")
		}

		mapValue := iter.Value()
		val := mapValue.Interface()
		switch mapValue.Kind() {
		case reflect.Map:
			encodedVal, err := encodeMap(reflect.ValueOf(val))
			if err != nil {
				return nil, err
			}
			resultMap[key] = encodedVal
		case reflect.Array:
			encodedVal, err := encodeArray(reflect.ValueOf(val))
			if err != nil {
				return nil, err
			}
			resultMap[key] = encodedVal
		case reflect.Struct:
			encodedVal, err := encodeStruct(reflect.ValueOf(val))
			if err != nil {
				return nil, err
			}
			resultMap[key] = encodedVal
		default:
			resultMap[key] = mapValue.Interface()
		}
	}

	return resultMap, nil
}

func encodeArray(input reflect.Value) ([]any, error) {
	if input.Kind() != reflect.Slice && input.Kind() != reflect.Slice {
		return nil, errors.New("input should be a slice")
	}

	var resultArray []any

	for i := 0; i < input.Len(); i++ {
		val := input.Index(i)
		switch val.Kind() {
		case reflect.Map:
			encodedVal, err := encodeMap(reflect.ValueOf(val))
			if err != nil {
				return nil, err
			}
			resultArray = append(resultArray, encodedVal)
		case reflect.Struct:
			encodedVal, err := encodeStruct(reflect.ValueOf(val))
			if err != nil {
				return nil, err
			}
			resultArray = append(resultArray, encodedVal)
		default:
			resultArray = append(resultArray, val.Interface())
		}
	}

	return resultArray, nil
}

func encodeStruct(input reflect.Value) (map[string]any, error) {
	if input.Kind() != reflect.Struct {
		return nil, errors.New("input should be a struct")
	}

	resultMap := make(map[string]any)

	for i := 0; i < input.NumField(); i++ {
		field := input.Type().Field(i)
		fieldValue := input.Field(i)

		key := strings.ToLower(field.Name) // Default key is lowercase field name
		tagValue := field.Tag.Get("logit")
		tagOptions := strings.Split(tagValue, ",")
		if len(tagOptions) > 0 {
			key = tagOptions[0] // Use the first tag option as key if available
		}

		if len(tagOptions) > 1 && tagOptions[1] == "omitempty" && fieldValue.IsZero() {
			continue // Skip if the omitempty tag is present and field value is zero
		}

		switch fieldValue.Kind() {
		case reflect.Array, reflect.Slice:
			arr := make([]any, fieldValue.Len())
			for j := 0; j < fieldValue.Len(); j++ {
				elem := fieldValue.Index(j)
				if elem.Kind() == reflect.Map || elem.Kind() == reflect.Struct {
					mapVal, err := encodeStruct(elem)
					if err != nil {
						return nil, err
					}
					arr[j] = mapVal
				} else {
					arr[j] = elem.Interface()
				}
			}
			resultMap[key] = arr
		case reflect.Map:
			mapVal, err := encodeMap(fieldValue)
			if err != nil {
				return nil, err
			}
			resultMap[key] = mapVal
		case reflect.Struct:
			mapVal, err := encodeStruct(fieldValue)
			if err != nil {
				return nil, err
			}
			resultMap[key] = mapVal
		default:
			resultMap[key] = fieldValue.Interface()
		}
	}

	return resultMap, nil
}
