package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateMatrixBot(t *testing.T) {
	assert := assert.New(t)

	result := StripCmd("!CMD test", "CMD")
	assert.Equal("test", result)

	result = StripCmd("test !CMD", "CMD")
	assert.Equal("test !CMD", result)

	result = StripCmd("!CMDtest", "CMD")
	assert.Equal("!CMDtest", result)

	result = StripCmd("!TEST test", "CMD")
	assert.Equal("!TEST test", result)

	result = StripCmd("!CMD2 test", "CMD")
	assert.Equal("!CMD2 test", result)
}