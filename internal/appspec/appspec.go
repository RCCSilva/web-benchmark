// Package appspec loads and validates the apps/*/benchmark.json files that
// describe each benchmark target, mirroring the previous src/apps.ts.
package appspec

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Language struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Package struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type benchmarkFile struct {
	Language Language  `json:"language"`
	Packages []Package `json:"packages"`
	Database bool      `json:"database"`
}

type BenchmarkApp struct {
	Tag      string
	Image    string
	Path     string
	Slug     string
	Name     string
	Database bool
	Language Language
	Packages []Package
}

// Load reads every apps/*/benchmark.json under appsDir, skipping any
// directory whose benchmark.json is missing, malformed, or incomplete.
func Load(appsDir, imagePrefix string) ([]BenchmarkApp, error) {
	entries, err := os.ReadDir(appsDir)
	if err != nil {
		return nil, fmt.Errorf("reading apps dir: %w", err)
	}

	tagVersion, err := randomHex(5)
	if err != nil {
		return nil, err
	}

	var apps []BenchmarkApp
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		app, ok := loadOne(filepath.Join(appsDir, entry.Name()), tagVersion, imagePrefix)
		if ok {
			apps = append(apps, app)
		}
	}

	return apps, nil
}

func loadOne(appPath, tagVersion, imagePrefix string) (BenchmarkApp, bool) {
	data, err := os.ReadFile(filepath.Join(appPath, "benchmark.json"))
	if err != nil {
		return BenchmarkApp{}, false
	}

	var spec benchmarkFile
	if err := json.Unmarshal(data, &spec); err != nil {
		return BenchmarkApp{}, false
	}
	if spec.Language.Name == "" || len(spec.Packages) == 0 {
		return BenchmarkApp{}, false
	}

	packageNames := make([]string, len(spec.Packages))
	packageLabels := make([]string, len(spec.Packages))
	for i, p := range spec.Packages {
		packageNames[i] = strings.ToLower(p.Name)
		packageLabels[i] = fmt.Sprintf("%s (%s)", p.Name, p.Version)
	}

	slug := strings.Join(packageNames, "")
	tag := strings.ToLower(fmt.Sprintf("web-benchmark-ecr:%s%s", slug, tagVersion))

	return BenchmarkApp{
		Tag:      tag,
		Image:    imagePrefix + tag,
		Path:     appPath,
		Slug:     slug,
		Name:     fmt.Sprintf("%s (%s) / %s", spec.Language.Name, spec.Language.Version, strings.Join(packageLabels, " / ")),
		Database: spec.Database,
		Language: spec.Language,
		Packages: spec.Packages,
	}, true
}

func randomHex(n int) (string, error) {
	buf := make([]byte, n)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
