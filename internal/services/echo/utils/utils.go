package utils

import (
	"fmt"
	"net/mail"

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
