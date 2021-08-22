package files

import (
	"fmt"
	"strings"
	"testing"
)

func TestSimplePatterns(t *testing.T) {
	t.Run("*.go pattern", func(t *testing.T) {
		// glob_test.go glob.go keep_or_remove_file.go watcher.go
		pattern := "*.go"
		files := GlobFiles(pattern)
		if len(files) != 4 {
			t.Errorf("A .go file hasn't been globbed, check pattern")
		}
	})

	t.Run(".git pattern", func(t *testing.T) {
		// .gitignore
		pattern := ".git"
		files := GlobFiles(pattern)
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
		// [files_test.go main.go _samples/a/a.go _samples/a/a2.go _samples/b/b.go]
		pattern := "**/**/*.go"
		files := GlobFiles(pattern)
		fmt.Println("TestDoubleStarPatterns", files)
		if len(files) != 5 {
			t.Errorf("A .go file hasn't been globbed when pattern is 3 level nested")
		}
	})

	t.Run("5 level nested double star pattern", func(t *testing.T) {
		// file_test.go main.go test/server.go
		// main.go main_test.go
		// _samples/a/a.go _samples/a/a2.go _samples/b/b.go
		// _samples/a/aa/aa.go _samples/a/aa/aa2.go _samples/b/ba/ba.go _samples/b/ba/ba2.go
		// _samples/a/aa/aaa/aaa.go _samples/a/aa/aaa/aaa2.go
		pattern := "**/**/**/**/*.go"
		files := GlobFiles(pattern)
		if len(files) != 11 {
			t.Errorf("A .go file hasn't been globbed when pattern is 5 level nested")
		}
	})
}

func TestMultiplePatternsWithWildcardPattern(t *testing.T) {
	t.Run("*.go and *.yml pattern", func(t *testing.T) {
		var filesCount int
		// file_test.go main.go .github/workflows/go.yml _samples/b/b.yml
		expression := "*.go **/**/*.yml"
		patterns := strings.Split(expression, " ")
		for i := range patterns {
			files := GlobFiles(patterns[i])
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
		// file_test.go main.go _samples/b/b.yml
		expression := "*.go **/**/go.yml"
		patterns := strings.Split(expression, " ")
		for i := range patterns {
			files := GlobFiles(patterns[i])
			filesCount += len(files)
		}
		if filesCount != 2 {
			t.Errorf("A file matching *.go or *.yml hasn't been globbed")
		}
	})
}
