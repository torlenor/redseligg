package webclient

import "net/http"

// APIResponse contains all the data from a successful API call
type APIResponse struct {
	Header     http.Header
	Body       []byte
	StatusCode int
}
