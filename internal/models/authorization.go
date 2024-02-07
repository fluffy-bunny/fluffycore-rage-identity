package models

type (
	AuthorizationRequest struct {
		ResponseType        string `param:"response_type" query:"response_type" form:"response_type" json:"response_type" xml:"response_type"`
		Scope               string `param:"scope" query:"scope" form:"scope" json:"scope" xml:"scope"`
		State               string `param:"state" query:"state" form:"state" json:"state" xml:"state"`
		RedirectURI         string `param:"redirect_uri" query:"redirect_uri" form:"redirect_uri" json:"redirect_uri" xml:"redirect_uri"`
		Audience            string `param:"audience" query:"audience" form:"audience" json:"audience" xml:"audience"`
		CodeChallenge       string `param:"code_challenge" query:"code_challenge" form:"code_challenge" json:"code_challenge" xml:"code_challenge"`
		CodeChallengeMethod string `param:"code_challenge_method" query:"code_challenge_method" form:"code_challenge_method" json:"code_challenge_method" xml:"code_challenge_method"`
	}
)
