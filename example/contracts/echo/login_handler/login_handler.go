package login_handler

import (
	login_models "github.com/fluffy-bunny/fluffycore-rage-identity/pkg/models/api/login_models"
	echo "github.com/labstack/echo/v4"
)

type (
	ILoginHandler interface {
		HandleLogin(c echo.Context, request *login_models.LoginRequest) (*login_models.LoginResponse, error)
	}
)
