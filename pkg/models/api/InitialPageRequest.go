package api

type (
	IDP struct {
		Name string `json:"name"`
		Slug string `json:"slug"`
	}
	InitialPageRequest struct {
		BaseRequest
	}
	InitialPageResponse struct {
		BaseResponse
		IDPs []IDP `json:"idps"`
	}
)
