package routes

import (
	"strings"

	app "github.com/maxence-charriere/go-app/v10/pkg/app"
)

type WellknownRoute string

var (
	WellknownRoute_Home           WellknownRoute = "/"
	WellknownRoute_CreateAccount  WellknownRoute = "/create-account"
	WellknownRoute_Password       WellknownRoute = "/password"
	WellknownRoute_Passkey        WellknownRoute = "/passkey"
	WellknownRoute_ForgotPassword WellknownRoute = "/forgot-password"
	WellknownRoute_ResetPassword  WellknownRoute = "/reset-password"
	WellknownRoute_VerifyCode     WellknownRoute = "/verify-code"
)

func GetFixedRoute(route WellknownRoute) string {
	return FixHRef(string(route))
}
func FixHRef(href string) string {
	rootPrefix := app.Getenv("GOAPP_ROOT_PREFIX")

	if rootPrefix == "" || rootPrefix == "/" {
		return href
	}
	rr := strings.TrimRight(rootPrefix, "/")
	return rr + href

}
