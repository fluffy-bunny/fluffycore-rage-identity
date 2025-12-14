package App

import (
	contracts_routes "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/oidc-login/internal/contracts/routes"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rs/zerolog"
)

func (s *service) Render() app.UI {
	return app.Div().Class("wizard-container").Body(
		s.renderHeader(s.AppContext),
		s.renderCurrentPage(),
		app.If(s.showCookieBanner, s.renderCookieBanner),
	)
}

func (s *service) renderCurrentPage() app.UI {
	switch s.currentPage {
	case contracts_routes.WellknownRoute_CreateAccount:
		return s.renderCreateAccountPage()
	case contracts_routes.WellknownRoute_Password:
		return s.renderPasswordPage()
	case contracts_routes.WellknownRoute_ResetPassword:
		return s.renderResetPasswordPage()
	case contracts_routes.WellknownRoute_VerifyCode:
		return s.renderVerifyCodePage()
	case contracts_routes.WellknownRoute_ForgotPassword:
		return s.renderForgotPasswordPage()
	default:
		return s.renderHomePage()
	}
}

func (s *service) renderHomePage() app.UI {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()
	log.Info().Msg("Rendering Home Page")
	return s.homeComposer
}

func (s *service) renderPasswordPage() app.UI {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()
	log.Info().Msg("Rendering Password Page")
	return s.passwordComposer
}

func (s *service) renderCreateAccountPage() app.UI {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()
	log.Info().Msg("Rendering Create Account Page")
	return s.createAccountComposer
}

func (s *service) renderForgotPasswordPage() app.UI {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()
	log.Info().Msg("Rendering Forgot Password Page")
	return s.forgotPasswordComposer
}

func (s *service) renderResetPasswordPage() app.UI {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()
	log.Info().Msg("Rendering Reset Password Page")
	return s.resetPasswordComposer
}

func (s *service) renderVerifyCodePage() app.UI {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "App").Logger()
	log.Info().Msg("Rendering Verify Code Page")
	return s.verifyCodeComposer
}
