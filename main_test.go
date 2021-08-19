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
	t.Run("3 level nested double star pattern", func(t *testing.T) {
		// file_test.go main.go test/server.go
		// [files_test.go main.go samples/a/a.go samples/a/a2.go samples/b/b.go]
		pattern := "**/**/*.go"
		files := globFiles(pattern)
		if len(files) != 5 {
			t.Errorf("A .go file hasn't been globbed when pattern is 3 level nested")
		}
	})

	t.Run("5 level nested double star pattern", func(t *testing.T) {
		// file_test.go main.go test/server.go
		// main.go main_test.go
		// samples/a/a.go samples/a/a2.go samples/b/b.go
		// samples/a/aa/aa.go samples/a/aa/aa2.go samples/b/ba/ba.go samples/b/ba/ba2.go
		// samples/a/aa/aaa/aaa.go samples/a/aa/aaa/aaa2.go
		pattern := "**/**/**/**/*.go"
		files := globFiles(pattern)
		if len(files) != 11 {
			t.Errorf("A .go file hasn't been globbed when pattern is 5 level nested")
		}
	})
}

func TestMultiplePatternsWithWildcardPattern(t *testing.T) {
	t.Run("*.go and *.yml pattern", func(t *testing.T) {
		var filesCount int
		// file_test.go main.go .github/workflows/go.yml samples/b/b.yml
		expression := "*.go **/**/*.yml"
		patterns := strings.Split(expression, " ")
		for i := range patterns {
			files := globFiles(patterns[i])
			filesCount += len(files)
		}
		if filesCount != 4 {
			t.Errorf("A file matching *.go or *.yml hasn't been globbed")
		}
	})
}

func TestMultiplePatternsWithoutWildcardPattern(t *testing.T) {
	t.Run("*.go and go.yml pattern", func(t *testing.T) {
		var filesCount int
		// file_test.go main.go samples/b/b.yml
		expression := "*.go **/**/go.yml"
		patterns := strings.Split(expression, " ")
		for i := range patterns {
			files := globFiles(patterns[i])
			filesCount += len(files)
		}
		if filesCount != 2 {
			t.Errorf("A file matching *.go or *.yml hasn't been globbed")
		}
	})
}
