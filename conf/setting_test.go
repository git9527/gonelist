package conf

import (
	"fmt"
	"gopkg.in/go-playground/assert.v1"
	"testing"
)

func TestLoadUserConfig(t *testing.T) {
	filePath := "../config.json"
	_ = LoadUserConfig(filePath)
	fmt.Println(UserSet)
	fmt.Println(*UserSet.Server)
	fmt.Println(UserSet.PassList[0])
	assert.Equal(t, UserSet.DomainBasedSubFolders.Pairs[1].Domain, "127.0.0.1:8000")
}

func TestGetBindAddr(t *testing.T) {
	fmt.Println(GetBindAddr(false, 8000))
}

func TestGetTokenPath(t *testing.T) {
	var (
		configPath string
		tokenPath  string
	)
	configPath = "/etc/gonelist/config.json"
	tokenPath = GetTokenPath(configPath)
	assert.Equal(t, tokenPath, "/etc/gonelist/.token")

	configPath = "config.json"
	tokenPath = GetTokenPath(configPath)
	assert.Equal(t, tokenPath, ".token")
}
