package api

import (
	"net/http"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/torlenor/redseligg/config"
)

type mockRouter struct {
	routes []*mux.Route
	paths  []string
}

func (m *mockRouter) HandleFunc(path string, f func(http.ResponseWriter, *http.Request)) *mux.Route {
	r := &mux.Route{}
	m.routes = append(m.routes, r)
	m.paths = append(m.paths, path)
	return r
}
func (m *mockRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}

func fakeFunk(http.ResponseWriter, *http.Request) {
	return
}

func TestCreatingNewAPI(t *testing.T) {
	assert := assert.New(t)

	config := config.API{}
	api, err := NewAPI(config, "")
	assert.Error(err)

	config.Port = "1234"
	api, err = NewAPI(config, "")
	assert.NoError(err)
	assert.NotNil(api)
}

func TestInit(t *testing.T) {
	assert := assert.New(t)
	config := config.API{
		Port: "1234",
	}

	api, err := NewAPI(config, "")
	assert.NoError(err)
	assert.NotNil(api)
	err = api.Init()
	assert.NoError(err)
	assert.Equal(":1234", api.server.Addr)

	config.IP = "0.0.0.0"
	api, err = NewAPI(config, "")
	assert.NoError(err)
	assert.NotNil(api)
	err = api.Init()
	assert.NoError(err)
	assert.Equal("0.0.0.0:1234", api.server.Addr)
}

func TestAttachModule(t *testing.T) {
	assert := assert.New(t)
	mockRouter := &mockRouter{}

	config := config.API{}
	api, err := NewAPICustom(config, "", mockRouter)
	assert.Error(err)

	config.Port = "1234"

	api, err = NewAPICustom(config, "", mockRouter)
	assert.NoError(err)
	assert.NotNil(api)

	assert.Equal(0, len(mockRouter.routes))
	assert.Equal(0, len(mockRouter.paths))

	cnt := 0
	api.AttachModuleGet("/get", fakeFunk)
	assert.Equal(cnt+1, len(mockRouter.routes))
	assert.Equal(cnt+1, len(mockRouter.paths))
	assert.Equal("/get", mockRouter.paths[0])
	methods, err := mockRouter.routes[cnt].GetMethods()
	assert.Equal(1, len(methods))
	assert.Equal("GET", methods[0])

	cnt++
	api.AttachModulePost("/post", fakeFunk)
	assert.Equal(cnt+1, len(mockRouter.routes))
	assert.Equal(cnt+1, len(mockRouter.paths))
	assert.Equal("/post", mockRouter.paths[1])
	methods, err = mockRouter.routes[cnt].GetMethods()
	assert.Equal(1, len(methods))
	assert.Equal("POST", methods[0])

	cnt++
	api.AttachModulePut("/put", fakeFunk)
	assert.Equal(cnt+1, len(mockRouter.routes))
	assert.Equal(cnt+1, len(mockRouter.paths))
	assert.Equal("/put", mockRouter.paths[cnt])
	methods, err = mockRouter.routes[cnt].GetMethods()
	assert.Equal(1, len(methods))
	assert.Equal("PUT", methods[0])

	cnt++
	api.AttachModuleDelete("/delete", fakeFunk)
	assert.Equal(cnt+1, len(mockRouter.routes))
	assert.Equal(cnt+1, len(mockRouter.paths))
	assert.Equal("/delete", mockRouter.paths[cnt])
	methods, err = mockRouter.routes[cnt].GetMethods()
	assert.Equal(1, len(methods))
	assert.Equal("DELETE", methods[0])

}
