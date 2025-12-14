package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	app "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type CookieOptions struct {
	Path     string
	MaxAge   int // in seconds, 0 means session cookie, -1 deletes cookie
	Domain   string
	Secure   bool
	HttpOnly bool   // Note: cannot be set via JavaScript
	SameSite string // "Strict", "Lax", or "None"
}

// DefaultCookieOptions returns sensible defaults
func DefaultCookieOptions() CookieOptions {
	return CookieOptions{
		Path:     "/",
		MaxAge:   3600, // 1 hour
		SameSite: "Lax",
	}
}

// GetCookie reads a cookie and decodes it into type T
// For structs: expects base64-encoded JSON
// For primitives: reads the value directly
func GetCookie[T any](name string) (T, error) {
	var result T

	// Get all cookies
	cookies := app.Window().Get("document").Get("cookie").String()

	// Find the cookie
	for _, cookie := range strings.Split(cookies, ";") {
		cookie = strings.TrimSpace(cookie)
		if strings.HasPrefix(cookie, name+"=") {
			value := strings.TrimPrefix(cookie, name+"=")

			// Decode based on type
			return decodeCookieValue[T](value)
		}
	}

	return result, fmt.Errorf("cookie %s not found", name)
}

// SetCookie writes a cookie with the given value
// For structs: encodes as base64 JSON
// For primitives: stores the value directly
func SetCookie[T any](name string, value T, options CookieOptions) error {
	// Encode the value
	encodedValue, err := encodeCookieValue(value)
	if err != nil {
		return fmt.Errorf("failed to encode cookie value: %w", err)
	}

	// Build cookie string
	cookieStr := fmt.Sprintf("%s=%s", name, encodedValue)

	if options.Path != "" {
		cookieStr += fmt.Sprintf("; path=%s", options.Path)
	}

	if options.MaxAge != 0 {
		cookieStr += fmt.Sprintf("; max-age=%d", options.MaxAge)
	}

	if options.Domain != "" {
		cookieStr += fmt.Sprintf("; domain=%s", options.Domain)
	}

	if options.Secure {
		cookieStr += "; secure"
	}

	if options.SameSite != "" {
		cookieStr += fmt.Sprintf("; SameSite=%s", options.SameSite)
	}

	// Set the cookie
	app.Window().Get("document").Set("cookie", cookieStr)

	return nil
}

// GetOrSetCookie gets a cookie value, or sets and returns the default value if not found
func GetOrSetCookie[T any](name string, defaultValue T, options CookieOptions) (T, error) {
	// Try to get existing cookie
	value, err := GetCookie[T](name)
	if err == nil {
		return value, nil
	}

	// Cookie doesn't exist, set the default value
	err = SetCookie(name, defaultValue, options)
	if err != nil {
		return defaultValue, fmt.Errorf("failed to set default cookie: %w", err)
	}

	return defaultValue, nil
} // DeleteCookie removes a cookie by setting max-age to -1
func DeleteCookie(name string, options CookieOptions) {
	options.MaxAge = -1
	// Use empty string as value when deleting
	SetCookie(name, "", options)
}

// encodeCookieValue encodes a value for storage in a cookie
func encodeCookieValue[T any](value T) (string, error) {
	// Get the type of T
	rt := reflect.TypeOf(value)

	// Handle nil pointers
	if rt != nil && rt.Kind() == reflect.Ptr {
		rv := reflect.ValueOf(value)
		if rv.IsNil() {
			return "", nil
		}
		// Dereference pointer
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

	// For structs and other complex types: JSON + base64
	jsonBytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(jsonBytes), nil
}

// decodeCookieValue decodes a cookie value into type T
func decodeCookieValue[T any](value string) (T, error) {
	var result T

	if value == "" {
		return result, fmt.Errorf("cookie value is empty")
	}

	// Get the type of T
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

	// For structs and other complex types: decode base64 + JSON
	decodedBytes, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return result, fmt.Errorf("failed to decode base64: %w", err)
	}

	if isPtr {
		result = reflect.New(rt).Interface().(T)
		err = json.Unmarshal(decodedBytes, result)
	} else {
		err = json.Unmarshal(decodedBytes, &result)
	}

	if err != nil {
		return result, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}
