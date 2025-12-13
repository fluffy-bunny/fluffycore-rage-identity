package LocalizerBundle

import (
	"testing"

	contracts_LocalizerBundle "github.com/fluffy-bunny/fluffycore-rage-identity/example/go-app/management/internal/contracts/LocalizerBundle"
	i18n "github.com/nicksnyder/go-i18n/v2/i18n"
	require "github.com/stretchr/testify/require"
)

func TestLocalizerBundle_Ctor(t *testing.T) {
	svc, err := stemService.Ctor()
	require.NoError(t, err)
	require.NotNil(t, svc)

	bundle := svc.GetBundle()
	require.NotNil(t, bundle)

	localizer := i18n.NewLocalizer(bundle)
	require.NotNil(t, localizer)

	for _, localKey := range contracts_LocalizerBundle.AllLocaleKeys {
		ss, err := localizer.LocalizeMessage(&i18n.Message{
			ID: string(localKey),
		})
		require.NoError(t, err)
		require.NotEmpty(t, ss)
	}

}
