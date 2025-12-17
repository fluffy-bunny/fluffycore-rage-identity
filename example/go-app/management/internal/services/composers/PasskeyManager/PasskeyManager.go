package PasskeyManager

import (
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	go_app_common "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/common"
	contracts_go_app_ManagementApiClient "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/contracts/ManagementApiClient"
	contracts_App "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/App"
	contracts_Localizer "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/Localizer"
	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/LocalizerBundle"
	services_ComposerBase "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/services/ComposerBase"
	models "github.com/fluffy-bunny/fluffycore-rage-identity/example/services/echo/account/models"
	"github.com/fluffy-bunny/fluffycore-rage-identity/pkg/go-app/common"
	app "github.com/maxence-charriere/go-app/v10/pkg/app"
	zerolog "github.com/rs/zerolog"
)

type (
	passkeyItem struct {
		ID           string
		FriendlyName string
		CreatedAt    string
		Transport    []string
		IsRenaming   bool
		EditName     string
	}

	service struct {
		services_ComposerBase.ComposerBase

		managementApiClient contracts_go_app_ManagementApiClient.IManagementApiClient
		passkeys            []passkeyItem
		errorMessage        string
		successMessage      string
		showError           bool
		showSuccess         bool
		isLoading           bool
		isAddingPasskey     bool
		isClaimedDomain     bool
		showDeleteConfirm   bool
		deleteCredentialID  string
		isDeleting          bool
	}
)

var stemService = (*service)(nil)

var _ contracts_App.IPasskeyManagerComposer = stemService

func (s *service) Ctor(
	container di.Container,
	appContext contracts_App.AppContext,
	localizer contracts_Localizer.ILocalizer,
	managementApiClient contracts_go_app_ManagementApiClient.IManagementApiClient,
) (contracts_App.IPasskeyManagerComposer, error) {

	return &service{
		ComposerBase: services_ComposerBase.ComposerBase{
			Container:  container,
			AppContext: appContext,
			Localizer:  localizer,
		},
		managementApiClient: managementApiClient,
		passkeys:            []passkeyItem{},
		isLoading:           true,
	}, nil
}

func AddScopedIPasskeyManagerComposer(cb di.ContainerBuilder) {
	di.AddScoped[contracts_App.IPasskeyManagerComposer](cb, stemService.Ctor)
}

func (s *service) OnMount(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PasskeyManager").Logger()
	log.Info().Msg("PasskeyManager page mounted")

	// Fetch user profile first to check if claimed domain
	ctx.Async(func() {
		response, err := s.managementApiClient.GetUserProfile(s.AppContext)
		ctx.Dispatch(func(ctx app.Context) {
			if err == nil && response != nil && response.Code == 200 && response.Response != nil {
				s.isClaimedDomain = response.Response.IsClaimedDomain
				log.Info().Bool("isClaimedDomain", s.isClaimedDomain).Msg("Profile loaded")

				// Only load passkeys if not a claimed domain
				if !s.isClaimedDomain {
					s.loadPasskeys(ctx)
				} else {
					s.isLoading = false
					ctx.Update()
				}
			} else {
				log.Error().Err(err).Msg("Failed to load profile")
				s.isLoading = false
				ctx.Update()
			}
		})
	})
}

// loadPasskeys fetches passkeys asynchronously (for initial page load)
func (s *service) loadPasskeys(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PasskeyManager").Logger()
	log.Info().Msg("🔵 loadPasskeys() - using ctx.Async + HTTP")

	ctx.Async(func() {
		response, err := s.managementApiClient.GetPasskeysHTTP(s.AppContext)

		ctx.Dispatch(func(ctx app.Context) {
			s.isLoading = false

			if err != nil {
				log.Error().Err(err).Msg("Failed to fetch passkeys")
				s.errorMessage = "Failed to load passkeys. Please try again."
				s.showError = true
			} else if response != nil && response.Code == 200 {
				var passkeys []models.PasskeyItem
				if response.Response != nil && response.Response.Passkeys != nil {
					passkeys = response.Response.Passkeys
				}

				log.Info().Int("count", len(passkeys)).Msg("Passkeys loaded")

				s.passkeys = make([]passkeyItem, len(passkeys))
				for i, passkey := range passkeys {
					s.passkeys[i] = passkeyItem{
						ID:           passkey.ID,
						FriendlyName: passkey.FriendlyName,
						CreatedAt:    passkey.CreatedAt,
						Transport:    passkey.Transport,
						IsRenaming:   false,
						EditName:     passkey.FriendlyName,
					}
				}
			} else {
				s.errorMessage = "Failed to load passkeys"
				s.showError = true
			}
			ctx.Update()
		})
	})
}

// loadPasskeysSync fetches passkeys synchronously (call from within ctx.Dispatch)
func (s *service) loadPasskeysSync(ctx app.Context) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PasskeyManager").Logger()
	app.Log("🔵 loadPasskeysSync() - using HTTP (not fetch)")

	response, err := s.managementApiClient.GetPasskeysHTTP(s.AppContext)
	app.Log("🔵 GetPasskeysHTTP returned - err:", err)

	s.isLoading = false

	if err != nil {
		log.Error().Err(err).Msg("Failed to fetch passkeys")
		app.Log("❌ Failed to fetch passkeys")
		s.errorMessage = "Failed to load passkeys. Please try again."
		s.showError = true
	} else if response != nil && response.Code == 200 {
		var passkeys []models.PasskeyItem
		if response.Response != nil && response.Response.Passkeys != nil {
			passkeys = response.Response.Passkeys
		}

		log.Info().Int("count", len(passkeys)).Msg("Passkeys loaded")
		app.Log("✅ Loaded", len(passkeys), "passkeys via HTTP")

		s.passkeys = make([]passkeyItem, len(passkeys))
		for i, passkey := range passkeys {
			s.passkeys[i] = passkeyItem{
				ID:           passkey.ID,
				FriendlyName: passkey.FriendlyName,
				CreatedAt:    passkey.CreatedAt,
				Transport:    passkey.Transport,
				IsRenaming:   false,
				EditName:     passkey.FriendlyName,
			}
		}
	} else {
		app.Log("❌ Unexpected response")
		s.errorMessage = "Failed to load passkeys"
		s.showError = true
	}

	app.Log("🔵 Calling ctx.Update()")
	ctx.Update()
	app.Log("✅ UI updated")
}

func (s *service) Render() app.UI {
	// Show unavailable message for claimed domain users
	if s.isClaimedDomain {
		return app.Div().Class("profile-container").Body(
			app.Div().Class("profile-header").Body(
				app.H1().Text("Passkeys"),
				app.P().Class("profile-subtitle").Text("Manage your passkeys for passwordless authentication"),
			),
			app.Div().Class("profile-card").Body(
				app.Div().Class("card-header").Body(
					app.H2().Text(s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasskeysNotAvailable)),
				),
				app.Div().Class("card-body").Body(
					app.P().Class("card-description").Style("color", "var(--text-secondary)").Text(
						s.Localizer.GetLocalizedString(contracts_LocalizerBundle.LocaleKeyPasskeysNotAvailableDescription),
					),
				),
			),
		)
	}

	return app.Div().Class("profile-container").Body(
		app.If(s.showError, s.renderErrorBanner),
		app.If(s.showSuccess, s.renderSuccessNotification),
		app.If(s.showDeleteConfirm, s.renderDeleteConfirmation),

		// Page Header
		app.Div().Class("profile-header").Body(
			app.H1().Text("Passkeys"),
			app.P().Class("profile-subtitle").Text("Manage your passkeys for passwordless authentication"),
		),

		// Content
		app.Div().Class("profile-cards").Body(
			app.If(s.isLoading, s.renderLoading),
			app.If(!s.isLoading && len(s.passkeys) == 0, s.renderEmptyState),
			app.If(!s.isLoading && len(s.passkeys) > 0, s.renderPasskeysList),
		),
	)
}

func (s *service) renderLoading() app.UI {
	return app.Div().Class("profile-card").Body(
		app.Div().Class("card-body").Style("text-align", "center").Body(
			app.P().Text("Loading passkeys..."),
		),
	)
}

func (s *service) renderEmptyState() app.UI {
	return app.Div().Class("profile-card").Body(
		app.Div().Class("card-header").Body(
			app.Div().Class("card-header-content").Body(
				app.Div().Class("card-icon password-icon").Body(
					app.Raw(`<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
								<circle cx="7" cy="7" r="2"></circle>
								<path d="M7 9v4a2 2 0 0 0 2 2h4"></path>
								<circle cx="19" cy="15" r="4"></circle>
								<path d="M19 11v-1"></path>
								<path d="M22 15h-1"></path>
								</svg>`),
				),
				app.Div().Class("card-title-group").Body(
					app.H2().Text("No Passkeys Configured"),
					app.P().Class("card-description").Text("Add a passkey to enable passwordless sign-in with your fingerprint, face, or security key"),
				),
			),
		),
		app.Div().Class("card-body").Body(
			app.Div().Class("button-group").Body(
				app.Button().
					Class("btn-primary").
					Disabled(s.isAddingPasskey).
					OnClick(s.handleAddPasskey).
					Text(func() string {
						if s.isAddingPasskey {
							return "Adding Passkey..."
						}
						return "Add Your First Passkey"
					}()),
			),
		),
	)
}

func (s *service) renderPasskeysList() app.UI {
	items := make([]app.UI, 0, len(s.passkeys)+1)

	// Add button card first
	items = append(items, s.renderAddPasskeyCard())

	// Then add all passkey cards
	for i := range s.passkeys {
		idx := i
		items = append(items, s.renderPasskeyCard(&s.passkeys[idx], idx))
	}

	return app.Div().Body(items...)
}

func (s *service) renderAddPasskeyCard() app.UI {
	return app.Div().Class("profile-card").Body(
		app.Div().Class("card-header").Body(
			app.Div().Class("card-header-content").Body(
				app.Div().Class("card-icon passkey-icon").Body(
					app.Raw(`<svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
						<circle cx="7" cy="7" r="3"></circle>
						<path d="M5 22v-5l-1-1v-4a2 2 0 0 1 2-2h2a2 2 0 0 1 2 2v4l-1 1v5"></path>
						<path d="M14 6h8"></path>
						<path d="M18 2v8"></path>
					</svg>`),
				),
				app.Div().Class("card-title-group").Body(
					app.H2().Text("Add New Passkey"),
					app.P().Class("card-description").Text("Register a new device or security key for passwordless authentication"),
				),
			),
		),
		app.Div().Class("card-body").Body(
			app.Div().Class("button-group").Body(
				app.Button().
					Class("btn-primary").
					Disabled(s.isAddingPasskey).
					OnClick(s.handleAddPasskey).
					Text(func() string {
						if s.isAddingPasskey {
							return "Adding Passkey..."
						}
						return "Add Passkey"
					}()),
			),
		),
	)
}

func (s *service) renderPasskeyCard(passkey *passkeyItem, index int) app.UI {
	return app.Div().Class("profile-card").Body(
		app.Div().Class("card-header").Body(
			app.Div().Class("card-header-content").Body(
				app.Div().Class("card-icon passkey-icon").Body(
					app.Raw(go_app_common.PasskeyIconSmallSVG)),
				app.Div().Class("card-title-group").Body(
					app.If(!passkey.IsRenaming, func() app.UI {
						return app.H2().Text(passkey.FriendlyName)
					}),
					app.If(passkey.IsRenaming, func() app.UI {
						return app.Div().Class("form-group").Style("margin", "0").Body(
							app.Input().
								Type("text").
								Value(passkey.EditName).
								OnInput(s.handleEditNameInput(index)).
								Placeholder("Enter passkey name"),
						)
					}),
				),
			),
		),
		app.Div().Class("card-body").Body(
			app.Div().Class("button-group").Body(
				app.If(!passkey.IsRenaming, func() app.UI {
					return app.Button().
						Class("btn-secondary").
						OnClick(s.handleStartRename(index)).
						Text("Rename")
				}),
				app.If(passkey.IsRenaming, func() app.UI {
					return app.Button().
						Class("btn-primary").
						OnClick(s.handleSaveRename(index)).
						Text("Save")
				}),
				app.If(passkey.IsRenaming, func() app.UI {
					return app.Button().
						Class("btn-secondary").
						OnClick(s.handleCancelRename(index)).
						Text("Cancel")
				}),
				app.If(!passkey.IsRenaming, func() app.UI {
					return app.Button().
						Class("btn-unlink").
						OnClick(s.handleDeletePasskey(passkey.ID)).
						Text("Delete")
				}),
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
				app.Span().Class("error-title").Text("Error"),
				app.Span().Class("error-message").Text(s.errorMessage),
			),
		),
		app.Button().
			Class("error-close").
			OnClick(s.handleCloseError).
			Body(
				app.Raw(`<svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
				</svg>`),
			),
	)
}

func (s *service) renderSuccessNotification() app.UI {
	return app.Div().Class("success-notification").Body(
		app.Div().Class("success-notification-icon").Body(
			app.Raw(`<svg width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
				<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
			</svg>`),
		),
		app.Div().Class("success-notification-text").Text(s.successMessage),
		app.Button().
			Class("success-notification-close").
			OnClick(s.handleCloseSuccess).
			Body(
				app.Raw(`<svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
					<path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
				</svg>`),
			),
	)
}

func (s *service) renderDeleteConfirmation() app.UI {
	return app.Div().Class("modal-overlay").OnClick(s.handleCancelDelete).Body(
		app.Div().Class("modal-content").OnClick(func(ctx app.Context, e app.Event) {
			e.PreventDefault()
			// Stop event from bubbling to overlay which would close the modal
			e.JSValue().Call("stopPropagation")
		}).Body(
			app.Div().Class("modal-header").Body(
				app.H3().Text("Delete Passkey"),
				app.Button().Class("modal-close").OnClick(s.handleCancelDelete).Body(
					app.Raw(`<svg width="20" height="20" viewBox="0 0 20 20" fill="currentColor">
						<path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
					</svg>`),
				),
			),
			app.Div().Class("modal-body").Body(
				app.P().Text("Are you sure you want to delete this passkey? This action cannot be undone."),
			),
			app.Div().Class("modal-footer").Body(
				app.Button().
					Class("btn-secondary").
					OnClick(s.handleCancelDelete).
					Text("Cancel"),
				app.Button().
					Class("btn-unlink").
					Disabled(s.isDeleting).
					OnClick(s.handleConfirmDelete).
					Text(func() string {
						if s.isDeleting {
							return "Deleting..."
						}
						return "Delete"
					}()),
			),
		),
	)
}

// Event Handlers

func (s *service) handleAddPasskey(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PasskeyManager").Logger()
	log.Info().Msg("Starting passkey registration")

	s.isAddingPasskey = true
	s.showError = false
	s.showSuccess = false
	ctx.Update()

	// Call the WebAuthn registration flow via JavaScript
	ctx.Async(func() {
		log.Info().Msg("Inside ctx.Async - checking for registerPasskey function")

		// Get the window object using app.Window(), then access registerPasskey
		registerPasskeyFunc := app.Window().Get("registerPasskey")
		log.Info().
			Bool("truthy", registerPasskeyFunc.Truthy()).
			Str("type", registerPasskeyFunc.Type().String()).
			Msg("registerPasskey function check")

		if !registerPasskeyFunc.Truthy() || registerPasskeyFunc.Type() != app.TypeFunction {
			log.Error().
				Bool("truthy", registerPasskeyFunc.Truthy()).
				Str("type", registerPasskeyFunc.Type().String()).
				Msg("registerPasskey function not available!")
			ctx.Dispatch(func(ctx app.Context) {
				s.isAddingPasskey = false
				s.errorMessage = "WebAuthn not available. Please ensure webauthn.js is loaded."
				s.showError = true
				ctx.Update()
			})
			return
		}

		log.Info().Msg("registerPasskey function found, calling it now")
		// Call the global registerPasskey JavaScript function defined in webauthn.js
		// This returns a Promise, so we need to await it
		resultChan := make(chan bool, 1)

		ctx.Dispatch(func(ctx app.Context) {
			log.Info().Msg("Inside ctx.Dispatch - about to call registerPasskey")
			// Call registerPasskey() via app.Window()
			promise := app.Window().Call("registerPasskey", "My Passkey")
			log.Info().Msg("registerPasskey called, setting up promise handlers")

			// Handle the promise resolution
			promise.Call("then", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
				log.Info().Msg("Promise resolved successfully")
				success := args[0].Bool()
				resultChan <- success
				return nil
			})).Call("catch", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
				log.Error().Msg("JavaScript registerPasskey error")
				resultChan <- false
				return nil
			}))
			log.Info().Msg("Promise handlers attached, waiting for result")
		})

		// Wait for the result
		success := <-resultChan

		ctx.Dispatch(func(ctx app.Context) {
			s.isAddingPasskey = false

			if success {
				log.Info().Msg("Passkey registration completed successfully")
				app.Log("✅ Passkey added successfully")
				s.successMessage = "Passkey added successfully!"
				s.showSuccess = true
				ctx.Update()

				// Auto-dismiss after 3 seconds
				go func() {
					time.Sleep(3 * time.Second)
					ctx.Dispatch(func(ctx app.Context) {
						s.showSuccess = false
						ctx.Update()
					})
				}()

				// Reload passkeys - use sync version to avoid nested ctx.Async → ctx.Dispatch
				app.Log("✅ Calling loadPasskeysSync() to refresh list")
				s.loadPasskeysSync(ctx)
			} else {
				log.Error().Msg("Passkey registration failed or was cancelled")
				s.errorMessage = "Failed to add passkey. Please try again."
				s.showError = true
				ctx.Update()
			}
		})
	})
}

func (s *service) handleStartRename(index int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		s.passkeys[index].IsRenaming = true
		s.passkeys[index].EditName = s.passkeys[index].FriendlyName
	}
}

func (s *service) handleCancelRename(index int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		s.passkeys[index].IsRenaming = false
		s.passkeys[index].EditName = s.passkeys[index].FriendlyName
	}
}

func (s *service) handleEditNameInput(index int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		s.passkeys[index].EditName = ctx.JSSrc().Get("value").String()
	}
}

func (s *service) handleSaveRename(index int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		log := zerolog.Ctx(s.AppContext).With().Str("component", "PasskeyManager").Logger()
		passkey := &s.passkeys[index]

		newName := passkey.EditName
		credentialID := passkey.ID

		log.Info().Str("credentialID", credentialID).Str("newName", newName).Msg("🔵 Renaming passkey - starting")
		app.Log("🔵 Renaming passkey - credentialID:", credentialID, "newName:", newName)

		s.showError = false
		s.showSuccess = false

		ctx.Async(func() {
			app.Log("🔵 About to call RenamePasskeyHTTP API")
			response, err := s.managementApiClient.RenamePasskeyHTTP(s.AppContext,
				&models.PasskeyRenameRequest{
					CredentialID: credentialID,
					FriendlyName: newName,
				})

			app.Log("🔵 RenamePasskeyHTTP API returned - err:", err, "response:", response)

			ctx.Dispatch(func(ctx app.Context) {
				if err != nil {
					log.Error().Err(err).Msg("❌ Failed to rename passkey")
					app.Log("❌ Failed to rename passkey - error:", err)
					s.errorMessage = "Failed to rename passkey"
					s.showError = true
				} else if response != nil && response.Code == 200 {
					log.Info().Msg("✅ Passkey renamed successfully")
					app.Log("✅ Passkey renamed successfully")
					s.passkeys[index].FriendlyName = newName
					s.passkeys[index].IsRenaming = false
					s.successMessage = "Passkey renamed successfully"
					s.showSuccess = true
					ctx.Update()

					// Auto-dismiss after 3 seconds
					go func() {
						time.Sleep(3 * time.Second)
						ctx.Dispatch(func(ctx app.Context) {
							s.showSuccess = false
							ctx.Update()
						})
					}()
				} else {
					app.Log("❌ Unexpected response code:", response)
					s.errorMessage = "Failed to rename passkey"
					s.showError = true
				}
			})
		})
	}
}

func (s *service) handleDeletePasskey(credentialID string) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		s.deleteCredentialID = credentialID
		s.showDeleteConfirm = true
		ctx.Update()
	}
}

func (s *service) handleConfirmDelete(ctx app.Context, e app.Event) {
	log := zerolog.Ctx(s.AppContext).With().Str("component", "PasskeyManager").Logger()
	credentialID := s.deleteCredentialID

	log.Info().Str("credentialID", credentialID).Msg("🔴 Deleting passkey - starting")
	app.Log("🔴 Deleting passkey from modal - credentialID:", credentialID)

	// Start deleting, but keep modal open
	s.isDeleting = true
	s.showError = false
	s.showSuccess = false
	ctx.Update()

	ctx.Async(func() {
		app.Log("🔴 Calling DeletePasskeyHTTP API...")
		deleteResp, deleteErr := s.managementApiClient.DeletePasskeyHTTP(s.AppContext,
			&models.PasskeyDeleteRequest{
				CredentialID: credentialID,
			})
		app.Log("🔴 Delete complete - err:", deleteErr, "response code:", deleteResp)

		var loadResp *common.WrappedResonseT[models.PasskeysResponse]
		var loadErr error

		// If delete successful, immediately load fresh passkey list
		if deleteErr == nil && deleteResp != nil && deleteResp.Code == 200 {
			app.Log("🔴 Delete successful, loading fresh passkeys...")
			loadResp, loadErr = s.managementApiClient.GetPasskeysHTTP(s.AppContext)
			app.Log("🔴 Load passkeys complete - err:", loadErr)
		}

		ctx.Dispatch(func(ctx app.Context) {
			app.Log("🔴 INSIDE ctx.Dispatch - updating UI after delete")

			// NOW close the modal
			s.showDeleteConfirm = false
			s.isDeleting = false

			if deleteErr != nil {
				log.Error().Err(deleteErr).Msg("❌ Failed to delete passkey")
				s.errorMessage = "Failed to delete passkey"
				s.showError = true
			} else if deleteResp != nil && deleteResp.Code == 200 {
				log.Info().Msg("✅ Passkey deleted successfully")
				app.Log("✅ Passkey deleted - showing success")

				s.successMessage = "Passkey deleted successfully"
				s.showSuccess = true

				// Auto-dismiss after 3 seconds
				go func() {
					time.Sleep(3 * time.Second)
					ctx.Dispatch(func(ctx app.Context) {
						s.showSuccess = false
						ctx.Update()
					})
				}()

				// Update passkeys list from the fresh load
				if loadErr == nil && loadResp != nil && loadResp.Code == 200 {
					var passkeys []models.PasskeyItem
					if loadResp.Response != nil && loadResp.Response.Passkeys != nil {
						passkeys = loadResp.Response.Passkeys
					}

					app.Log("✅ Loaded", len(passkeys), "passkeys after delete")
					s.passkeys = make([]passkeyItem, len(passkeys))
					for i, passkey := range passkeys {
						s.passkeys[i] = passkeyItem{
							ID:           passkey.ID,
							FriendlyName: passkey.FriendlyName,
							CreatedAt:    passkey.CreatedAt,
							Transport:    passkey.Transport,
							IsRenaming:   false,
							EditName:     passkey.FriendlyName,
						}
					}
				} else {
					app.Log("❌ Failed to load passkeys after delete")
				}
			} else {
				app.Log("❌ Unexpected delete response")
				s.errorMessage = "Failed to delete passkey"
				s.showError = true
			}
			ctx.Update()
		})
	})
}

func (s *service) handleCancelDelete(ctx app.Context, e app.Event) {
	s.showDeleteConfirm = false
	s.deleteCredentialID = ""
}

func (s *service) handleCloseError(ctx app.Context, e app.Event) {
	s.showError = false
}

func (s *service) handleCloseSuccess(ctx app.Context, e app.Event) {
	s.showSuccess = false
	ctx.Update()
}

// Helper function to join transport methods
func joinTransports(transports []string) string {
	if len(transports) == 0 {
		return "Unknown"
	}
	result := transports[0]
	for i := 1; i < len(transports); i++ {
		result += ", " + transports[i]
	}
	return result
}
