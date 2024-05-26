package manifest

type (
	IDP struct {
		Slug        string `json:"slug"`
		 
	}
	Manifest struct {
		SocialIdps []IDP `json:"social_idps"`
	}
)
