package config

import (
	"goverwatch/helpers"
	"log"
	"path/filepath"
	"testing"
)

func TestConfigurationParsing(t *testing.T) {
	t.Run("test yaml configuration", func(t *testing.T) {
		absolutePath, err := filepath.Abs("./samples/test-config.yaml")
		if err != nil {
			log.Fatal(err)
		}
		testConfig(absolutePath, t)
	})

	t.Run("test json configuration", func(t *testing.T) {
		absolutePath, err := filepath.Abs("./samples/test-config.json")
		if err != nil {
			log.Fatal(err)
		}
		testConfig(absolutePath, t)
	})
}

func testConfig(absolutePath string, t *testing.T) {
	conf := ParseConfigFile(absolutePath)
	patternsToGlob := conf.Files
	patternsToIgnore := conf.IgnoredFiles
	command := conf.Command

	wantedFiles := []string{"*.go", "**/**.go", "**/**.go", "**/**/**.go"}
	wantedCommand := "go build"
	wantedIgnoredFiles := []string{"**/a*.go", "**/**/b*.go"}
	if command != wantedCommand {
		t.Errorf("issue with command got %s want %s", command, wantedCommand)
	}
	if !helpers.Equal(patternsToGlob, wantedFiles) {
		t.Errorf("issue with files got %s want %s", patternsToGlob, wantedFiles)
	}
	if !helpers.Equal(patternsToIgnore, wantedIgnoredFiles) {
		t.Errorf("issue with ignored files got %s want %s", patternsToIgnore, wantedIgnoredFiles)
	}
}
