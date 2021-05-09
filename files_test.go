package main

import (
	"strings"
	"testing"
)

func TestSimplePatterns(t *testing.T) {
	t.Run("*.go pattern", func(t *testing.T) {
		pattern := "*.go"
		files := globFiles(pattern)
		if len(files) != 2 {
			t.Errorf("A .go file hasn't been globbed, check pattern")
		}
	})

	t.Run(".git pattern", func(t *testing.T) {
		pattern := ".git"
		files := globFiles(pattern)
		for i := range files {
			if !(strings.Contains(files[i], ".git")) {
				t.Errorf("As we ignore directories, it should only return one file : .gitignore")
			}
		}
	})
}

func TestDoubleStarPatterns(t *testing.T) {
	t.Run("**/*.go pattern", func(t *testing.T) {
		pattern := "**/*.go"
		files := globFiles(pattern)
		if len(files) != 3 {
			t.Errorf("A .go file hasn't been globbed, check pattern")
		}
	})

	t.Run(".git/**/*.sample pattern", func(t *testing.T) {
		pattern := ".git/**/*.sample"
		files := globFiles(pattern)
		if len(files) != 11 {
			t.Errorf("A .sample file hasn't been globbed")
		}
	})

	t.Run("**/**/*.sample pattern", func(t *testing.T) {
		pattern := "**/**/*.sample"
		files := globFiles(pattern)
		if len(files) != 11 {
			t.Errorf("A .sample file hasn't been globbed")
		}
	})
}
