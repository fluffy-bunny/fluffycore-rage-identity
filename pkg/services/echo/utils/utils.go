package utils

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/mail"
	"time"

	contracts_identity "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/identity"
	echo "github.com/labstack/echo/v4"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
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
func GetLocalizerFromEchoContext(b *i18n.Bundle, e echo.Context) *i18n.Localizer {
	accept := e.Request().Header.Get("Accept-Language")
	return i18n.NewLocalizer(b, accept)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GenerateRandomAlphaNumericString(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

type VerificationCodeResult struct {
	PlainCode  string
	HashedCode string
}

// GenerateHashedVerificationCode generates a random alphanumeric verification code
// and returns both the plain text code (for sending in email) and its hashed version (for storing in cookie)
func GenerateHashedVerificationCode(ctx context.Context, passwordHasher contracts_identity.IPasswordHasher, length int) (*VerificationCodeResult, error) {
	plainCode := GenerateRandomAlphaNumericString(length)

	hashResponse, err := passwordHasher.HashPassword(ctx, &contracts_identity.HashPasswordRequest{
		Password: plainCode,
	})
	if err != nil {
		return nil, err
	}

	return &VerificationCodeResult{
		PlainCode:  plainCode,
		HashedCode: hashResponse.HashedPassword,
	}, nil
}

// VerifyVerificationCode verifies that the provided plain text code matches the hashed code
func VerifyVerificationCode(ctx context.Context, passwordHasher contracts_identity.IPasswordHasher, plainCode string, hashedCode string) error {
	return passwordHasher.VerifyPassword(ctx, &contracts_identity.VerifyPasswordRequest{
		Password:       plainCode,
		HashedPassword: hashedCode,
	})
}
