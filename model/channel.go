package model

// Channel describes a certain channel of the platform (can be private or public).
// Not all fields are filled by every platform, but ID and Name will shall always be there
// and the platforms should make sure they can uniquely identify a channel by it.
type Channel struct {
	ID   string `json:"id"`
	Name string `json:"name"`

	Topic   string `json:"header"`
	Purpose string `json:"purpose"`

	IsPrivate bool `json:"is_private"`
}
