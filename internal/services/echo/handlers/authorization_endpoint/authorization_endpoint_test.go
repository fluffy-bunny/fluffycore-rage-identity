package authorization_endpoint

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIdpUrnExtraction(t *testing.T) {

	template := "urn:mastodon:idp:google"
	info, err := extractIdpSlug(template)
	require.NoError(t, err)
	require.Equal(t, "google", info["idp_slug"])

	template = "urn:mastodon:idp:"
	info, err = extractIdpSlug(template)
	require.NoError(t, err)
	require.Equal(t, "", info["idp_slug"])

	template = "urn:mastodon:idp"
	info, err = extractIdpSlug(template)
	require.Error(t, err)
	require.Empty(t, info)

	template = "invalid_template"
	info, err = extractIdpSlug(template)
	require.Error(t, err)
	require.Empty(t, info)

}

func TestRootCandidateUrnExtraction(t *testing.T) {

	template := "urn:mastodon:root_candidate:123"
	info, err := extractRootCandidate(template)
	require.NoError(t, err)
	require.Equal(t, "123", info["user_id"])

	template = "urn:mastodon:root_candidate:"
	info, err = extractRootCandidate(template)
	require.NoError(t, err)
	require.Equal(t, "", info["user_id"])

	template = "gargabe_template"
	info, err = extractRootCandidate(template)
	require.Error(t, err)
	require.Empty(t, info)

	template = "urn:mastodon:idp:google"
	info, err = extractRootCandidate(template)
	require.Error(t, err)
	require.Empty(t, info)

	template = "urn:mastodon:root_candidate"
	info, err = extractRootCandidate(template)
	require.Error(t, err)
	require.Empty(t, info)

}
