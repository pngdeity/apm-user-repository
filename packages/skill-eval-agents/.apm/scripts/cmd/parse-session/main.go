package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	os.Exit(run())
}

func run() int {
	for _, a := range os.Args[1:] {
		if a == "--help" || a == "-h" {
			printParseUsage()
			return 0
		}
	}

	cli := flag.String("cli", "", "CLI that produced the session: opencode or gemini")
	input := flag.String("input", "", "Path to session output file, or '-' for stdin")
	skillName := flag.String("skill-name", "", "Name of the skill to look for (optional)")
	flag.Parse()

	if *cli == "" {
		printParseUsage()
		return 1
	}

	var data []byte
	var err error

	if *input == "" || *input == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(*input)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading input: %v\n", err)
		return 3
	}

	var activated bool
	var evidence string
	var tokenCount int
	var durationMs int64

	switch *cli {
	case "opencode":
		activated, evidence, tokenCount, durationMs = parseOpenCodeSession(data, *skillName)
	case "gemini":
		activated, evidence, tokenCount, durationMs = parseGeminiSession(data, *skillName)
	default:
		fmt.Fprintf(os.Stderr, "unknown CLI: %s (must be opencode or gemini)\n", *cli)
		return 1
	}

	result := map[string]interface{}{
		"skill_activated":    activated,
		"activation_evidence": evidence,
		"token_count":        tokenCount,
		"duration_ms":        durationMs,
	}

	output, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling result: %v\n", err)
		return 3
	}
	fmt.Println(string(output))
	return 0
}

func printParseUsage() {
	fmt.Fprintf(os.Stderr, "usage: parse-session --cli <opencode|gemini> --input <file> [--skill-name <name>]\n")
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	fmt.Fprintf(os.Stderr, "  --cli          CLI that produced the session: opencode or gemini\n")
	fmt.Fprintf(os.Stderr, "  --input        Path to session output file, or '-' for stdin\n")
	fmt.Fprintf(os.Stderr, "  --skill-name   Name of the skill to look for (optional; if empty, matches any skill)\n")
	fmt.Fprintf(os.Stderr, "  --help, -h     Show this help\n")
	fmt.Fprintf(os.Stderr, "\nReads CLI session output and detects skill activation.\nOutputs JSON with activation evidence, token count, and duration.\n")
	fmt.Fprintf(os.Stderr, "\nexit codes: 0=success, 1=CLI not found, 2=parse error, 3=I/O error\n")
}

func parseOpenCodeSession(data []byte, skillName string) (bool, string, int, int64) {
	lines := strings.Split(string(data), "\n")
	tokenCount := 0
	durationMs := int64(0)

	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		var event map[string]interface{}
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			continue
		}

		if usage, ok := event["usage"].(map[string]interface{}); ok {
			if total, ok := usage["total_tokens"].(float64); ok {
				tokenCount = int(total)
			} else {
				if in, ok := usage["input_tokens"].(float64); ok {
					tokenCount += int(in)
				}
				if out, ok := usage["output_tokens"].(float64); ok {
					tokenCount += int(out)
				}
			}
		}

		if dur, ok := event["duration_ms"].(float64); ok {
			durationMs = int64(dur)
		}

		activated, ev := checkOpenCodeEventForSkill(event, i+1, skillName)
		if activated {
			return true, ev, tokenCount, durationMs
		}
	}
	return false, "", tokenCount, durationMs
}

func checkOpenCodeEventForSkill(event map[string]interface{}, lineNum int, skillName string) (bool, string) {
	if t, ok := event["type"].(string); ok {
		switch t {
		case "tool_use", "tool_call":
			name := ""
			if n, ok := event["name"].(string); ok {
				name = n
			}
			if fn, ok := event["function"].(map[string]interface{}); ok {
				if n, ok := fn["name"].(string); ok {
					name = n
				}
			}
			lower := strings.ToLower(name)
			if lower == "skill" || lower == "load_skill" || lower == "activate_skill" {
				return true, fmt.Sprintf("tool_use event showing %s call at line %d", name, lineNum)
			}
		case "skill", "skill_loaded", "skill_activated":
			return true, fmt.Sprintf("skill event type '%s' at line %d", t, lineNum)
		}
	}
	if _, ok := event["skill"]; ok {
		return true, fmt.Sprintf("skill field present in event at line %d", lineNum)
	}
	if sk, ok := event["skill_name"].(string); ok {
		if skillName == "" || strings.EqualFold(sk, skillName) {
			return true, fmt.Sprintf("skill_name '%s' found in event at line %d", sk, lineNum)
		}
	}
	return false, ""
}

func parseGeminiSession(data []byte, skillName string) (bool, string, int, int64) {
	lines := strings.Split(string(data), "\n")
	tokenCount := 0
	durationMs := int64(0)

	for i, line := range lines {
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
				activated, ev := checkGeminiToolUseForSkill(event, i+1, skillName)
				if activated {
					return true, ev, tokenCount, durationMs
				}
			case "result":
				if usage, ok := event["usage"].(map[string]interface{}); ok {
					if in, ok := usage["input_tokens"].(float64); ok {
						tokenCount += int(in)
					}
					if out, ok := usage["output_tokens"].(float64); ok {
						tokenCount += int(out)
					}
					if total, ok := usage["total_tokens"].(float64); ok {
						tokenCount = int(total)
					}
				}
				if dur, ok := event["duration_ms"].(float64); ok {
					durationMs = int64(dur)
				}
				if dur, ok := event["duration"].(float64); ok {
					durationMs = int64(dur)
				}
			}
		}

		if tu, ok := event["tool_use"].(map[string]interface{}); ok {
			activated, ev := checkGeminiToolUseForSkill(tu, i+1, skillName)
			if activated {
				return true, ev, tokenCount, durationMs
			}
		}
	}

	return false, "", tokenCount, durationMs
}

func checkGeminiToolUseForSkill(event map[string]interface{}, lineNum int, skillName string) (bool, string) {
	name := ""
	if n, ok := event["name"].(string); ok {
		name = n
	}
	if fn, ok := event["function"].(map[string]interface{}); ok {
		if n, ok := fn["name"].(string); ok {
			name = n
		}
	}
	lower := strings.ToLower(name)
	if lower == "skill" || lower == "load_skill" || lower == "activate_skill" {
		if skillName == "" {
			return true, fmt.Sprintf("tool_use event showing %s call at line %d", name, lineNum)
		}
		detectedSkill := extractSkillNameFromEvent(event)
		if detectedSkill != "" && strings.EqualFold(detectedSkill, skillName) {
			return true, fmt.Sprintf("tool_use event showing %s call at line %d", name, lineNum)
		}
	}
	return false, ""
}

func extractSkillNameFromEvent(event map[string]interface{}) string {
	if sn, ok := event["skill_name"].(string); ok && sn != "" {
		return sn
	}
	for _, key := range []string{"input", "arguments", "args"} {
		switch v := event[key].(type) {
		case map[string]interface{}:
			for _, f := range []string{"skill", "name", "skill_name"} {
				if s, ok := v[f].(string); ok && s != "" {
					return s
				}
			}
		case string:
			var m map[string]interface{}
			if json.Unmarshal([]byte(v), &m) == nil {
				for _, f := range []string{"skill", "name", "skill_name"} {
					if s, ok := m[f].(string); ok && s != "" {
						return s
					}
				}
			}
		}
	}
	return ""
}
