package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type candidateResult struct {
	Name           string  `json:"name"`
	CompositeScore float64 `json:"composite_score"`
	PassRate       float64 `json:"pass_rate"`
	TriggerRate    float64 `json:"trigger_rate"`
	TokenDelta     float64 `json:"token_delta"`
}

type thresholdCheck struct {
	Name         string  `json:"name"`
	PassRate     float64 `json:"pass_rate"`
	TriggerRate  float64 `json:"trigger_rate"`
	PassRateOk   bool    `json:"pass_rate_ok"`
	TriggerRateOk bool   `json:"trigger_rate_ok"`
	Passed       bool    `json:"passed"`
}

type failureReason struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

type selectedOutput struct {
	Selected        string            `json:"selected"`
	SelectedPath    string            `json:"selected_path"`
	CompositeScore  float64           `json:"composite_score"`
	Rankings        []candidateResult `json:"rankings"`
	ThresholdChecks []thresholdCheck  `json:"threshold_checks"`
	FailureReasons  []failureReason   `json:"failure_reasons"`
}

func main() {
	os.Exit(run())
}

func run() int {
	for _, a := range os.Args[1:] {
		if a == "--help" || a == "-h" {
			printSelectUsage()
			return 0
		}
	}

	workspace := flag.String("workspace", "", "Path to eval workspace directory")
	candidatesFlag := flag.String("candidates", "original,revision-A,revision-B,revision-C", "Comma-separated candidate names")
	flag.Parse()

	if *workspace == "" {
		printSelectUsage()
		return 1
	}

	candidateNames := strings.Split(*candidatesFlag, ",")
	for i := range candidateNames {
		candidateNames[i] = strings.TrimSpace(candidateNames[i])
	}

	var results []candidateResult
	var failureReasons []failureReason

	for _, name := range candidateNames {
		candidateDir := filepath.Join(*workspace, name)
		if _, err := os.Stat(candidateDir); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "warning: candidate directory not found: %s\n", candidateDir)
			failureReasons = append(failureReasons, failureReason{Name: name, Reason: "directory_not_found"})
			continue
		}
		cr, frs := evaluateCandidate(candidateDir, name)
		results = append(results, cr)
		failureReasons = append(failureReasons, frs...)
	}

	if len(results) == 0 {
		fmt.Fprintf(os.Stderr, "no candidate directories found in workspace\n")
		outData, _ := json.MarshalIndent(selectedOutput{
			FailureReasons: failureReasons,
		}, "", "  ")
		selectedPath := filepath.Join(*workspace, "selected.json")
		os.WriteFile(selectedPath, outData, 0644)
		fmt.Println(string(outData))
		return 3
	}

	for _, cr := range results {
		fmt.Fprintf(os.Stderr, "candidate %s: composite=%.4f pass_rate=%.4f trigger_rate=%.4f\n",
			cr.Name, cr.CompositeScore, cr.PassRate, cr.TriggerRate)
	}

	const minPassRate = 0.5
	const minTriggerRate = 0.5
	var thresholdChecks []thresholdCheck
	passedMap := make(map[string]bool)

	for i := range results {
		cr := &results[i]
		tc := thresholdCheck{
			Name:         cr.Name,
			PassRate:     cr.PassRate,
			TriggerRate:  cr.TriggerRate,
			PassRateOk:   cr.PassRate >= minPassRate,
			TriggerRateOk: cr.TriggerRate >= minTriggerRate,
		}
		tc.Passed = tc.PassRateOk && tc.TriggerRateOk
		thresholdChecks = append(thresholdChecks, tc)
		passedMap[cr.Name] = tc.Passed

		if !tc.Passed {
			var reasons []string
			if !tc.PassRateOk {
				reasons = append(reasons, "pass_rate below 0.5")
			}
			if !tc.TriggerRateOk {
				reasons = append(reasons, "trigger_rate below 0.5")
			}
			failureReasons = append(failureReasons, failureReason{
				Name:   cr.Name,
				Reason: strings.Join(reasons, "; "),
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].CompositeScore > results[j].CompositeScore
	})

	var selectedName string
	var selectedScore float64
	for _, cr := range results {
		if passedMap[cr.Name] {
			selectedName = cr.Name
			selectedScore = cr.CompositeScore
			break
		}
	}

	if selectedName == "" {
		outData, err := json.MarshalIndent(selectedOutput{
			Rankings:        results,
			ThresholdChecks: thresholdChecks,
			FailureReasons:  failureReasons,
		}, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "error marshaling output: %v\n", err)
			return 3
		}
		selectedPath := filepath.Join(*workspace, "selected.json")
		if err := os.WriteFile(selectedPath, outData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "error writing selected.json: %v\n", err)
			return 3
		}
		fmt.Println(string(outData))
		return 2
	}

	output := selectedOutput{
		Selected:        selectedName,
		SelectedPath:    filepath.Join(*workspace, selectedName),
		CompositeScore:  math.Round(selectedScore*100) / 100,
		Rankings:        results,
		ThresholdChecks: thresholdChecks,
		FailureReasons:  failureReasons,
	}

	outData, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling output: %v\n", err)
		return 3
	}

	selectedPath := filepath.Join(*workspace, "selected.json")
	if err := os.WriteFile(selectedPath, outData, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing selected.json: %v\n", err)
		return 3
	}

	fmt.Println(string(outData))
	return 0
}

func evaluateCandidate(dir, name string) (candidateResult, []failureReason) {
	benchmarkPath := filepath.Join(dir, "benchmark.json")
	data, err := os.ReadFile(benchmarkPath)
	if err != nil {
		return candidateResult{
			Name:           name,
			CompositeScore: 0,
			PassRate:       0,
			TriggerRate:    0,
			TokenDelta:     0,
		}, []failureReason{{Name: name, Reason: "missing_data"}}
	}

	var benchmark struct {
		RunSummary struct {
			WithSkill struct {
				PassRate    struct{ Mean float64 } `json:"pass_rate"`
				TimeSeconds struct{ Mean float64 } `json:"time_seconds"`
				Tokens      struct{ Mean float64 } `json:"tokens"`
			} `json:"with_skill"`
			WithoutSkill struct {
				PassRate    struct{ Mean float64 } `json:"pass_rate"`
				TimeSeconds struct{ Mean float64 } `json:"time_seconds"`
				Tokens      struct{ Mean float64 } `json:"tokens"`
			} `json:"without_skill"`
			Delta struct {
				PassRate    float64 `json:"pass_rate"`
				TimeSeconds float64 `json:"time_seconds"`
				Tokens      float64 `json:"tokens"`
			} `json:"delta"`
		} `json:"run_summary"`
	}

	if err := json.Unmarshal(data, &benchmark); err != nil {
		fmt.Fprintf(os.Stderr, "warning: error parsing benchmark for %s: %v\n", name, err)
		return candidateResult{
			Name:           name,
			CompositeScore: 0,
			PassRate:       0,
			TriggerRate:    0,
			TokenDelta:     0,
		}, []failureReason{{Name: name, Reason: "missing_data"}}
	}

	summary := benchmark.RunSummary
	passRate := summary.WithSkill.PassRate.Mean
	tokenDelta := summary.Delta.Tokens

	triggerRate := readTriggerRate(dir)

	withTokens := summary.WithSkill.Tokens.Mean
	withoutTokens := summary.WithoutSkill.Tokens.Mean
	tokenEfficiencyDelta := 0.0
	if withoutTokens > 0 {
		raw := 1.0 - (withTokens / withoutTokens)
		tokenEfficiencyDelta = clamp(raw, 0, 1)
	}

	composite := (passRate * 0.5) + (triggerRate * 0.3) + (tokenEfficiencyDelta * 0.2)

	return candidateResult{
		Name:           name,
		CompositeScore: math.Round(composite*100) / 100,
		PassRate:       math.Round(passRate*100) / 100,
		TriggerRate:    math.Round(triggerRate*100) / 100,
		TokenDelta:     tokenDelta,
	}, nil
}

func readTriggerRate(dir string) float64 {
	triggerPath := filepath.Join(dir, "trigger.json")
	data, err := os.ReadFile(triggerPath)
	if err != nil {
		skillDir := filepath.Join(dir, "with_skill")
		return findTriggerRateFromSession(skillDir)
	}
	var trigger struct {
		TriggerRate float64 `json:"trigger_rate"`
		Activated   bool    `json:"activated"`
	}
	if err := json.Unmarshal(data, &trigger); err != nil {
		return 0
	}
	if trigger.TriggerRate > 0 {
		return trigger.TriggerRate
	}
	if trigger.Activated {
		return 1.0
	}
	return 0
}

func findTriggerRateFromSession(dir string) float64 {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return 0
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}

	var activations, total int
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".json" && entry.Name() != "grading.json" && entry.Name() != "timing.json" {
			data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				continue
			}
			var session struct {
				SkillActivated bool `json:"skill_activated"`
			}
			if json.Unmarshal(data, &session) == nil {
				total++
				if session.SkillActivated {
					activations++
				}
			}
		}
	}

	if total == 0 {
		return 0
	}
	return float64(activations) / float64(total)
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func printSelectUsage() {
	fmt.Fprintf(os.Stderr, "usage: select-best --workspace <eval-workspace-dir> [--candidates <names>]\n")
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	fmt.Fprintf(os.Stderr, "  --workspace    Path to eval workspace directory\n")
	fmt.Fprintf(os.Stderr, "  --candidates   Comma-separated candidate names (default: original,revision-A,revision-B,revision-C)\n")
	fmt.Fprintf(os.Stderr, "  --help, -h     Show this help\n")
	fmt.Fprintf(os.Stderr, "\nSelects the best skill variant based on composite scoring with threshold checks.\n")
	fmt.Fprintf(os.Stderr, "\nexit codes: 0=success, 1=prep error, 2=no candidate passed thresholds, 3=I/O error\n")
}
