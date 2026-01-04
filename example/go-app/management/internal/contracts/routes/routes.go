package routes

import (
	"strings"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WellknownRoute string

var (
	WellknownRoute_Home            WellknownRoute = "/"
	WellknownRoute_Profile         WellknownRoute = "/profile"
	WellknownRoute_PasswordManager WellknownRoute = "/password-manager"
	WellknownRoute_PasskeyManager  WellknownRoute = "/passkey-manager"
	WellknownRoute_LinkedAccounts  WellknownRoute = "/linked-accounts"
	WellknownRoute_Preferences     WellknownRoute = "/preferences"
)

func GetFixedRoute(route WellknownRoute) string {
	rootPrefix := app.Getenv("GOAPP_ROOT_PREFIX")

	if rootPrefix == "" || rootPrefix == "/" {
		return string(route)
	}
	rr := strings.TrimRight(rootPrefix, "/")
	return rr + string(route)
}
