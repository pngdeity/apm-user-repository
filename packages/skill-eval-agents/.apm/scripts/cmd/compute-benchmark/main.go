package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type gradingResult struct {
	Pass   bool    `json:"pass"`
	Assert string  `json:"assert,omitempty"`
	Score  float64 `json:"score,omitempty"`
}

type timingResult struct {
	DurationMs int64   `json:"duration_ms"`
	DurationS  float64 `json:"duration_s,omitempty"`
	TokenCount int     `json:"token_count,omitempty"`
}

type stats struct {
	Mean   float64 `json:"mean"`
	Stddev float64 `json:"stddev"`
}

type configStats struct {
	PassRate    stats `json:"pass_rate"`
	TimeSeconds stats `json:"time_seconds"`
	Tokens      stats `json:"tokens"`
}

type benchmarkOutput struct {
	RunSummary struct {
		WithSkill    configStats `json:"with_skill"`
		WithoutSkill configStats `json:"without_skill"`
		Delta        struct {
			PassRate    float64 `json:"pass_rate"`
			TimeSeconds float64 `json:"time_seconds"`
			Tokens      float64 `json:"tokens"`
		} `json:"delta"`
	} `json:"run_summary"`
}

func main() {
	os.Exit(run())
}

func run() int {
	for _, a := range os.Args[1:] {
		if a == "--help" || a == "-h" {
			printComputeUsage()
			return 0
		}
	}

	workspace := flag.String("workspace", "", "Path to eval workspace directory")
	outputPath := flag.String("output", "benchmark.json", "Output JSON file path (relative to workspace)")
	flag.Parse()

	if *workspace == "" {
		printComputeUsage()
		return 1
	}

	workspaceDir := *workspace

	withSkillGrades, withSkillTimings := collectResults(filepath.Join(workspaceDir, "with_skill"))
	withoutSkillGrades, withoutSkillTimings := collectResults(filepath.Join(workspaceDir, "without_skill"))

	fmt.Fprintf(os.Stderr, "found %d grading files (with_skill: %d, without_skill: %d), %d timing files (with_skill: %d, without_skill: %d)\n",
		len(withSkillGrades)+len(withoutSkillGrades), len(withSkillGrades), len(withoutSkillGrades),
		len(withSkillTimings)+len(withoutSkillTimings), len(withSkillTimings), len(withoutSkillTimings))

	output := benchmarkOutput{}

	output.RunSummary.WithSkill = computeConfigStats(withSkillGrades, withSkillTimings)
	output.RunSummary.WithoutSkill = computeConfigStats(withoutSkillGrades, withoutSkillTimings)

	output.RunSummary.Delta.PassRate = output.RunSummary.WithSkill.PassRate.Mean - output.RunSummary.WithoutSkill.PassRate.Mean
	output.RunSummary.Delta.TimeSeconds = output.RunSummary.WithSkill.TimeSeconds.Mean - output.RunSummary.WithoutSkill.TimeSeconds.Mean
	output.RunSummary.Delta.Tokens = output.RunSummary.WithSkill.Tokens.Mean - output.RunSummary.WithoutSkill.Tokens.Mean

	result, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error marshaling benchmark: %v\n", err)
		return 3
	}

	benchmarkPath := *outputPath
	if !filepath.IsAbs(benchmarkPath) {
		benchmarkPath = filepath.Join(workspaceDir, benchmarkPath)
	}
	if err := os.WriteFile(benchmarkPath, result, 0644); err != nil {
		fmt.Fprintf(os.Stderr, "error writing benchmark: %v\n", err)
		return 3
	}

	fmt.Println(string(result))
	return 0
}

func printComputeUsage() {
	fmt.Fprintf(os.Stderr, "usage: compute-benchmark --workspace <eval-workspace-dir> [--output <path>]\n")
	fmt.Fprintf(os.Stderr, "\nFlags:\n")
	fmt.Fprintf(os.Stderr, "  --workspace   Path to eval workspace directory\n")
	fmt.Fprintf(os.Stderr, "  --output      Output JSON file path (default: benchmark.json relative to workspace)\n")
	fmt.Fprintf(os.Stderr, "  --help, -h    Show this help\n")
	fmt.Fprintf(os.Stderr, "\nAggregates grading.json and timing.json files into a benchmark JSON.\n")
	fmt.Fprintf(os.Stderr, "\nexit codes: 0=success, 1=prep error, 2=parse error, 3=I/O error\n")
}

func collectResults(dir string) ([]gradingResult, []timingResult) {
	var grades []gradingResult
	var timings []timingResult

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return grades, timings
	}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}
		base := filepath.Base(path)
		if strings.HasSuffix(base, "grading.json") {
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: error reading %s: %v\n", path, err)
				return nil
			}
			var g gradingResult
			if err := json.Unmarshal(data, &g); err != nil {
				fmt.Fprintf(os.Stderr, "warning: error parsing %s: %v\n", path, err)
				return nil
			}
			grades = append(grades, g)
		}
		if strings.HasSuffix(base, "timing.json") {
			data, err := os.ReadFile(path)
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: error reading %s: %v\n", path, err)
				return nil
			}
			var t timingResult
			if err := json.Unmarshal(data, &t); err != nil {
				fmt.Fprintf(os.Stderr, "warning: error parsing %s: %v\n", path, err)
				return nil
			}
			timings = append(timings, t)
		}
		return nil
	})

	return grades, timings
}

func computeConfigStats(grades []gradingResult, timings []timingResult) configStats {
	var cs configStats

	passRates := make([]float64, 0)
	for _, g := range grades {
		if g.Pass {
			passRates = append(passRates, 1.0)
		} else {
			passRates = append(passRates, 0.0)
		}
	}
	if len(passRates) == 0 {
		passRates = append(passRates, math.NaN())
	}
	cs.PassRate = computeStats(passRates)

	times := make([]float64, 0)
	for _, t := range timings {
		if t.DurationS > 0 {
			times = append(times, t.DurationS)
		} else if t.DurationMs > 0 {
			times = append(times, float64(t.DurationMs)/1000.0)
		}
	}
	if len(times) == 0 {
		times = append(times, math.NaN())
	}
	cs.TimeSeconds = computeStats(times)

	tokens := make([]float64, 0)
	for _, t := range timings {
		if t.TokenCount > 0 {
			tokens = append(tokens, float64(t.TokenCount))
		}
	}
	if len(tokens) == 0 {
		tokens = append(tokens, math.NaN())
	}
	cs.Tokens = computeStats(tokens)

	return cs
}

func computeStats(values []float64) stats {
	if len(values) == 0 {
		return stats{Mean: math.NaN(), Stddev: math.NaN()}
	}

	sum := 0.0
	for _, v := range values {
		sum += v
	}
	mean := sum / float64(len(values))

	variance := 0.0
	for _, v := range values {
		diff := v - mean
		variance += diff * diff
	}
	if len(values) > 1 {
		variance /= float64(len(values) - 1)
	}
	stddev := math.Sqrt(variance)

	mean = math.Round(mean*1000) / 1000
	stddev = math.Round(stddev*1000) / 1000

	return stats{Mean: mean, Stddev: stddev}
}
