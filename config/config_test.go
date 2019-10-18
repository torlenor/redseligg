package config

import (
	"log"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var expectedConfig = Config{
	General: General{
		API: API{
			Enabled: true,
			IP:      "127.0.0.1",
			Port:    "8000",
		},
	},
	Bots: BotsConfig{
		Discord: DiscordConfig{
			Token:  "INSERT_DISCORD_TOKEN_HERE",
			ID:     "INSERT_DISCORD_CLIENT_ID_HERE",
			Secret: "INSERT_DISCORD_SECRET_HERE",
		},
		Matrix: MatrixConfig{
			Enabled:  true,
			Server:   "INSERT_MATRIX_SERVER_HERE",
			Username: "INSERT_MATRIX_USER_HERE",
			Password: "INSERT_MATRIX_PASSWORD_HERE",
			Token:    "INSERT_MATRIX_TOKEN_HERE",
		},
		Mattermost: MattermostConfig{
			Server:   "INSERT_MATTERMOST_SERVER_HERE",
			Username: "INSERT_MATTERMOST_USER_HERE",
			Password: "INSERT_MATTERMOST_PASSWORT_HERE",
			UseToken: true,
			Token:    "TOKEN",
		},
		Slack: SlackConfig{
			Workspace: "INSERT_SLACK_WORKSPACE_HERE",
			Token:     "INSERT_BOT_TOKEN_HERE",
		},
	},
}

func TestConfig_ParseFromFile(t *testing.T) {
	assert := assert.New(t)

	_, filename, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		log.Fatal(err)
	}

	actualConfig, err := ParseFromFile(dir + "/../test/testdata/config.toml")
	assert.NoError(err)

	assert.Equal(expectedConfig, actualConfig)

}

func TestConfig_ParseFromFile_File_Corrupt(t *testing.T) {
	assert := assert.New(t)

	_, filename, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		log.Fatal(err)
	}

	_, err = ParseFromFile(dir + "/../test/testdata/config_corrupt.toml")
	assert.Error(err)
}

func TestConfig_ParseFromFile_File_Does_Not_Exist(t *testing.T) {
	assert := assert.New(t)

	_, filename, _, _ := runtime.Caller(0)
	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		log.Fatal(err)
	}

	_, err = ParseFromFile(dir + "/../test/testdata/config_does_not_exist.toml")
	assert.Error(err)
}
