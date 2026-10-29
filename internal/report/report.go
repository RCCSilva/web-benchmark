// Package report writes the RESULT.md table and archives it into results/,
// mirroring the previous src/result.ts.
package report

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/RCCSilva/web-benchmark/internal/appspec"
	"github.com/RCCSilva/web-benchmark/internal/loadtest"
)

const header = `| Language | Packages | Database | Average Latency (ms) | Median Latency (ms) | Average Requests/s | Median Requests/s |Success Rate |
|----------|-----------------|----------|------------:|------------:|-------------:| -----------:| -------------:|
`

const resultFile = "RESULT.md"

// BuildResultHeader (re)creates RESULT.md with just the table header.
func BuildResultHeader() error {
	return os.WriteFile(resultFile, []byte(header), 0o644)
}

// ProcessResult appends one row describing app's benchmark result.
func ProcessResult(app appspec.BenchmarkApp, result loadtest.Result) error {
	packageNames := make([]string, len(app.Packages))
	for i, p := range app.Packages {
		packageNames[i] = p.Name
	}

	hasDatabase := "✗"
	if app.Database {
		hasDatabase = "✓"
	}

	successRate := 100.0
	if result.TotalRequests > 0 {
		successRate = (1 - float64(result.Non2xx)/float64(result.TotalRequests)) * 100
	}

	row := fmt.Sprintf(
		"| %s (%s) | %s | %s | %.2f | %.2f | %.2f | %.2f | %.2f%% | \n",
		app.Language.Name, app.Language.Version,
		strings.Join(packageNames, " / "),
		hasDatabase,
		result.AverageLatencyMs,
		result.MedianLatencyMs,
		result.AverageRequestsPerSec,
		result.MedianRequestsPerSec,
		successRate,
	)

	f, err := os.OpenFile(resultFile, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(row)
	return err
}

// CopyResultToResults copies RESULT.md into results/<timestamp>_RESULT.md.
func CopyResultToResults() error {
	src, err := os.Open(resultFile)
	if err != nil {
		return err
	}
	defer src.Close()

	name := time.Now().UTC().Format("2006-01-02T15:04:05") + "_RESULT.md"
	dst, err := os.Create(filepath.Join("results", name))
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}
