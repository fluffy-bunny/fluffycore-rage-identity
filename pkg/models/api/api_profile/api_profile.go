package api_profile

type Profile struct {
	Subject         string `json:"subject"`
	Email           string `json:"email"`
	EmailVerified   bool   `json:"emailVerified"`
	GivenName       string `json:"givenName"`
	FamilyName      string `json:"familyName"`
	PhoneNumber     string `json:"phoneNumber"`
	IsClaimedDomain bool   `json:"isClaimedDomain"`
}
