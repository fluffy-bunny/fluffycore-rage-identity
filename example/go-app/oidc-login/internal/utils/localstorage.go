package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	app "github.com/maxence-charriere/go-app/v10/pkg/app"
)

// GetLocalStorage reads a value from localStorage and decodes it into type T
// For structs: expects JSON
// For primitives: reads the value directly
func GetLocalStorage[T any](key string) (T, error) {
	var result T

	storage := app.Window().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return result, fmt.Errorf("localStorage is not available")
	}

	value := storage.Call("getItem", key)
	if value.IsNull() || value.IsUndefined() {
		return result, fmt.Errorf("key %s not found in localStorage", key)
	}

	valueStr := value.String()
	if valueStr == "" {
		return result, fmt.Errorf("key %s has empty value", key)
	}

	return decodeLocalStorageValue[T](valueStr)
}

// SetLocalStorage writes a value to localStorage
// For structs: encodes as JSON
// For primitives: stores the value directly
func SetLocalStorage[T any](key string, value T) error {
	storage := app.Window().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return fmt.Errorf("localStorage is not available")
	}

	encodedValue, err := encodeLocalStorageValue(value)
	if err != nil {
		return fmt.Errorf("failed to encode value: %w", err)
	}

	storage.Call("setItem", key, encodedValue)
	return nil
}

// GetOrSetLocalStorage gets a value from localStorage, or sets and returns the default value if not found
func GetOrSetLocalStorage[T any](key string, defaultValue T) (T, error) {
	// Try to get existing value
	value, err := GetLocalStorage[T](key)
	if err == nil {
		return value, nil
	}

	// Key doesn't exist, set the default value
	err = SetLocalStorage(key, defaultValue)
	if err != nil {
		return defaultValue, fmt.Errorf("failed to set default value: %w", err)
	}

	return defaultValue, nil
}

// RemoveLocalStorage removes a key from localStorage
func RemoveLocalStorage(key string) error {
	storage := app.Window().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return fmt.Errorf("localStorage is not available")
	}

	storage.Call("removeItem", key)
	return nil
}

// ClearLocalStorage removes all keys from localStorage
func ClearLocalStorage() error {
	storage := app.Window().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return fmt.Errorf("localStorage is not available")
	}

	storage.Call("clear")
	return nil
}

// HasLocalStorageKey checks if a key exists in localStorage
func HasLocalStorageKey(key string) bool {
	storage := app.Window().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return false
	}

	value := storage.Call("getItem", key)
	return !value.IsNull() && !value.IsUndefined()
}

// GetLocalStorageKeys returns all keys in localStorage
func GetLocalStorageKeys() []string {
	storage := app.Window().Get("localStorage")
	if storage.IsUndefined() || storage.IsNull() {
		return []string{}
	}

	length := storage.Get("length")
	if length.IsUndefined() || length.IsNull() {
		return []string{}
	}

	keys := []string{}
	for i := 0; i < length.Int(); i++ {
		keyValue := storage.Call("key", i)
		if !keyValue.IsNull() && !keyValue.IsUndefined() {
			keys = append(keys, keyValue.String())
		}
	}

	return keys
}

// encodeLocalStorageValue encodes a value for storage in localStorage
func encodeLocalStorageValue[T any](value T) (string, error) {
	rt := reflect.TypeOf(value)

	// Handle nil pointers
	if rt != nil && rt.Kind() == reflect.Ptr {
		rv := reflect.ValueOf(value)
		if rv.IsNil() {
			return "", nil
		}
		value = rv.Elem().Interface().(T)
		rt = reflect.TypeOf(value)
	}

	// Handle primitives
	switch rt.Kind() {
	case reflect.String:
		return reflect.ValueOf(value).String(), nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(reflect.ValueOf(value).Int(), 10), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(reflect.ValueOf(value).Uint(), 10), nil
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(reflect.ValueOf(value).Float(), 'f', -1, 64), nil
	case reflect.Bool:
		return strconv.FormatBool(reflect.ValueOf(value).Bool()), nil
	}

	// For structs and other complex types: JSON
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// decodeLocalStorageValue decodes a localStorage value into type T
func decodeLocalStorageValue[T any](value string) (T, error) {
	var result T

	if value == "" {
		return result, fmt.Errorf("value is empty")
	}

	rt := reflect.TypeOf(result)

	// Handle pointer types
	isPtr := false
	if rt != nil && rt.Kind() == reflect.Ptr {
		isPtr = true
		rt = rt.Elem()
	}

	// Handle primitives
	switch rt.Kind() {
	case reflect.String:
		if isPtr {
			str := value
			return reflect.ValueOf(&str).Interface().(T), nil
		}
		return reflect.ValueOf(value).Interface().(T), nil

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		intVal, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return result, err
		}
		if isPtr {
			v := reflect.New(rt)
			v.Elem().SetInt(intVal)
			return v.Interface().(T), nil
		}
		v := reflect.New(rt).Elem()
		v.SetInt(intVal)
		return v.Interface().(T), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		uintVal, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return result, err
		}
		if isPtr {
			v := reflect.New(rt)
			v.Elem().SetUint(uintVal)
			return v.Interface().(T), nil
		}
		v := reflect.New(rt).Elem()
		v.SetUint(uintVal)
		return v.Interface().(T), nil

	case reflect.Float32, reflect.Float64:
		floatVal, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return result, err
		}
		if isPtr {
			v := reflect.New(rt)
			v.Elem().SetFloat(floatVal)
			return v.Interface().(T), nil
		}
		v := reflect.New(rt).Elem()
		v.SetFloat(floatVal)
		return v.Interface().(T), nil

	case reflect.Bool:
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return result, err
		}
		if isPtr {
			v := reflect.New(rt)
			v.Elem().SetBool(boolVal)
			return v.Interface().(T), nil
		}
		v := reflect.New(rt).Elem()
		v.SetBool(boolVal)
		return v.Interface().(T), nil
	}

	// For structs and other complex types: decode JSON
	if isPtr {
		result = reflect.New(rt).Interface().(T)
		err := json.Unmarshal([]byte(value), result)
		if err != nil {
			return result, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	} else {
		err := json.Unmarshal([]byte(value), &result)
		if err != nil {
			return result, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}
	}

	return result, nil
}
