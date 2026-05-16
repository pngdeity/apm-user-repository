package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const timeout = 120 * time.Second

func main() {
	os.Exit(run())
}

func run() int {
	for _, a := range os.Args[1:] {
		if a == "--help" || a == "-h" {
			printInvokeUsage()
			return 0
		}
	}

	prompt := flag.String("prompt", "", "Prompt to send to the agent")
	cli := flag.String("cli", "", "CLI to invoke: opencode or gemini")
	workspace := flag.String("workspace", "", "Workspace directory")
	skill := flag.String("skill", "", "Path to skill directory/file")
	newSession := flag.Bool("new-session", true, "Start a new session for isolation (gemini only)")
	flag.Parse()

	if *prompt == "" || *cli == "" || *workspace == "" {
		printInvokeUsage()
		return 1
	}

	switch *cli {
	case "opencode":
		return runOpenCode(*prompt, *workspace, *skill)
	case "gemini":
		return runGemini(*prompt, *workspace, *skill, *newSession)
	default:
		fmt.Fprintf(os.Stderr, "unknown CLI: %s (must be opencode or gemini)\n", *cli)
		return 1
	}
}

func printInvokeUsage() {
	fmt.Fprintf(os.Stderr, "usage: invoke-cli --prompt <string> --cli <opencode|gemini> --workspace <dir> [--skill <path>] [--new-session]\n")
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	fmt.Fprintf(os.Stderr, "  --prompt        Prompt to send to the agent\n")
	fmt.Fprintf(os.Stderr, "  --cli           CLI to invoke: opencode or gemini\n")
	fmt.Fprintf(os.Stderr, "  --workspace     Workspace directory\n")
	fmt.Fprintf(os.Stderr, "  --skill         Path to skill directory/file (optional)\n")
	fmt.Fprintf(os.Stderr, "  --new-session   Start a new session for isolation (default: true, gemini only)\n")
	fmt.Fprintf(os.Stderr, "  --help, -h      Show this help\n")
	fmt.Fprintf(os.Stderr, "\nexit codes: 0=success, 1=CLI/prep error, 2=timeout, 3=parse error\n")
}

func runOpenCode(prompt, workspace, skillPath string) int {
	absWorkspace, err := filepath.Abs(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error resolving workspace path: %v\n", err)
		return 1
	}

	tmpDir, err := os.MkdirTemp("", "opencode-eval-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating temp workspace: %v\n", err)
		return 1
	}
	defer os.RemoveAll(tmpDir)

	skillActivated := false
	if skillPath != "" {
		agentsDir := filepath.Join(tmpDir, ".agents", "skills")
		if err := os.MkdirAll(agentsDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error creating .agents/skills dir: %v\n", err)
			return 1
		}
		skillName := filepath.Base(skillPath)
		dest := filepath.Join(agentsDir, skillName)
		if err := copyDir(skillPath, dest); err != nil {
			fmt.Fprintf(os.Stderr, "error copying skill: %v\n", err)
			return 1
		}
	}

	if _, err := exec.LookPath("opencode"); err != nil {
		fmt.Fprintf(os.Stderr, "opencode CLI not found in PATH\n")
		return 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if skillPath != "" {
		agentsDir := filepath.Join(absWorkspace, ".agents", "skills")
		if err := os.MkdirAll(agentsDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error creating .agents/skills in workspace: %v\n", err)
			return 1
		}
		skillName := filepath.Base(skillPath)
		dest := filepath.Join(agentsDir, skillName)
		if err := copyDir(skillPath, dest); err != nil {
			fmt.Fprintf(os.Stderr, "error copying skill to workspace: %v\n", err)
			return 1
		}
	}

	args := []string{"run", prompt}
	if absWorkspace != "" {
		args = append(args, "--dir", absWorkspace)
	}
	args = append(args, "--format", "json", "--dangerously-skip-permissions")

	cmd := exec.CommandContext(ctx, "opencode", args...)
	cmd.Dir = tmpDir
	cmd.Env = append(os.Environ(), "NO_COLOR=1")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err = cmd.Run()
	elapsed := time.Since(start)
	durationMs := elapsed.Milliseconds()

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Fprintf(os.Stderr, "invocation timed out after %v\n", timeout)
		return 2
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "opencode run error: %v\nstderr: %s\n", err, stderr.String())
	}

	stdoutBytes := stdout.Bytes()

	sessionID := ""
	lines := strings.Split(strings.TrimSpace(stdout.String()), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		var event map[string]interface{}
		if json.Unmarshal([]byte(line), &event) == nil {
			if sid, ok := event["session_id"].(string); ok && sid != "" {
				sessionID = sid
				break
			}
		}
	}

	sessionOutput := ""
	tokenCount := 0
	if sessionID != "" {
		expCtx, expCancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer expCancel()

		expCmd := exec.CommandContext(expCtx, "opencode", "export", sessionID)
		expCmd.Dir = tmpDir
		var expOut bytes.Buffer
		var expErr bytes.Buffer
		expCmd.Stdout = &expOut
		expCmd.Stderr = &expErr

		if expErr2 := expCmd.Run(); expErr2 == nil {
			sessionOutput = strings.TrimSpace(expOut.String())
			var exportData map[string]interface{}
			if json.Unmarshal(expOut.Bytes(), &exportData) == nil {
				if usage, ok := exportData["usage"].(map[string]interface{}); ok {
					if total, ok := usage["input_tokens"].(float64); ok {
						tokenCount += int(total)
					}
					if total, ok := usage["output_tokens"].(float64); ok {
						tokenCount += int(total)
					}
					if total, ok := usage["total_tokens"].(float64); ok {
						tokenCount = int(total)
					}
				}
			}
		}
	}

	skillActivated = detectSkillActivationOpenCode(stdoutBytes)

	result := map[string]interface{}{
		"cli":             "opencode",
		"skill_activated": skillActivated,
		"token_count":     tokenCount,
		"duration_ms":     durationMs,
		"session_output":  sessionOutput,
	}

	output, _ := json.Marshal(result)
	fmt.Println(string(output))
	return 0
}

func runGemini(prompt, workspace, skillPath string, newSession bool) int {
	absWorkspace, err := filepath.Abs(workspace)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error resolving workspace path: %v\n", err)
		return 1
	}

	tmpDir, err := os.MkdirTemp("", "gemini-eval-")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating temp workspace: %v\n", err)
		return 1
	}
	defer os.RemoveAll(tmpDir)

	if skillPath != "" {
		agentsDir := filepath.Join(tmpDir, ".agents", "skills")
		if err := os.MkdirAll(agentsDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error creating .agents/skills dir: %v\n", err)
			return 1
		}
		skillName := filepath.Base(skillPath)
		dest := filepath.Join(agentsDir, skillName)
		if err := copyDir(skillPath, dest); err != nil {
			fmt.Fprintf(os.Stderr, "error copying skill: %v\n", err)
			return 1
		}

		agentsDir2 := filepath.Join(absWorkspace, ".agents", "skills")
		if err := os.MkdirAll(agentsDir2, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "error creating .agents/skills in workspace: %v\n", err)
			return 1
		}
		dest2 := filepath.Join(agentsDir2, skillName)
		if err := copyDir(skillPath, dest2); err != nil {
			fmt.Fprintf(os.Stderr, "error copying skill to workspace: %v\n", err)
			return 1
		}
	}

	if _, err := exec.LookPath("gemini"); err != nil {
		fmt.Fprintf(os.Stderr, "gemini CLI not found in PATH\n")
		return 1
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	args := []string{"-p", prompt, "--output-format", "stream-json"}
	if newSession {
		args = append(args, "--new-session")
	}
	cmd := exec.CommandContext(ctx, "gemini", args...)
	cmd.Dir = tmpDir
	cmd.Env = append(os.Environ(), "NO_COLOR=1")

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err = cmd.Run()
	elapsed := time.Since(start)
	durationMs := elapsed.Milliseconds()

	if ctx.Err() == context.DeadlineExceeded {
		fmt.Fprintf(os.Stderr, "invocation timed out after %v\n", timeout)
		return 2
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "gemini run error: %v\nstderr: %s\n", err, stderr.String())
	}

	stdoutBytes := stdout.Bytes()

	skillActivated := detectSkillActivationGemini(stdoutBytes)
	tokenCount := extractTokenCountGemini(stdoutBytes)
	if streamDuration := extractDurationGemini(stdoutBytes); streamDuration > 0 {
		durationMs = streamDuration
	}

	result := map[string]interface{}{
		"cli":             "gemini",
		"skill_activated": skillActivated,
		"token_count":     tokenCount,
		"duration_ms":     durationMs,
		"session_output":  strings.TrimSpace(stdout.String()),
	}

	output, _ := json.Marshal(result)
	fmt.Println(string(output))
	return 0
}

func detectSkillActivationOpenCode(data []byte) bool {
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var event map[string]interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		if t, ok := event["type"].(string); ok {
			switch t {
			case "tool_use", "tool_call", "tool_result":
				if name, ok := event["name"].(string); ok {
					lower := strings.ToLower(name)
					if lower == "skill" || lower == "load_skill" || lower == "activate_skill" {
						return true
					}
				}
				if fn, ok := event["function"].(map[string]interface{}); ok {
					if name, ok := fn["name"].(string); ok {
						lower := strings.ToLower(name)
						if lower == "skill" || lower == "load_skill" || lower == "activate_skill" {
							return true
						}
					}
				}
			case "skill", "skill_loaded", "skill_activated":
				return true
			}
		}
		if _, ok := event["skill"]; ok {
			return true
		}
	}
	return false
}

func detectSkillActivationGemini(data []byte) bool {
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var event map[string]interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		if t, ok := event["type"].(string); ok {
			switch t {
			case "tool_use":
				if name, ok := event["name"].(string); ok {
					lower := strings.ToLower(name)
					if lower == "skill" || lower == "load_skill" || lower == "activate_skill" {
						return true
					}
				}
				if fn, ok := event["function"].(map[string]interface{}); ok {
					if name, ok := fn["name"].(string); ok {
						lower := strings.ToLower(name)
						if lower == "skill" || lower == "load_skill" || lower == "activate_skill" {
							return true
						}
					}
				}
			case "result", "done":
				if name, ok := event["triggered_skills"].([]interface{}); ok && len(name) > 0 {
					return true
				}
			}
		}
		if tu, ok := event["tool_use"].(map[string]interface{}); ok {
			if name, ok := tu["name"].(string); ok {
				lower := strings.ToLower(name)
				if lower == "skill" || lower == "load_skill" || lower == "activate_skill" {
					return true
				}
			}
		}
	}
	return false
}

func extractTokenCountGemini(data []byte) int {
	lines := strings.Split(string(data), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		var event map[string]interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		if t, ok := event["type"].(string); ok && t == "result" {
			if usage, ok := event["usage"].(map[string]interface{}); ok {
				if input, ok := usage["input_tokens"].(float64); ok {
					total := int(input)
					if output, ok := usage["output_tokens"].(float64); ok {
						total += int(output)
					}
					return total
				}
				if total, ok := usage["total_tokens"].(float64); ok {
					return int(total)
				}
			}
		}
	}
	return 0
}

func extractDurationGemini(data []byte) int64 {
	lines := strings.Split(string(data), "\n")
	for i := len(lines) - 1; i >= 0; i-- {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		var event map[string]interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}
		if t, ok := event["type"].(string); ok && t == "result" {
			if dur, ok := event["duration_ms"].(float64); ok {
				return int64(dur)
			}
			if dur, ok := event["duration"].(float64); ok {
				return int64(dur)
			}
		}
	}
	return 0
}

func copyDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := copyDir(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}
