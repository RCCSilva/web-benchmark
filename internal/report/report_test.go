package report

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/RCCSilva/web-benchmark/internal/appspec"
	"github.com/RCCSilva/web-benchmark/internal/loadtest"
)

const wantHeader = `| Language | Packages | Database | Average Latency (ms) | Median Latency (ms) | Average Requests/s | Median Requests/s |Success Rate |
|----------|-----------------|----------|------------:|------------:|-------------:| -----------:| -------------:|
`

func TestBuildResultHeaderMatchesExistingFormat(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if err := os.Mkdir("results", 0o755); err != nil {
		t.Fatal(err)
	}

	if err := BuildResultHeader(); err != nil {
		t.Fatalf("BuildResultHeader: %v", err)
	}

	got, err := os.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != wantHeader {
		t.Errorf("header mismatch:\ngot:  %q\nwant: %q", got, wantHeader)
	}
}

func TestProcessResultAndCopy(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	if err := os.Mkdir("results", 0o755); err != nil {
		t.Fatal(err)
	}
	if err := BuildResultHeader(); err != nil {
		t.Fatal(err)
	}

	app := appspec.BenchmarkApp{
		Language: appspec.Language{Name: "Go", Version: "1.22"},
		Packages: []appspec.Package{{Name: "http", Version: "1.22"}},
		Database: false,
	}
	result := loadtest.Result{
		AverageLatencyMs:      4.15,
		MedianLatencyMs:       3.0,
		AverageRequestsPerSec: 106628.80,
		MedianRequestsPerSec:  107007.0,
		TotalRequests:         1000,
		Non2xx:                0,
	}

	if err := ProcessResult(app, result); err != nil {
		t.Fatalf("ProcessResult: %v", err)
	}

	contents, err := os.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}

	wantRow := "| Go (1.22) | http | ✗ | 4.15 | 3.00 | 106628.80 | 107007.00 | 100.00% | \n"
	if string(contents) != wantHeader+wantRow {
		t.Errorf("unexpected RESULT.md contents:\n%s", contents)
	}

	if err := CopyResultToResults(); err != nil {
		t.Fatalf("CopyResultToResults: %v", err)
	}

	entries, err := os.ReadDir("results")
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 archived result file, got %d", len(entries))
	}
	if filepath.Ext(entries[0].Name()) != ".md" {
		t.Errorf("expected .md archived file, got %s", entries[0].Name())
	}
}
