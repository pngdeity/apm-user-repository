package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

type rootConfig struct {
	Marketplace struct {
		Packages []struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			Subdir      string `yaml:"subdir"`
		} `yaml:"packages"`
	} `yaml:"marketplace"`
}

type pkgConfig struct {
	Type string `yaml:"type"`
}

type pkgEntry struct {
	Name        string
	Type        string
	Description string
}

func main() {
	os.Exit(run())
}

func run() int {
	checkOnly := flag.Bool("check", false, "report mismatches without modifying README.md")
	flag.Parse()

	repo, err := findRepoRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}

	rootPath := filepath.Join(repo, "apm.yml")
	rootData, err := os.ReadFile(rootPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot read %s: %v\n", rootPath, err)
		return 1
	}

	var root rootConfig
	if err := yaml.Unmarshal(rootData, &root); err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot parse %s: %v\n", rootPath, err)
		return 1
	}

	var entries []pkgEntry
	for _, pkg := range root.Marketplace.Packages {
		pkgType := "unknown"
		pkgPath := filepath.Join(repo, pkg.Subdir, "apm.yml")
		if pkgData, err := os.ReadFile(pkgPath); err == nil {
			var cfg pkgConfig
			if yaml.Unmarshal(pkgData, &cfg) == nil && cfg.Type != "" {
				pkgType = cfg.Type
			}
		}
		entries = append(entries, pkgEntry{
			Name:        pkg.Name,
			Type:        pkgType,
			Description: pkg.Description,
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name < entries[j].Name
	})

	var sb strings.Builder
	sb.WriteString("| Package | Type | Description |\n")
	sb.WriteString("|---------|------|-------------|\n")
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("| **%s** | %s | %s |\n", e.Name, e.Type, e.Description))
	}
	generated := sb.String()

	readmePath := filepath.Join(repo, "README.md")
	readmeData, err := os.ReadFile(readmePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot read %s: %v\n", readmePath, err)
		return 1
	}
	readmeContent := string(readmeData)

	startMarker := "<!-- PACKAGE_TABLE_START -->"
	endMarker := "<!-- PACKAGE_TABLE_END -->"

	startIdx := strings.Index(readmeContent, startMarker)
	endIdx := strings.Index(readmeContent, endMarker)

	if startIdx == -1 || endIdx == -1 || endIdx <= startIdx {
		fmt.Fprintf(os.Stderr, "error: README.md is missing PACKAGE_TABLE_START / PACKAGE_TABLE_END markers\n")
		return 1
	}

	newContent := readmeContent[:startIdx+len(startMarker)] + "\n" + generated + readmeContent[endIdx:]

	if readmeContent == newContent {
		fmt.Println("README.md package table is up to date.")
		return 0
	}

	if *checkOnly {
		fmt.Println("README.md package table is out of date. Run 'just sync-readme' to regenerate.")
		return 1
	}

	if err := os.WriteFile(readmePath, []byte(newContent), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error: cannot write %s: %v\n", readmePath, err)
		return 1
	}

	fmt.Printf("README.md package table regenerated (%d packages).\n", len(entries))
	return 0
}

func findRepoRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not inside a git repository")
		}
		dir = parent
	}
}
