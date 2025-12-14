package LinkedAccounts

import (
	"strconv"

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
					s.accounts[i] = linkedAccount{
						Identity:    identity.Subject,
						Provider:    identity.Provider,
						Email:       identity.Email,
						LinkedDate:  identity.LinkedAt,
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

	return app.Range(accountCards).Slice(func(i int) app.UI {
		return accountCards[i]
	})
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
				app.Div().Class("card-icon").Class(s.getProviderIconClass(account.Provider)).Body(
					app.Raw(s.getProviderIcon(account.Provider)),
				),
				app.Div().Class("card-title-group").Body(
					app.H2().Text(account.Provider),
					app.P().Class("card-description").Text(account.Email),
				),
			),
			app.If(!s.isClaimedDomain, func() app.UI {
				return app.Button().
					Class("btn-unlink").
					Disabled(account.IsUnlinking).
					DataSet("index", index).
					OnClick(s.handleUnlink).
					Text(func() string {
						if account.IsUnlinking {
							return s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyUnlinkingDotDot)
						}
						return s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyUnlink)
					}())
			}),
		),
		app.Div().Class("card-body").Body(
			app.Div().Class("info-rows").Body(
				app.Div().Class("info-row").Body(
					app.Span().Class("info-label").Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyLinkedOn)),
					app.Span().Class("info-value").Text(account.LinkedDate),
				),
			),
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

func (s *service) getProviderIconClass(provider string) string {
	switch provider {
	case "Google":
		return "google-icon"
	case "Microsoft":
		return "microsoft-icon"
	case "GitHub":
		return "github-icon"
	case "Facebook":
		return "facebook-icon"
	default:
		return "linked-accounts-icon"
	}
}

func (s *service) getProviderIcon(provider string) string {
	switch provider {
	case "Google":
		return `<svg xmlns="http://www.w3.org/2000/svg"  viewBox="0 0 48 48" width="48px" height="48px"><path fill="#FFC107" d="M43.611,20.083H42V20H24v8h11.303c-1.649,4.657-6.08,8-11.303,8c-6.627,0-12-5.373-12-12c0-6.627,5.373-12,12-12c3.059,0,5.842,1.154,7.961,3.039l5.657-5.657C34.046,6.053,29.268,4,24,4C12.955,4,4,12.955,4,24c0,11.045,8.955,20,20,20c11.045,0,20-8.955,20-20C44,22.659,43.862,21.35,43.611,20.083z"/><path fill="#FF3D00" d="M6.306,14.691l6.571,4.819C14.655,15.108,18.961,12,24,12c3.059,0,5.842,1.154,7.961,3.039l5.657-5.657C34.046,6.053,29.268,4,24,4C16.318,4,9.656,8.337,6.306,14.691z"/><path fill="#4CAF50" d="M24,44c5.166,0,9.86-1.977,13.409-5.192l-6.19-5.238C29.211,35.091,26.715,36,24,36c-5.202,0-9.619-3.317-11.283-7.946l-6.522,5.025C9.505,39.556,16.227,44,24,44z"/><path fill="#1976D2" d="M43.611,20.083H42V20H24v8h11.303c-0.792,2.237-2.231,4.166-4.087,5.571c0.001-0.001,0.002-0.001,0.003-0.002l6.19,5.238C36.971,39.205,44,34,44,24C44,22.659,43.862,21.35,43.611,20.083z"/></svg>`
	case "Microsoft":
		return `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 48 48" width="36px" height="36px"><path fill="#ff5722" d="M6 6H22V22H6z" transform="rotate(-180 14 14)"/><path fill="#4caf50" d="M26 6H42V22H26z" transform="rotate(-180 34 14)"/><path fill="#ffc107" d="M26 26H42V42H26z" transform="rotate(-180 34 34)"/><path fill="#03a9f4" d="M6 26H22V42H6z" transform="rotate(-180 14 34)"/></svg>`
	case "GitHub":
		return `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
    <path d="M9 19c-5 1.5-5-2.5-7-3m14 6v-3.87a3.37 3.37 0 0 0-.94-2.61c3.14-.35 6.44-1.54 6.44-7A5.44 5.44 0 0 0 20 4.77 5.07 5.07 0 0 0 19.91 1S18.73.65 16 2.48a13.38 13.38 0 0 0-7 0C6.27.65 5.09 1 5.09 1A5.07 5.07 0 0 0 5 4.77a5.44 5.44 0 0 0-1.5 3.78c0 5.42 3.3 6.61 6.44 7A3.37 3.37 0 0 0 9 18.13V22"/>
</svg>`

	default:
		return `<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
			<path d="M16 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"></path>
			<circle cx="8.5" cy="7" r="4"></circle>
			<line x1="20" y1="8" x2="20" y2="14"></line>
			<line x1="23" y1="11" x2="17" y2="11"></line>
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
