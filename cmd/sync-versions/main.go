package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/Masterminds/semver/v3"
	"gopkg.in/yaml.v3"
)

type rootConfig struct {
	Marketplace struct {
		Packages []struct {
			Name    string `yaml:"name"`
			Version string `yaml:"version"`
			Subdir  string `yaml:"subdir"`
		} `yaml:"packages"`
	} `yaml:"marketplace"`
}

type pkgConfig struct {
	Version string `yaml:"version"`
}

func main() {
	os.Exit(run())
}

func run() int {
	checkOnly := flag.Bool("check", false, "report mismatches without modifying files")
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

	mismatches := 0
	for _, pkg := range root.Marketplace.Packages {
		if pkg.Version == "" {
			fmt.Printf("[ref]   %s: pinned to ref (no version constraint)\n", pkg.Name)
			continue
		}

		constraint, err := semver.NewConstraint(pkg.Version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[err]   %s: invalid constraint %q: %v\n", pkg.Name, pkg.Version, err)
			mismatches++
			continue
		}

		pkgPath := filepath.Join(repo, pkg.Subdir, "apm.yml")
		if _, statErr := os.Stat(pkgPath); os.IsNotExist(statErr) {
			fmt.Fprintf(os.Stderr, "[warn]  %s: package apm.yml not found at %s\n", pkg.Name, pkg.Subdir)
			continue
		}

		pkgData, err := os.ReadFile(pkgPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[err]   %s: cannot read %s: %v\n", pkg.Name, pkgPath, err)
			mismatches++
			continue
		}

		var pkgCfg pkgConfig
		if err := yaml.Unmarshal(pkgData, &pkgCfg); err != nil {
			fmt.Fprintf(os.Stderr, "[err]   %s: cannot parse %s: %v\n", pkg.Name, pkgPath, err)
			mismatches++
			continue
		}

		currentVer, err := semver.NewVersion(pkgCfg.Version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[err]   %s: invalid version %q: %v\n", pkg.Name, pkgCfg.Version, err)
			mismatches++
			continue
		}

		if constraint.Check(currentVer) {
			fmt.Printf("[ok]    %s: %s  (constraint: %s)\n", pkg.Name, currentVer, pkg.Version)
			continue
		}

		newVer := extractBaseVersion(pkg.Version)
		tag := "[fix]"
		if *checkOnly {
			tag = "[mismatch]"
		}
		fmt.Printf("%s %s: %s -> %s  (constraint: %s)\n", tag, pkg.Name, currentVer, newVer, pkg.Version)
		mismatches++

		if !*checkOnly {
			newData := replaceVersion(string(pkgData), pkgCfg.Version, newVer)
			if err := os.WriteFile(pkgPath, []byte(newData), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "[err]   %s: cannot write %s: %v\n", pkg.Name, pkgPath, err)
			}
		}
	}

	fmt.Printf("\n%s\n", strings.Repeat("-", 40))
	if mismatches > 0 {
		if *checkOnly {
			fmt.Printf("Found %d version mismatch(es). Run without --check to auto-fix, then run 'apm pack'.\n", mismatches)
		} else {
			fmt.Printf("Fixed %d package(s). Run 'apm pack' to regenerate marketplace.json.\n", mismatches)
		}
		return 1
	}
	fmt.Println("All package versions satisfy root constraints.")
	return 0
}

// extractBaseVersion strips the caret (and any leading operator) from a version
// constraint to get the minimum valid version. E.g., "^0.3.0" -> "0.3.0".
func extractBaseVersion(constraint string) string {
	s := strings.TrimSpace(constraint)
	s = strings.TrimLeftFunc(s, func(r rune) bool {
		return !unicode.IsDigit(r)
	})
	return s
}

// replaceVersion replaces the version line in the apm.yml content.
func replaceVersion(content, oldVer, newVer string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "version:") {
			lines[i] = fmt.Sprintf("version: %s", newVer)
			break
		}
	}
	return strings.Join(lines, "\n")
}

// findRepoRoot walks up from the working directory to find the git repository root.
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
