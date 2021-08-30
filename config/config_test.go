package config

import (
	"errors"
	"goverwatch/helpers"
	"log"
	"path/filepath"
	"testing"
)

func TestConfigurationParsing(t *testing.T) {
	t.Run("test yaml configuration", func(t *testing.T) {
		absolutePath, err := filepath.Abs("./_samples/test-config.yaml")
		if err != nil {
			log.Fatal(err)
		}

		ext, err := getConfigFileExtension(absolutePath)
		if errors.Is(err, extensionErr) && ext != ".yaml" {
			t.Errorf("extension should be yaml")
		}

		testConfig(absolutePath, t)
	})

	t.Run("test json configuration", func(t *testing.T) {
		absolutePath, err := filepath.Abs("./_samples/test-config.json")
		if err != nil {
			log.Fatal(err)
		}

		ext, err := getConfigFileExtension(absolutePath)
		if errors.Is(err, extensionErr) && ext != ".json" {
			t.Errorf("extension should be json")
		}

		testConfig(absolutePath, t)
	})

	t.Run("test unexpected config format", func(t *testing.T) {
		absolutePath, err := filepath.Abs("./config.go")
		if err != nil {
			log.Fatal(err)
		}
		_, err = getConfigFileExtension(absolutePath)

		if !errors.Is(err, extensionErr) {
			t.Errorf("should raise an error while retrieving an extension that is not yaml nor json")
		}
	})

	t.Run("test errored json file", func(t *testing.T) {
		absolutePath, err := filepath.Abs("./_samples/test-error.json")
		if err != nil {
			log.Fatal(err)
		}

		_, parseErr := ParseConfigFile(absolutePath)
		if parseErr == nil {
			t.Error("should raise an error due to empty json file")
		}
	})
}

func testConfig(absolutePath string, t *testing.T) {
	conf, err := ParseConfigFile(absolutePath)
	if err != nil {
		t.Error(err)
	}
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
