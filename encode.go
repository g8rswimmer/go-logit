package logit

import (
	"errors"
	"fmt"
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
	if k := input.Kind(); k != reflect.Map {
		return nil, fmt.Errorf("input must be a map but is %v", k)
	}

	if k := input.Type().Key().Kind(); k != reflect.String {
		return nil, fmt.Errorf("map keys must be strings but are %v", k)
	}

	resultMap := map[string]any{}

	iter := input.MapRange()
	for iter.Next() {
		mapKey := iter.Key()

		var key string
		kk := mapKey.Kind()
		switch kk {
		case reflect.String:
			key = mapKey.String()
		default:
			return nil, fmt.Errorf("map keys must be strings but are %v", kk)
		}

		mapValue := iter.Value()
		if mapValue.Kind() == reflect.Interface {
			mapValue = reflect.ValueOf(mapValue.Interface())
		}
		if !mapValue.IsValid() {
			return nil, fmt.Errorf("the map key %s value, %v, is not valid or supported at this time", key, mapValue.Kind())
		}
		val := mapValue.Interface()
		switch mapValue.Kind() {
		case reflect.Map:
			encodedVal, err := encodeMap(reflect.ValueOf(val))
			if err != nil {
				return nil, errors.Join(fmt.Errorf("the map key %s had an error encoding a map", key), err)
			}
			resultMap[key] = encodedVal
		case reflect.Array, reflect.Slice:
			encodedVal, err := encodeArray(reflect.ValueOf(val))
			if err != nil {
				return nil, errors.Join(fmt.Errorf("the map key %s had an error encoding an array or slice", key), err)
			}
			resultMap[key] = encodedVal
		case reflect.Struct:
			encodedVal, err := encodeStruct(reflect.ValueOf(val))
			if err != nil {
				return nil, errors.Join(fmt.Errorf("the map key %s had an error encoding a struct", key), err)
			}
			resultMap[key] = encodedVal
		default:
			resultMap[key] = mapValue.Interface()
		}
	}

	return resultMap, nil
}

func encodeArray(input reflect.Value) ([]any, error) {
	if input.Kind() != reflect.Slice && input.Kind() != reflect.Array {
		return nil, fmt.Errorf("input should be a slice or array, but got %v", input.Kind())
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
			encodedVal, err := encodeStruct(val)
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
		if len(tagOptions) > 0 && len(tagOptions[0]) > 0 {
			key = tagOptions[0] // Use the first tag option as key if available
		}

		if len(tagOptions) > 1 && tagOptions[1] == "omitempty" {
			continue // Skip if the omitempty tag is present and field value is zero
		}
		switch fieldValue.Kind() {
		case reflect.Array, reflect.Slice:
			arrVale, err := encodeArray(fieldValue)
			if err != nil {
				return nil, err
			}
			resultMap[key] = arrVale
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
			if fieldValue.CanInterface() {
				resultMap[key] = fieldValue.Interface()
			}
		}
	}

	return resultMap, nil
}
