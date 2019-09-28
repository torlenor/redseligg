package mattermost

// ErrorResponse is a response from Mattermost API if something went wrong
type ErrorResponse struct {
	ID         string `json:"id"`
	Message    string `json:"message"`
	RequestID  string `json:"request_id"`
	StatusCode int    `json:"status_code"`
	IsOauth    bool   `json:"is_oauth"`
}

// UserObject contains informations about a Mattermost User
type UserObject struct {
	ID             string `json:"id"`
	CreateAt       int64  `json:"create_at"`
	UpdateAt       int64  `json:"update_at"`
	DeleteAt       int    `json:"delete_at"`
	Username       string `json:"username"`
	AuthData       string `json:"auth_data"`
	AuthService    string `json:"auth_service"`
	Email          string `json:"email"`
	EmailVerified  bool   `json:"email_verified"`
	Nickname       string `json:"nickname"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Position       string `json:"position"`
	Roles          string `json:"roles"`
	AllowMarketing bool   `json:"allow_marketing"`
	NotifyProps    struct {
		Channel      string `json:"channel"`
		Comments     string `json:"comments"`
		Desktop      string `json:"desktop"`
		DesktopSound string `json:"desktop_sound"`
		Email        string `json:"email"`
		FirstName    string `json:"first_name"`
		MentionKeys  string `json:"mention_keys"`
		Push         string `json:"push"`
		PushStatus   string `json:"push_status"`
	} `json:"notify_props"`
	LastPasswordUpdate int64  `json:"last_password_update"`
	LastPictureUpdate  int64  `json:"last_picture_update"`
	FailedAttempts     int    `json:"failed_attempts"`
	Locale             string `json:"locale"`
	Timezone           struct {
		AutomaticTimezone    string `json:"automaticTimezone"`
		ManualTimezone       string `json:"manualTimezone"`
		UseAutomaticTimezone string `json:"useAutomaticTimezone"`
	} `json:"timezone"`
}
