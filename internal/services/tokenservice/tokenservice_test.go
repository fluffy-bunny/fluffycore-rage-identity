package tokenservice

import (
	"context"
	"encoding/json"
	"strings"
	"testing"
	"time"

	di "github.com/fluffy-bunny/fluffy-dozm-di"
	fluffycore_contracts_jwtminter "github.com/fluffy-bunny/fluffycore/contracts/jwtminter"
	fluffycore_services_claims "github.com/fluffy-bunny/fluffycore/services/claims"
	fluffycore_services_jwtminter "github.com/fluffy-bunny/fluffycore/services/jwtminter"
	fluffycore_services_keymaterial "github.com/fluffy-bunny/fluffycore/services/keymaterial"

	require "github.com/stretchr/testify/require"
)

const signingKeysTemplate = `{
    "signing_keys": [
        {
            "private_key": "-----BEGIN EC PRIVATE KEY-----\nMHcCAQEEIFA+8y3M5qxkjuI7HOUAPVgrsjUnu5kwRjsZlbCmyabCoAoGCCqGSM49\nAwEHoUQDQgAEYMrUm/S5+d+euQHrrzQMWJSFafSYcgUE0RYjfI7sErK75hSdE0aj\nPNMXaaDG395zD18VxjsmqPTWom17ncVnnw==\n-----END EC PRIVATE KEY-----\n",
            "public_key": "-----BEGIN EC  PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEYMrUm/S5+d+euQHrrzQMWJSFafSY\ncgUE0RYjfI7sErK75hSdE0ajPNMXaaDG395zD18VxjsmqPTWom17ncVnnw==\n-----END EC  PUBLIC KEY-----\n",
            "not_before": "{not_before}",
            "not_after": "{not_after}",
            "password": "",
            "kid": "0b2cd2e54c924ce89f010f242862367d",
            "public_jwk": {
                "alg": "ES256",
                "crv": "P-256",
                "kid": "0b2cd2e54c924ce89f010f242862367d",
                "kty": "EC",
                "use": "sig",
                "x": "YMrUm_S5-d-euQHrrzQMWJSFafSYcgUE0RYjfI7sErI",
                "y": "u-YUnRNGozzTF2mgxt_ecw9fFcY7Jqj01qJte53FZ58"
            },
            "private_jwk": {
                "alg": "ES256",
                "crv": "P-256",
                "d": "UD7zLczmrGSO4jsc5QA9WCuyNSe7mTBGOxmVsKbJpsI",
                "kid": "0b2cd2e54c924ce89f010f242862367d",
                "kty": "EC",
                "use": "sig",
                "x": "YMrUm_S5-d-euQHrrzQMWJSFafSYcgUE0RYjfI7sErI",
                "y": "u-YUnRNGozzTF2mgxt_ecw9fFcY7Jqj01qJte53FZ58"
            }
        }
    ]
}`

var signingKeys = ""

func getSigningKeysJSON() string {
	now := time.Now()
	notBefore := now.Add(-1 * time.Hour)
	notAfter := now.Add(24 * time.Hour)

	nbf := notBefore.Format("2006-01-02T15:04:05Z")
	naf := notAfter.Format("2006-01-02T15:04:05Z")
	signingKeys = strings.Replace(signingKeysTemplate, "{not_before}", nbf, -1)
	signingKeys = strings.Replace(signingKeys, "{not_after}", naf, -1)
	return signingKeys

}
func TestMintToken(t *testing.T) {

	b := di.Builder()
	b.ConfigureOptions(func(o *di.Options) {
		o.ValidateScopes = true
		o.ValidateOnBuild = true
	})
	keyMaterialJSON := getSigningKeysJSON()
	keymaterial := &fluffycore_contracts_jwtminter.KeyMaterial{}
	err := json.Unmarshal([]byte(keyMaterialJSON), keymaterial)
	require.NoError(t, err)
	di.AddInstance[*fluffycore_contracts_jwtminter.KeyMaterial](b, keymaterial)
	// order maters for Singleton and Transient, they are both app scoped and the last one wins
	fluffycore_services_jwtminter.AddSingletonIJWTMinter(b)
	fluffycore_services_keymaterial.AddSingletonIKeyMaterial(b)
	container := b.Build()

	jwtMinter := di.Get[fluffycore_contracts_jwtminter.IJWTMinter](container)
	require.NotNil(t, jwtMinter)

	now := time.Now()
	expiration := now.Add(24 * time.Hour).Unix()
	claims := fluffycore_services_claims.NewClaims()
	claims.Set("sub", "1234567890")
	claims.Set("name", "John Doe")
	claims.Set("iss", "http://example.com")
	claims.Set("exp", expiration)

	token, err := jwtMinter.MintToken(context.TODO(), claims)
	require.NoError(t, err)
	require.NotEmpty(t, token)

}
