package oauth

// Profile is a profile with google.
type Profile struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Link          string `json:"link"`
	Gender        string `json:"gender"`
	Locale        string `json:"locale"`
	PictureURL    string `json:"picture"`
}

// Username returns the <username>@fqdn component
// of the email address.
func (p Profile) Username() string {
	return p.Email
}
