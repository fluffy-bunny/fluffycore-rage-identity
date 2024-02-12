package utils

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"

	argon2id "github.com/alexedwards/argon2id"
	echo "github.com/labstack/echo/v4"
)

func GetMyRootPath(c echo.Context) string {
	return fmt.Sprintf("%s://%s", c.Scheme(), c.Request().Host)
}
func IsValidEmailAddress(address string) (string, bool) {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return "", false
	}
	return addr.Address, true
}
func GeneratePasswordHash(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}
func ComparePasswordHash(password string, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}
func DeleteCookie(c echo.Context, name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	c.SetCookie(cookie)
}
func SetCookieInterface(c echo.Context, cookie *http.Cookie, value interface{}) {
	cookieData, _ := json.Marshal(value)
	encodedValue := base64.StdEncoding.EncodeToString([]byte(cookieData))
	cookie.Value = encodedValue
	c.SetCookie(cookie)
}

func GetCookieInterface(c echo.Context, name string, v any) error {
	cookie, err := c.Cookie(name)
	if err != nil {
		return err
	}
	decodedValue, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return err
	}
	return json.Unmarshal(decodedValue, v)
}
