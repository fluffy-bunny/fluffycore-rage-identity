/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package serve

import (
	"fmt"
	"strings"

	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	shared "github.com/fluffy-bunny/fluffycore-rage-oidc/cmd/oidc-client/shared"
	cobra_utils "github.com/fluffy-bunny/fluffycore-rage-oidc/internal/cobra_utils"
	fluffycore_utils "github.com/fluffy-bunny/fluffycore/utils"
	req "github.com/imroc/req/v3"
	zerolog "github.com/rs/zerolog"
	cobra "github.com/spf13/cobra"
	context "golang.org/x/net/context"
	oauth2 "golang.org/x/oauth2"
)

// serveCmd represents the about command
var serveCmd = &cobra.Command{
	Use:               "serve",
	PersistentPreRunE: cobra_utils.ParentPersistentPreRunE,
	Short:             "serves the client server",
	Long:              `serves the client server`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(fluffycore_utils.PrettyJSON(shared.AppConfig))
		Serve()
	},
}
var (
	callbackPath = "/auth/callback"
)

func InitCommand(parent *cobra.Command) {
	serveCmd.PersistentFlags().StringVar(&shared.AppConfig.Authority, "authority", "", "the authority of the provider i.e. https://accounts.google.com")
	serveCmd.PersistentFlags().StringVar(&shared.AppConfig.ClientId, "client_id", "", "the client id")
	serveCmd.PersistentFlags().StringVar(&shared.AppConfig.ClientSecret, "client_secret", "", "the client secret")
	serveCmd.PersistentFlags().IntVar(&shared.AppConfig.Port, "port", 5556, "the port to listen on")
	serveCmd.PersistentFlags().StringArrayVar(&shared.AppConfig.ACRValues, "acr_values", []string{}, "the acr_values")
	parent.AddCommand(serveCmd)
}
func randString(nByte int) (string, error) {
	b := make([]byte, nByte)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(b), nil
}
func setCallbackCookie(w http.ResponseWriter, r *http.Request, name, value string) {
	c := &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   int(time.Hour.Seconds()),
		Secure:   r.TLS != nil,
		HttpOnly: true,
	}
	http.SetCookie(w, c)
}

func Serve() {
	ctx := context.Background()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// create a logger and add it to the context
	logz := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()

	ctx = logz.WithContext(ctx)
	log := zerolog.Ctx(ctx)

	provider, err := oidc.NewProvider(ctx, shared.AppConfig.Authority)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to query provider.")
	}
	tokenUrl := provider.Endpoint().TokenURL

	oidcConfig := &oidc.Config{
		ClientID: shared.AppConfig.ClientId,
	}
	verifier := provider.Verifier(oidcConfig)

	config := oauth2.Config{
		ClientID:     shared.AppConfig.ClientId,
		ClientSecret: shared.AppConfig.ClientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  fmt.Sprintf("http://localhost:%d%s", shared.AppConfig.Port, callbackPath),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", oidc.ScopeOfflineAccess},
	}

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		state, err := randString(16)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		nonce, err := randString(16)
		if err != nil {
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return
		}
		setCallbackCookie(w, r, "state", state)
		setCallbackCookie(w, r, "nonce", nonce)
		authCodeOptions := []oauth2.AuthCodeOption{
			oidc.Nonce(nonce),
		}
		if len(shared.AppConfig.ACRValues) > 0 {
			authCodeOptions = append(authCodeOptions, AcrValues(shared.AppConfig.ACRValues...))
		}
		authRequestURL := config.AuthCodeURL(state, authCodeOptions...)
		fmt.Println(authRequestURL)
		http.Redirect(w, r, authRequestURL, http.StatusFound)
	})

	http.HandleFunc(callbackPath, func(w http.ResponseWriter, r *http.Request) {
		state, err := r.Cookie("state")
		if err != nil {
			http.Error(w, "state not found", http.StatusBadRequest)
			return
		}
		if r.URL.Query().Get("state") != state.Value {
			http.Error(w, "state did not match", http.StatusBadRequest)
			return
		}

		oauth2Token, err := config.Exchange(ctx, r.URL.Query().Get("code"))
		if err != nil {
			http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		log.Info().Interface("oauth2Token", oauth2Token).Msg("callbackPath")
		rawIDToken, ok := oauth2Token.Extra("id_token").(string)
		if !ok {
			http.Error(w, "No id_token field in oauth2 token.", http.StatusInternalServerError)
			return
		}
		idToken, err := verifier.Verify(ctx, rawIDToken)
		if err != nil {
			http.Error(w, "Failed to verify ID Token: "+err.Error(), http.StatusInternalServerError)
			return
		}
		nonce, err := r.Cookie("nonce")
		if err != nil {
			http.Error(w, "nonce not found", http.StatusBadRequest)
			return
		}
		if idToken.Nonce != nonce.Value {
			http.Error(w, "nonce did not match", http.StatusBadRequest)
			return
		}

		//oauth2Token.AccessToken = "*REDACTED*"

		reqClient := req.C()
		resp2, err := reqClient.R().
			SetBasicAuth(shared.AppConfig.ClientId, shared.AppConfig.ClientSecret).
			SetFormData(map[string]string{
				"grant_type":    "refresh_token",
				"refresh_token": oauth2Token.RefreshToken,
			}).Post(tokenUrl)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		generic := make(map[string]interface{})
		json.Unmarshal(resp2.Bytes(), &generic)

		log.Info().Interface("generic", generic).Msg("")

		resp := struct {
			OAuth2Token   *oauth2.Token
			IDTokenClaims *json.RawMessage // ID Token payload is just JSON.
			IDToken       string
		}{oauth2Token, new(json.RawMessage), rawIDToken}

		if err := idToken.Claims(&resp.IDTokenClaims); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := json.MarshalIndent(resp, "", "    ")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(data)
	})

	addr := fmt.Sprintf("localhost:%d", shared.AppConfig.Port)
	log.Info().Msgf("listening on http://%s/", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}
}
func AcrValues(acr ...string) oauth2.AuthCodeOption {
	acrValues := strings.Join(acr, " ")
	return oauth2.SetAuthURLParam("acr_values", acrValues)
}
