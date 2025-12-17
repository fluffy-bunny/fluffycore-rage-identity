package LinkedAccounts

import (
	"strconv"
	"strings"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	contracts_go_app_ManagementApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/ManagementApiClient"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/App"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/Localizer"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/LocalizerBundle"
	services_ComposerBase "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/ComposerBase"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

type (
	linkedAccount struct {
		Identity    string
		Provider    string
		Email       string
		LinkedDate  string
		IsUnlinking bool
	}

	service struct {
		services_ComposerBase.ComposerBase

		managementApiClient contracts_go_app_ManagementApiClient.IManagementApiClient
		accounts            []linkedAccount
		isClaimedDomain     bool
		errorMessage        string
		showError           bool
		isLoading           bool
	}
)

var stemService = (*service)(nil)

var _ contracts_App.ILinkedAccountsComposer = stemService

func (s *service) Ctor(
	container di.Container,
	appContext contracts_App.AppContext,
	localizer contracts_Localizer.ILocalizer,
	managementApiClient contracts_go_app_ManagementApiClient.IManagementApiClient,
) (contracts_App.ILinkedAccountsComposer, error) {

	return &service{
		ComposerBase: services_ComposerBase.ComposerBase{
			Container:  container,
			AppContext: appContext,
			Localizer:  localizer,
		},
		managementApiClient: managementApiClient,
		accounts:            []linkedAccount{},
		isLoading:           true,
	}, nil
}

func AddScopedILinkedAccountsComposer(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.ILinkedAccountsComposer](cb, stemService.Ctor)
}

func (s *service) OnMount(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "LinkedAccounts").Logger()
	log.Info().Msg("LinkedAccounts page mounted, fetching linked accounts")

	// Fetch linked accounts from API
	ctx.Async(func() {
		response, err := s.managementApiClient.GetUserLinkedAccounts(s.AppContext)
		ctx.Dispatch(func(ctx app.Context) {
			s.isLoading = false

			if err != nil {
				log.Error().Err(err).Msg("Failed to fetch linked accounts")
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyFailedToLoadLinkedAccounts)
				s.showError = true
				ctx.Update()
				return
			}

			if response != nil && response.Code == 200 && response.Response != nil {
				s.isClaimedDomain = response.Response.IsClaimedDomain
				log.Info().
					Int("count", len(response.Response.Identities)).
					Bool("isClaimedDomain", s.isClaimedDomain).
					Msg("Linked accounts loaded")

				// Convert API response to internal linkedAccount struct
				s.accounts = make([]linkedAccount, len(response.Response.Identities))
				for i, identity := range response.Response.Identities {
					linkedDate := ""
					if identity.CreatedOn > 0 {
						linkedDate = strconv.FormatInt(identity.CreatedOn, 10)
					}
					s.accounts[i] = linkedAccount{
						Identity:    identity.Subject,
						Provider:    identity.Provider,
						Email:       identity.Email,
						LinkedDate:  linkedDate,
						IsUnlinking: false,
					}
				}
			} else {
				log.Warn().Int("code", response.Code).Msg("Failed to fetch linked accounts")
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyFailedToLoadLinkedAccounts)
				s.showError = true
			}

			ctx.Update()
		})
	})
}

func (s *service) OnNav(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "LinkedAccounts").Logger()
	log.Info().Msg("LinkedAccounts page navigated")
}

func (s *service) OnDismount() {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "LinkedAccounts").Logger()
	log.Info().Msg("LinkedAccounts page dismounted")
}

func (s *service) Render() app.UI {
	return app.Div().Class("profile-container").Body(
		app.If(s.showError, s.renderErrorBanner),

		// Page Header
		app.Div().Class("profile-header").Body(
			app.H1().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyLinkedAccounts)),
			app.P().Class("profile-subtitle").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyManageYourLinkedAccounts)),
		),

		// Accounts Cards Container
		app.Div().Class("profile-cards").Body(
			app.If(s.isLoading, func() app.UI {
				return app.Div().Class("profile-card").Body(
					app.Div().Class("card-body").Style("text-align", "center").Style("padding", "40px 20px").Body(
						app.P().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyLoadingLinkedAccountsDotDot)),
					),
				)
			}).Else(func() app.UI {
				return s.renderAccountsList()
			}),
		),
	)
}

func (s *service) renderAccountsList() app.UI {
	if len(s.accounts) == 0 {
		return s.renderEmptyState()
	}

	accountCards := make([]app.UI, len(s.accounts))
	for i := range s.accounts {
		accountCards[i] = s.renderAccountCard(i)
	}

	return app.Div().Style("display", "flex").Style("flex-direction", "column").Style("gap", "1rem").Body(accountCards...)
}

func (s *service) renderAccountCard(index int) app.UI {
	account := s.accounts[index]
	log := zerolog.Ctx(s.AppContext).With().Str("component", "LinkedAccounts").Logger()
	log.Info().
		Int("index", index).
		Bool("isClaimedDomain", s.isClaimedDomain).
		Str("provider", account.Provider).
		Msg("Rendering account card")

	return app.Div().Class("profile-card").Body(
		app.Div().Class("card-header").Body(
			app.Div().Class("card-header-content").Body(
				app.Div().
					Style("background", "white").
					Style("border-radius", "50%").
					Style("width", "48px").
					Style("height", "48px").
					Style("min-width", "48px").
					Style("min-height", "48px").
					Style("display", "flex").
					Style("align-items", "center").
					Style("justify-content", "center").
					Style("padding", "8px").
					Style("box-shadow", "0 2px 8px rgba(0,0,0,0.1)").
					Body(
						app.Raw(s.getProviderIcon(account.Provider)),
					),
				app.Div().Class("card-title-group").Body(
					app.H2().Text(account.Provider+" ("+account.Email+")"),
					app.P().Class("card-description").Text("Linked on "+account.LinkedDate),
				),
			),
		),
		app.Div().Class("card-body").Body(
			app.If(!s.isClaimedDomain, func() app.UI {
				return app.Div().Class("button-group").Body(
					app.Button().
						Class("btn-unlink").
						Disabled(account.IsUnlinking).
						DataSet("index", index).
						OnClick(s.handleUnlink).
						Text(func() string {
							if account.IsUnlinking {
								return s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyUnlinkingDotDot)
							}
							return s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyUnlink)
						}()),
				)
			}),
		),
	)
}

func (s *service) renderEmptyState() app.UI {
	return app.Div().Class("profile-card").Body(
		app.Div().Class("card-body").Body(
			app.Div().Style("text-align", "center").Style("padding", "40px 20px").Body(
				app.Div().Class("home-feature-icon").Class("accounts-icon").Style("margin", "0 auto 24px").Body(
					app.Raw(`<svg width="48" height="48" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path>
						<path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path>
					</svg>`),
				),
				app.H3().Style("margin", "0 0 8px 0").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyNoLinkedAccounts)),
				app.P().Class("card-description").Style("color", "var(--text-secondary)").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyNoLinkedAccountsDescription)),
			),
		),
	)
}

func (s *service) renderErrorBanner() app.UI {
	return app.Div().Class("error-banner").Body(
		app.Div().Class("error-content").Body(
			app.Div().Class("error-icon").Body(
				app.Raw(`<svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
				</svg>`),
			),
			app.Div().Class("error-text").Body(
				app.Span().Class("error-title").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyError)),
				app.Span().Class("error-message").Text(s.errorMessage),
			),
		),
	)
}

func (s *service) getProviderIcon(provider string) string {
	lowerProvider := strings.ToLower(provider)
	if strings.Contains(lowerProvider, "google") {
		return `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
			<path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
			<path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
			<path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z" fill="#FBBC05"/>
			<path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
		</svg>`
	} else if strings.Contains(lowerProvider, "microsoft") {
		return `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
			<path d="M11.4 11.4H2V2h9.4v9.4z" fill="#F25022"/>
			<path d="M22 11.4h-9.4V2H22v9.4z" fill="#7FBA00"/>
			<path d="M11.4 22H2v-9.4h9.4V22z" fill="#00A4EF"/>
			<path d="M22 22h-9.4v-9.4H22V22z" fill="#FFB900"/>
		</svg>`
	} else if strings.Contains(lowerProvider, "github") {
		return `<svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
			<path d="M12 2C6.477 2 2 6.477 2 12c0 4.42 2.865 8.17 6.839 9.49.5.092.682-.217.682-.482 0-.237-.008-.866-.013-1.7-2.782.603-3.369-1.34-3.369-1.34-.454-1.156-1.11-1.463-1.11-1.463-.908-.62.069-.608.069-.608 1.003.07 1.531 1.03 1.531 1.03.892 1.529 2.341 1.087 2.91.831.092-.646.35-1.086.636-1.336-2.22-.253-4.555-1.11-4.555-4.943 0-1.091.39-1.984 1.029-2.683-.103-.253-.446-1.27.098-2.647 0 0 .84-.269 2.75 1.025A9.578 9.578 0 0112 6.836c.85.004 1.705.114 2.504.336 1.909-1.294 2.747-1.025 2.747-1.025.546 1.377.203 2.394.1 2.647.64.699 1.028 1.592 1.028 2.683 0 3.842-2.339 4.687-4.566 4.935.359.309.678.919.678 1.852 0 1.336-.012 2.415-.012 2.743 0 .267.18.578.688.48C19.137 20.167 22 16.418 22 12c0-5.523-4.477-10-10-10z"/>
		</svg>`
	} else {
		// Default chain/link icon for unknown providers
		return `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#374151" stroke-width="2">
			<path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71"></path>
			<path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71"></path>
		</svg>`
	}
}

// Event handlers
func (s *service) handleUnlink(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "LinkedAccounts").Logger()

	indexStr := ctx.JSSrc().Get("dataset").Get("index").String()
	index, err := strconv.Atoi(indexStr)
	if err != nil || index < 0 || index >= len(s.accounts) {
		log.Warn().Str("indexStr", indexStr).Int("index", index).Msg("Invalid account index")
		return
	}

	account := &s.accounts[index]
	identity := account.Identity
	log.Info().Str("provider", account.Provider).Str("email", account.Email).Str("identity", identity).Msg("Unlinking account")

	account.IsUnlinking = true
	s.showError = false
	ctx.Update()

	// Call API to unlink account
	ctx.Async(func() {
		response, err := s.managementApiClient.DeleteUserLinkedAccount(s.AppContext, identity)
		ctx.Dispatch(func(ctx app.Context) {
			if err != nil {
				log.Error().Err(err).Str("identity", identity).Msg("Failed to unlink account")
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyFailedToUnlinkAccount)
				s.showError = true
				account.IsUnlinking = false
				ctx.Update()
				return
			}

			if response != nil && response.Code == 200 {
				log.Info().Str("provider", account.Provider).Msg("Account unlinked successfully")
				// Remove account from list
				s.accounts = append(s.accounts[:index], s.accounts[index+1:]...)
				s.showError = false
			} else {
				log.Warn().Int("code", response.Code).Msg("Failed to unlink account")
				s.errorMessage = s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyFailedToUnlinkAccount)
				s.showError = true
				account.IsUnlinking = false
			}
			ctx.Update()
		})
	})
}
