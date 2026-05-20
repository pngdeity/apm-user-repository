package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type rootConfig struct {
	Marketplace struct {
		Packages []struct {
			Name   string `yaml:"name"`
			Source string `yaml:"source"`
			Subdir string `yaml:"subdir"`
			Ref    string `yaml:"ref"`
		} `yaml:"packages"`
	} `yaml:"marketplace"`
}

type compareResponse struct {
	Status string `json:"status"`
	Files  []struct {
		Filename string `json:"filename"`
		Status   string `json:"status"`
	} `json:"files"`
	Commits []struct {
		SHA    string `json:"sha"`
		Commit struct {
			Message string `json:"message"`
		} `json:"commit"`
	} `json:"commits"`
}

func main() {
	os.Exit(run())
}

func run() int {
	checkOnly := flag.Bool("check", false, "report staleness without modifying files")
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

	selfSource := "pngdeity/apm-user-repository"
	staleCount := 0
	fixedCount := 0

	for _, pkg := range root.Marketplace.Packages {
		if pkg.Source == "" || pkg.Source == selfSource {
			continue
		}
		if pkg.Source == "./" || strings.HasPrefix(pkg.Source, "./") {
			continue
		}
		if pkg.Ref == "" {
			continue
		}

		newSHA, err := checkUpstream(pkg.Source, pkg.Subdir, pkg.Ref)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[err]   %s: %v\n", pkg.Name, err)
			staleCount++
			continue
		}

		if newSHA == "" || newSHA == pkg.Ref {
			fmt.Printf("[ok]    %s: ref %s (HEAD, current)\n", pkg.Name, pkg.Ref[:8])
			continue
		}

		tag := "[fix]"
		if *checkOnly {
			tag = "[stale]"
		}
		fmt.Printf("%s %s: %s -> %s  (files changed in %s)\n", tag, pkg.Name, pkg.Ref[:8], newSHA[:8], pkg.Subdir)
		staleCount++

		if !*checkOnly {
			newData := replaceRef(string(rootData), pkg.Ref, newSHA)
			if err := os.WriteFile(rootPath, []byte(newData), 0644); err != nil {
				fmt.Fprintf(os.Stderr, "[err]   %s: cannot write %s: %v\n", pkg.Name, rootPath, err)
			} else {
				fixedCount++
				rootData = []byte(newData)
			}
		}
	}

	if staleCount > 0 {
		fmt.Printf("\n%s\n", strings.Repeat("-", 40))
		if *checkOnly {
			fmt.Printf("Found %d stale external ref(s). Run without --check to auto-update, then run 'apm pack'.\n", staleCount)
		} else {
			fmt.Printf("Updated %d external ref(s). Run 'apm pack' to regenerate marketplace.json.\n", fixedCount)
		}
		return 1
	}

	fmt.Printf("%s\nAll external refs are current.\n", strings.Repeat("-", 40))
	return 0
}

func checkUpstream(source, subdir, pinnedSHA string) (string, error) {
	parts := strings.Split(source, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid source %q (expected owner/repo)", source)
	}
	owner := parts[0]
	repoName := parts[1]

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/compare/%s...HEAD",
		owner, repoName, pinnedSHA)

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("GH_TOKEN")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("cannot create request: %w", err)
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("cannot read response: %w", err)
	}

	if resp.StatusCode == 404 {
		return "", fmt.Errorf("repo %s not found or ref %s unreachable", source, pinnedSHA[:8])
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var cr compareResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return "", fmt.Errorf("cannot parse response: %w", err)
	}

	if cr.Status == "identical" {
		return "", nil
	}

	subdirChanged := false
	for _, f := range cr.Files {
		if strings.HasPrefix(f.Filename, subdir+"/") || f.Filename == subdir {
			subdirChanged = true
			break
		}
	}

	if !subdirChanged {
		fmt.Printf("        (no changes in %s; %d total commits since pinned)\n", subdir, len(cr.Commits))
		return "", nil
	}

	if len(cr.Commits) > 0 {
		return cr.Commits[len(cr.Commits)-1].SHA, nil
	}
	return "", fmt.Errorf("subdir %s changed but no commits returned", subdir)
}

func replaceRef(content, oldRef, newRef string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "ref:") && strings.Contains(trimmed, oldRef) {
			fields := strings.Fields(line)
			for j, f := range fields {
				if f == oldRef {
					fields[j] = newRef
				}
			}
			lines[i] = strings.Join(fields, " ")
			break
		}
	}
	return strings.Join(lines, "\n")
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
