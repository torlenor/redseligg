package storage

import "errors"

// ErrNotFound is returned when no data could be found
var ErrNotFound = errors.New("Could not find requested data")
