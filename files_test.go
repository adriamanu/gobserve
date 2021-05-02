package main

import (
	"strings"
	"testing"
)

func TestGlobbingPatterns(t *testing.T) {
	t.Run("test star pattern", func(t *testing.T) {
		pattern := "*.go"
		files := retrieveFilesToWatch(pattern)
		if len(files) != 2 {
			t.Errorf("A .go file hasn't been globbed, check pattern")
		}
	})

	t.Run("test dot files", func(t *testing.T) {
		pattern := ".*"
		files := retrieveFilesToWatch(pattern)
		for i := range retrieveFilesToWatch(".*") {
			if !strings.Contains(files[i], ".git") {
				t.Errorf("dotfiles should only be related to git")
			}
		}
	})
}
