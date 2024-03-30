package api

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitialPageRequest(t *testing.T) {
	request := InitialPageRequest{
		BaseRequest: BaseRequest{
			RequestType: "InitialPageRequest",
			Version:     "1",
		},
	}
	requestJson, _ := json.Marshal(request)
	t.Log(string(requestJson))

	// marshal as base to then determine the type
	var base BaseRequest
	err := json.Unmarshal(requestJson, &base)
	require.NoError(t, err)
	require.Equal(t, "InitialPageRequest", base.RequestType)
	require.Equal(t, "1", base.Version)

	finalRequest := InitialPageRequest{}
	err = json.Unmarshal(requestJson, &finalRequest)
	require.NoError(t, err)
	require.Equal(t, "InitialPageRequest", finalRequest.RequestType)
	require.Equal(t, "1", finalRequest.Version)
}
