package model

// User holds all relevant information about a User.
// Not all fields are filled by every platform, but ID and Name will shall always be there
// and the platforms should make sure they can uniquely identify a user by it.
type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Nickname  string `json:"nickname"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`

	Tz     string `json:"tz"`
	Locale string `json:"locale"`

	IsAdmin bool `json:"is_admin"`
	IsBot   bool `json:"is_bot"`
	IsOwner bool `json:"is_owner"`

	// IsMod is a Redseligg property indicating certain rights of that user inside of Redseligg and its plugins
	IsMod bool `json:"is_mod"`
}

// IsValid indicates if a User object is valid
func (u User) IsValid() bool {
	if len(u.ID) == 0 || len(u.Name) == 0 {
		return false
	}

	return true
}
