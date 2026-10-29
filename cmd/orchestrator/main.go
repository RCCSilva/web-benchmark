// Command orchestrator builds every benchmark app's Docker image, runs each
// one (plus a Postgres instance) on Kubernetes, load-tests it, and records
// the results — the Go replacement for the previous src/index.ts.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"

	"github.com/RCCSilva/web-benchmark/internal/appspec"
	"github.com/RCCSilva/web-benchmark/internal/config"
	"github.com/RCCSilva/web-benchmark/internal/imagebuild"
	"github.com/RCCSilva/web-benchmark/internal/k8sclient"
	"github.com/RCCSilva/web-benchmark/internal/runner"
)

const appsDir = "apps"

func main() {
	_ = godotenv.Load()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	k8s, err := k8sclient.New(cfg.K8sNamespace)
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cleanupDone := make(chan struct{})
	go func() {
		<-ctx.Done()
		fmt.Println()
		if err := k8s.Cleanup(context.Background()); err != nil {
			fmt.Println("Cleanup after interrupt failed:", err)
		}
		close(cleanupDone)
	}()

	err = run(ctx, cfg, k8s)
	if ctx.Err() != nil {
		// Interrupted: wait for the goroutine above to finish tearing down
		// pods/services before exiting, instead of racing it.
		<-cleanupDone
		os.Exit(1)
	}
	if err != nil {
		log.Fatal("Failed: ", err)
	}
	fmt.Println("Finished!")
}

func run(ctx context.Context, cfg config.Config, k8s *k8sclient.Client) error {
	apps, err := appspec.Load(appsDir, cfg.ImagePrefix)
	if err != nil {
		return err
	}

	if err := k8s.Cleanup(ctx); err != nil {
		return err
	}

	if err := imagebuild.BuildAndUploadImage(ctx, apps); err != nil {
		return err
	}

	r := runner.New(k8s, cfg)
	if err := r.RunAll(ctx, apps); err != nil {
		return err
	}

	return k8s.Cleanup(ctx)
}
