package cookies

import (
	"encoding/json"
	"time"

	contracts_config "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/contracts/config"
	fluffycore_contracts_cookies "github.com/fluffy-bunny/fluffycore/echo/contracts/cookies"
	status "github.com/gogo/status"
	echo "github.com/labstack/echo/v4"
	codes "google.golang.org/grpc/codes"
)

type (
	CustomCookieBase struct {
		fluffycore_contracts_cookies.ICookies
	}
)

func GetCookie[T any](c echo.Context,
	cookies fluffycore_contracts_cookies.ICookies, name string, data *T) error {
	getCookieResponse, err := cookies.GetCookie(c, name)
	if err != nil {
		return err
	}
	if getCookieResponse.Value == nil {
		return status.Errorf(codes.NotFound, "%s not found", name)
	}

	bb, err := json.Marshal(getCookieResponse.Value)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bb, data)
	return err
}
func SetCookieByRequest[T any](c echo.Context,
	config *contracts_config.CookieConfig,
	cookies fluffycore_contracts_cookies.ICookies,
	request *fluffycore_contracts_cookies.SetCookieRequest,
	data T) error {

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	value := make(map[string]interface{})
	err = json.Unmarshal(b, &value)
	if err != nil {
		return err
	}
	request.Path = "/"
	request.Expires = time.Now().Add(30 * time.Minute)
	request.Domain = config.Domain
	request.Value = value
	_, err = cookies.SetCookie(c, request)
	if err != nil {
		return err
	}
	return nil
}

func SetCookie[T any](c echo.Context,
	config *contracts_config.CookieConfig,
	cookies fluffycore_contracts_cookies.ICookies, name string, data T) error {

	return SetCookieByRequest(c, config, cookies,
		&fluffycore_contracts_cookies.SetCookieRequest{
			Name:     name,
			HttpOnly: false,
		},
		data)

}
