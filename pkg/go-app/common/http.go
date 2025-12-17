package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	fluffycore_go_app_cookies "github.com/fluffy-bunny/fluffycore/go-app/cookies"
	zerolog "github.com/rs/zerolog"
)

// GetCSRFToken retrieves the CSRF token from the _csrf cookie
func GetCSRFToken() string {
	csrfToken, err := fluffycore_go_app_cookies.GetCookie[string]("_csrf")
	if err != nil {
		return ""
	}
	return csrfToken
}

// BuildCustomHeaders creates custom headers map with CSRF token if it exists
func BuildCustomHeaders() map[string]string {
	csrfToken := GetCSRFToken()
	if csrfToken != "" {
		return map[string]string{
			"X-Csrf-Token": csrfToken,
		}
	}
	return nil
}

type CallInput struct {
	Method        string
	Url           string
	Data          any
	CustomHeaders map[string]string
}
type WrappedResonseT[T any] struct {
	Response *T  `json:"response"`
	Code     int `json:"code"`
}

func HTTPFetchWrappedResponseT[T any](ctx context.Context, input *CallInput) (*WrappedResonseT[T], error) {
	responseData, code, err := HTTPDataT[T](ctx, input)
	if err != nil {
		return nil, err
	}
	return &WrappedResonseT[T]{
		Response: responseData,
		Code:     code,
	}, nil
}

func HTTPDataT[T any](ctx context.Context, input *CallInput) (*T, int, error) {
	log := zerolog.Ctx(ctx).With().
		Str("component", "HTTPDataT").
		Interface("input", input).
		Logger()
	req, err := http.NewRequest(input.Method, input.Url, nil)
	if err != nil {
		log.Error().Err(err).Msg("HTTPDataT: NewRequest error")
		return nil, 0, fmt.Errorf("failed to create request: %w", err)
	}

	for key, value := range input.CustomHeaders {
		req.Header.Set(key, value)
	}
	// http.DefaultClient in WASM automatically includes cookies (credentials: 'include')
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Error().Err(err).Msg("HTTPDataT: request error")
		return nil, 0, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("failed to read response: %w", err)
	}

	var result T
	// create a reader for the body
	r := bytes.NewReader(body)
	if err := json.NewDecoder(r).Decode(&result); err != nil {
		log.Error().Err(err).Msg("FetchDataT: decode error")
		return nil, resp.StatusCode, fmt.Errorf("failed to decode response: %w", err)
	}
	log.Debug().Interface("result", result).Msg("FetchDataT: decoded result")

	return &result, resp.StatusCode, nil
}
