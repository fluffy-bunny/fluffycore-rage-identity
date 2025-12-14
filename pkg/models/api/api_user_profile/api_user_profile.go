package api_user_profile

type (
	UserProfile struct {
		Email       string `json:"email"`
		GivenName   string `json:"givenName"`
		FamilyName  string `json:"familyName"`
		PhoneNumber string `json:"phoneNumber"`
	}
)
