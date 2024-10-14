package api

type (
	BaseRequest struct {
		RequestType string `param:"request_type" query:"request_type" form:"request_type" json:"request_type" xml:"request_type"`
		Version     string `param:"version" query:"version" form:"version" json:"version" xml:"version"`
	}
	BaseResponse struct {
		Errors []string `json:"errors,omitempty"`
	}
	UnautorizedResponse struct {
		Path string `json:"path"`
	}
	AuthorizedResponse struct {
		Path string `json:"path"`
	}
)
