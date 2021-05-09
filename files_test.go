package main

import (
	"strings"
	"testing"
)

func TestGlobbingPatterns(t *testing.T) {
	t.Run("test star pattern", func(t *testing.T) {
		pattern := "*.go"
		files := globFiles(pattern)
		if len(files) != 3 {
			t.Errorf("A .go file hasn't been globbed, check pattern")
		}
	})

	t.Run("test dot files", func(t *testing.T) {
		pattern := ".*"
		files := globFiles(pattern)
		for i := range files {
			if !(strings.Contains(files[i], ".git") || strings.Contains(files[i], ".travis")) {
				t.Errorf("dotfiles should only be related to git or travis")
			}
		}
	})
}
