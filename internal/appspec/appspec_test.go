package appspec

import (
	"path/filepath"
	"testing"
)

func TestLoadRealApps(t *testing.T) {
	apps, err := Load(filepath.Join("..", "..", "apps"), "")
	if err != nil {
		t.Fatalf("Load: %v", err)
	}

	if len(apps) == 0 {
		t.Fatal("expected at least one app to load from apps/")
	}

	seenSlugs := map[string]bool{}
	for _, app := range apps {
		if app.Slug == "" {
			t.Errorf("app %q has empty slug", app.Path)
		}
		if app.Tag == "" {
			t.Errorf("app %q has empty tag", app.Path)
		}
		if app.Image == "" {
			t.Errorf("app %q has empty image", app.Path)
		}
		if app.Language.Name == "" {
			t.Errorf("app %q has empty language name", app.Path)
		}
		seenSlugs[app.Slug] = true
	}

	if len(seenSlugs) != len(apps) {
		t.Errorf("expected unique slugs per app, got %d apps but %d unique slugs", len(apps), len(seenSlugs))
	}
}
