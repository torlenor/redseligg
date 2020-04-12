package model

// Reaction is a reaction on a message, e.g., with an emoji.
type Reaction struct {
	Message MessageIdentifier // Uniquely identifies the message on which the reaction has been made

	Type     string // "added" or "removed"
	Reaction string // e.g., "wink"
	User     User   // The User who made the reaction to the Message
}
