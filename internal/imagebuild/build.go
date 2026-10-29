// Package imagebuild builds and tags each benchmark app's Docker image using
// the Docker Engine SDK, replacing the previous src/docker.ts which shelled
// out to the `docker` CLI. client.FromEnv reads DOCKER_HOST/TLS env vars the
// same way the CLI does, so `eval $(minikube docker-env)` still applies.
package imagebuild

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/moby/patternmatcher"
	"github.com/moby/patternmatcher/ignorefile"

	"github.com/RCCSilva/web-benchmark/internal/appspec"
)

// BuildAndUploadImage builds and tags the Docker image for every app, in
// order, matching the sequential behaviour of the previous implementation.
func BuildAndUploadImage(ctx context.Context, apps []appspec.BenchmarkApp) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return fmt.Errorf("creating docker client: %w", err)
	}
	defer cli.Close()

	for _, app := range apps {
		fmt.Printf("Building %s\n", app.Name)
		if err := buildOne(ctx, cli, app); err != nil {
			return fmt.Errorf("failed to build %s: %w", app.Tag, err)
		}
	}

	return nil
}

func buildOne(ctx context.Context, cli *client.Client, app appspec.BenchmarkApp) error {
	buildContext, err := tarDirectory(app.Path)
	if err != nil {
		return fmt.Errorf("packing build context: %w", err)
	}

	resp, err := cli.ImageBuild(ctx, buildContext, types.ImageBuildOptions{
		Tags:       []string{app.Tag},
		Dockerfile: "Dockerfile",
	})
	if err != nil {
		return fmt.Errorf("image build request: %w", err)
	}
	defer resp.Body.Close()

	return jsonmessage.DisplayJSONMessagesStream(resp.Body, io.Discard, 0, false, nil)
}

// tarDirectory packs dir into a tar archive suitable for use as a Docker
// build context, honouring .dockerignore the same way `docker build` does.
func tarDirectory(dir string) (io.Reader, error) {
	patterns, err := readDockerignore(dir)
	if err != nil {
		return nil, err
	}

	matcher, err := patternmatcher.New(patterns)
	if err != nil {
		return nil, fmt.Errorf("parsing .dockerignore: %w", err)
	}

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	err = filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(dir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		ignored, err := matcher.Matches(filepath.ToSlash(rel))
		if err != nil {
			return err
		}
		if ignored {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = filepath.ToSlash(rel)

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tw, file)
		return err
	})
	if err != nil {
		return nil, err
	}

	if err := tw.Close(); err != nil {
		return nil, err
	}

	return &buf, nil
}

func readDockerignore(dir string) ([]string, error) {
	file, err := os.Open(filepath.Join(dir, ".dockerignore"))
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return ignorefile.ReadAll(file)
}
