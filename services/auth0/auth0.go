package auth0

import "time"

type Identity struct {
	Connection   string `json:"connection,omitempty"`
	UserID       string `json:"user_id,omitempty"`
	Provider     string `json:"provider,omitempty"`
	IsSocial     bool   `json:"isSocial,omitempty"`
	ExpiresIn    int    `json:"expires_in,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	AccessToken  string `json:"access_token,omitempty"`
}

type User struct {

	// The users identifier.
	ID string `json:"user_id,omitempty"`

	// The connection the user belongs to.
	Connection string `json:"connection,omitempty"`

	// The user's email
	Email string `json:"email,omitempty"`

	// The users name
	Name string `json:"name,omitempty"`

	// The users given name
	GivenName string `json:"given_name,omitempty"`

	// The users family name
	FamilyName string `json:"family_name,omitempty"`

	// The user's username. Only valid if the connection requires a username
	Username string `json:"username,omitempty"`

	// The user's nickname
	Nickname string `json:"nickname,omitempty"`

	// The user's password (mandatory for non SMS connections)
	Password string `json:"password,omitempty"`

	// The user's phone number (following the E.164 recommendation), only valid
	// for users to be added to SMS connections.
	PhoneNumber string `json:"phone_number,omitempty"`

	// The time the user is created.
	CreatedAt time.Time `json:"created_at,omitempty"`

	// The last time the user is updated.
	UpdatedAt time.Time `json:"updated_at,omitempty"`

	// The last time the user has logged in.
	LastLogin time.Time `json:"last_login,omitempty"`

	// UserMetadata holds data that the user has read/write access to (e.g.
	// color_preference, blog_url, etc).
	UserMetadata map[string]interface{} `json:"user_metadata,omitempty"`

	Identities []Identity `json:"identities,omitempty"`

	// True if the user's email is verified, false otherwise. If it is true then
	// the user will not receive a verification email, unless verify_email: true
	// was specified.
	EmailVerified bool `json:"email_verified,omitempty"`

	// If true, the user will receive a verification email after creation, even
	// if created with email_verified set to true. If false, the user will not
	// receive a verification email, even if created with email_verified set to
	// false. If unspecified, defaults to the behavior determined by the value
	// of email_verified.
	VerifyEmail bool `json:"verify_email,omitempty"`

	// True if the user's phone number is verified, false otherwise. When the
	// user is added to a SMS connection, they will not receive an verification
	// SMS if this is true.
	PhoneVerified bool `json:"phone_verified,omitempty"`

	// AppMetadata holds data that the user has read-only access to (e.g. roles,
	// permissions, vip, etc).
	AppMetadata map[string]interface{} `json:"app_metadata,omitempty"`

	// The user's picture url
	Picture string `json:"picture,omitempty"`
}
