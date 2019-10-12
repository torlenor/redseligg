package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIDProvider_Get(t *testing.T) {
	assert := assert.New(t)

	idProvider := IDProvider{}

	assert.Equal(1, idProvider.Get())
	assert.Equal(2, idProvider.Get())
	assert.Equal(3, idProvider.Get())
}
