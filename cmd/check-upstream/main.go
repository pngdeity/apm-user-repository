package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type rootConfig struct {
	Marketplace struct {
		Packages []pkgEntry `yaml:"packages"`
	} `yaml:"marketplace"`
}

type pkgEntry struct {
	Name   string `yaml:"name"`
	Source string `yaml:"source"`
	Subdir string `yaml:"subdir"`
	Ref    string `yaml:"ref"`
}

type marketplaceJSON struct {
	Plugins []struct {
		Name   string `json:"name"`
		Source struct {
			SHA string `json:"sha"`
		} `json:"source"`
	} `json:"plugins"`
}

type compareResponse struct {
	Status  string `json:"status"`
	Commits []struct {
		SHA    string `json:"sha"`
		Commit struct {
			Message string `json:"message"`
		} `json:"commit"`
	} `json:"commits"`
	Files []struct {
		Filename string `json:"filename"`
	} `json:"files"`
}

var shaPattern = regexp.MustCompile(`^[a-f0-9]{40}$`)

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

	mktPath := filepath.Join(repo, ".claude-plugin", "marketplace.json")
	mktSHA := readMarketplaceSHA(mktPath)

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

		isSHA := shaPattern.MatchString(pkg.Ref)

		if isSHA {
			fixed := handleSHARef(pkg, &rootData, rootPath, checkOnly)
			if fixed < 0 {
				staleCount++
			} else {
				fixedCount += fixed
			}
		} else {
			// Named ref (branch or tag): resolve to SHA, compare against marketplace.json
			stale := handleNamedRef(pkg, mktSHA)
			if stale {
				staleCount++
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

// handleSHARef checks if a commit-SHA-pinned ref is stale (files changed
// under subdir since the pinned commit). Returns number of fixes (0 or 1)
// on success, or -1 on error.
func handleSHARef(pkg pkgEntry, rootData *[]byte, rootPath string, checkOnly *bool) int {
	newSHA, err := resolveSHARef(pkg.Source, pkg.Subdir, pkg.Ref)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[err]   %s: %v\n", pkg.Name, err)
		return -1
	}

	if newSHA == "" || newSHA == pkg.Ref {
		fmt.Printf("[ok]    %s: ref %s (HEAD, current)\n", pkg.Name, pkg.Ref[:8])
		return 0
	}

	tag := "[fix]"
	if *checkOnly {
		tag = "[stale]"
	}
	fmt.Printf("%s %s: %s -> %s  (files changed in %s)\n", tag, pkg.Name, pkg.Ref[:8], newSHA[:8], pkg.Subdir)

	if *checkOnly {
		return 0
	}

	newData := replaceRef(string(*rootData), pkg.Ref, newSHA)
	if err := os.WriteFile(rootPath, []byte(newData), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "[err]   %s: cannot write %s: %v\n", pkg.Name, rootPath, err)
		return -1
	}
	*rootData = []byte(newData)
	return 1
}

// handleNamedRef checks if a branch/tag-pinned ref's resolved SHA differs
// from what was last written into marketplace.json. Returns true if stale.
func handleNamedRef(pkg pkgEntry, mktSHA map[string]string) bool {
	parts := strings.Split(pkg.Source, "/")
	if len(parts) < 2 {
		fmt.Fprintf(os.Stderr, "[err]   %s: invalid source %q\n", pkg.Name, pkg.Source)
		return true
	}

	currentSHA, err := resolveNamedRef(parts[0], parts[1], pkg.Ref)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[err]   %s: %v\n", pkg.Name, err)
		return true
	}

	lastSHA, known := mktSHA[pkg.Name]
	if !known {
		fmt.Printf("[info]  %s: new package (no prior marketplace SHA); run apm pack\n", pkg.Name)
		return true
	}

	if currentSHA == lastSHA {
		fmt.Printf("[ok]    %s: ref %s (%s, current)\n", pkg.Name, pkg.Ref, currentSHA[:8])
		return false
	}

	fmt.Printf("[stale] %s: %s/%s -> %s  (branch %s moved, run apm pack)\n",
		pkg.Name, pkg.Ref, lastSHA[:8], currentSHA[:8], pkg.Ref)
	return true
}

func resolveSHARef(source, subdir, pinnedSHA string) (string, error) {
	parts := strings.Split(source, "/")
	if len(parts) < 2 {
		return "", fmt.Errorf("invalid source %q (expected owner/repo)", source)
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/compare/%s...HEAD",
		parts[0], parts[1], pinnedSHA)

	cr, err := fetchCompare(url)
	if err != nil {
		return "", err
	}

	if cr.Status == "identical" {
		return "", nil
	}

	for _, f := range cr.Files {
		if strings.HasPrefix(f.Filename, subdir+"/") || f.Filename == subdir {
			if len(cr.Commits) > 0 {
				return cr.Commits[len(cr.Commits)-1].SHA, nil
			}
			return "", fmt.Errorf("subdir %s changed but no commits returned", subdir)
		}
	}

	fmt.Printf("        (no changes in %s; %d total commits since pinned)\n", subdir, len(cr.Commits))
	return "", nil
}

func resolveNamedRef(owner, repo, ref string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/ref/heads/%s", owner, repo, ref)
	cr, err := fetchRef(url)
	if err != nil {
		url = fmt.Sprintf("https://api.github.com/repos/%s/%s/git/ref/tags/%s", owner, repo, ref)
		cr, err = fetchRef(url)
	}
	if err != nil {
		return "", fmt.Errorf("cannot resolve ref %q: %w", ref, err)
	}
	return cr.SHA, nil
}

type refResponse struct {
	Object struct {
		SHA string `json:"sha"`
	} `json:"object"`
	SHA string `json:"sha"`
}

func fetchRef(url string) (*refResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	setAuth(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var rr refResponse
	if err := json.Unmarshal(body, &rr); err != nil {
		return nil, err
	}
	if rr.Object.SHA != "" {
		rr.SHA = rr.Object.SHA
	}
	return &rr, nil
}

func fetchCompare(url string) (*compareResponse, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	setAuth(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("repo or ref not found (HTTP 404)")
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API returned %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var cr compareResponse
	if err := json.Unmarshal(body, &cr); err != nil {
		return nil, err
	}
	return &cr, nil
}

func setAuth(req *http.Request) {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("GH_TOKEN")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
}

func readMarketplaceSHA(path string) map[string]string {
	m := make(map[string]string)
	data, err := os.ReadFile(path)
	if err != nil {
		return m
	}
	var mkt marketplaceJSON
	if err := json.Unmarshal(data, &mkt); err != nil {
		return m
	}
	for _, p := range mkt.Plugins {
		if p.Source.SHA != "" {
			m[p.Name] = p.Source.SHA
		}
	}
	return m
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
