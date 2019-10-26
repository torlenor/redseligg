package matrix

import (
	"net/http"
	"testing"

	"github.com/pkg/errors"
)

type mockHTTPrequest struct {
	request  http.Request
	response http.Response
	letFail  bool
}

func (mr *mockHTTPrequest) Do(request *http.Request) (*http.Response, error) {
	mr.request = *request
	if mr.letFail {
		return nil, errors.New("Failed")
	}
	return &mr.response, nil
}

func Test_matrixAPI_call(t *testing.T) {
	mr := &mockHTTPrequest{}
	doHTTPRequest = mr.Do

	// TODO: Add tests
}
