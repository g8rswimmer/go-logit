package logit

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

func encode(input any) (any, error) {
	v := reflect.ValueOf(input)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if k := v.Kind(); k != reflect.Map && k != reflect.Struct {
		return nil, fmt.Errorf("input must be a map or struct but is %v", k)
	}
	switch v.Kind() {
	case reflect.Map:
		return encodeMap(v)
	default:
		return encodeStruct(v)
	}
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
		if mapValue.Kind() == reflect.Ptr {
			mapValue = mapValue.Elem()
		}
		if !mapValue.IsValid() {
			return nil, fmt.Errorf("the map key %s value, %v, is not valid or supported at this time", key, mapValue.Kind())
		}
		if v, has := retrieveValue(mapValue); has {
			resultMap[key] = v
			continue
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
			resultMap[key] = val
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
		if val.Kind() == reflect.Ptr {
			val = val.Elem()
		}
		if v, has := retrieveValue(val); has {
			resultArray = append(resultArray, v)
			continue
		}
		switch val.Kind() {
		case reflect.Map:
			encodedVal, err := encodeMap(val)
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
		case reflect.Array, reflect.Slice:
			encodedVal, err := encodeArray(val)
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
	if input.Kind() == reflect.Ptr {
		input = input.Elem()
	}

	resultMap := make(map[string]any)

	for i := 0; i < input.NumField(); i++ {
		field := input.Type().Field(i)
		fieldValue := input.Field(i)

		ft := encodeFieldTag(field)

		if ft.Omit {
			continue // Skip if the omitempty tag is present and field value is zero
		}
		if fieldValue.CanInterface() && fieldValue.Kind() == reflect.Ptr {
			fieldValue = fieldValue.Elem()
		}
		if val, has := handleTags(ft, fieldValue); has {
			resultMap[ft.Name] = val
			continue
		}
		if val, has := retrieveValue(fieldValue); has {
			resultMap[ft.Name] = val
			continue
		}
		switch fieldValue.Kind() {
		case reflect.Array, reflect.Slice:
			arrVale, err := encodeArray(fieldValue)
			if err != nil {
				return nil, err
			}
			resultMap[ft.Name] = arrVale
		case reflect.Map:
			mapVal, err := encodeMap(fieldValue)
			if err != nil {
				return nil, err
			}
			resultMap[ft.Name] = mapVal
		case reflect.Struct:
			mapVal, err := encodeStruct(fieldValue)
			if err != nil {
				return nil, err
			}
			resultMap[ft.Name] = mapVal
		default:
			if fieldValue.CanInterface() {
				resultMap[ft.Name] = fieldValue.Interface()
			}
		}
	}

	return resultMap, nil
}

func retrieveValue(input reflect.Value) (any, bool) {
	if !input.CanInterface() {
		return nil, false
	}
	if input.Kind() == reflect.Ptr {
		input = input.Elem()
	}
	val := input.Interface()
	switch v := val.(type) {
	case time.Time:
		return v.Format(time.RFC3339), true
	case *time.Time:
		return v.Format(time.RFC3339), true
	default:
		return nil, false
	}
}

func handleTags(ft *fieldTag, input reflect.Value) (any, bool) {
	switch {
	case ft.Obfuscate:
		return "xxxx", true
	}
	return nil, false
}
