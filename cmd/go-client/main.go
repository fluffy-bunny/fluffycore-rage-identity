/*
This is an example application to demonstrate parsing an ID Token.
*/
package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	oidc "github.com/coreos/go-oidc/v3/oidc"
	req "github.com/imroc/req/v3"
	zerolog "github.com/rs/zerolog"
	"golang.org/x/net/context"
	oauth2 "golang.org/x/oauth2"
)

var (
	clientID     = os.Getenv("OAUTH2_CLIENT_ID")
	clientSecret = os.Getenv("OAUTH2_CLIENT_SECRET")
	port         = os.Getenv("PORT")
	authority    = os.Getenv("AUTHORITY")
	callbackPath = "/auth/callback"
)

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
		HttpOnly: false,
	}
	http.SetCookie(w, c)
}

func main() {
	ctx := context.Background()
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	// create a logger and add it to the context
	logz := zerolog.New(os.Stdout).With().Caller().Timestamp().Logger()

	ctx = logz.WithContext(ctx)
	log := zerolog.Ctx(ctx)

	provider, err := oidc.NewProvider(ctx, authority)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to query provider.")
	}
	tokenUrl := provider.Endpoint().TokenURL

	oidcConfig := &oidc.Config{
		ClientID: clientID,
	}
	verifier := provider.Verifier(oidcConfig)

	config := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		RedirectURL:  fmt.Sprintf("http://localhost:%s%s", port, callbackPath),
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

		http.Redirect(w, r, config.AuthCodeURL(state, oidc.Nonce(nonce)), http.StatusFound)
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
			SetBasicAuth(clientID, clientSecret).
			SetFormData(map[string]string{
				"grant_type":    "refresh_token",
				"refresh_token": oauth2Token.RefreshToken,
			}).Post(tokenUrl)
		if err != nil {
			log.Fatal().Err(err).Msg("")
		}
		generic := make(map[string]interface{})
		err = json.Unmarshal(resp2.Bytes(), &generic)
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

	addr := fmt.Sprintf("localhost:%s", port)
	log.Info().Msgf("listening on http://%s/", addr)
	err = http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal().Err(err).Msg("ListenAndServe")
	}
}
