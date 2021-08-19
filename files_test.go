package main

import (
	"strings"
	"testing"
)

func TestSimplePatterns(t *testing.T) {
	t.Run("*.go pattern", func(t *testing.T) {
		// file_test.go main.go
		pattern := "*.go"
		files := globFiles(pattern)
		if len(files) != 2 {
			t.Errorf("A .go file hasn't been globbed, check pattern")
		}
	})

	t.Run(".git pattern", func(t *testing.T) {
		// .gitignore
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
		// file_test.go main.go test/server.go
		pattern := "**/*.go"
		files := globFiles(pattern)
		if len(files) != 3 {
			t.Errorf("A .go file hasn't been globbed, check pattern")
		}
	})
}

func TestMultiplePatterns(t *testing.T) {
	t.Run("*.go and *.yml pattern", func(t *testing.T) {
		var filesCount int
		// file_test.go main.go .github/workflows/go.yml
		expression := "*.go **/**/.yml"
		patterns := strings.Split(expression, " ")
		for i := range patterns {
			files := globFiles(patterns[i])
			filesCount += len(files)
		}
		if filesCount != 3 {
			t.Errorf("A file matching *.go or *.yml hasn't been globbed")
		}
	})
}
