// Package runner wires together database + app provisioning and the load
// test for each benchmark app, mirroring the previous
// src/runner/{run-test,create-app,create-database}.ts.
package runner

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/RCCSilva/web-benchmark/internal/appspec"
	"github.com/RCCSilva/web-benchmark/internal/config"
	"github.com/RCCSilva/web-benchmark/internal/dbwait"
	"github.com/RCCSilva/web-benchmark/internal/k8sclient"
	"github.com/RCCSilva/web-benchmark/internal/loadtest"
	"github.com/RCCSilva/web-benchmark/internal/report"
)

const loadTestConnections = 50

type Runner struct {
	K8s *k8sclient.Client
	Cfg config.Config
}

func New(k8s *k8sclient.Client, cfg config.Config) *Runner {
	return &Runner{K8s: k8s, Cfg: cfg}
}

// RunAll benchmarks every app in sequence, matching the previous one-at-a-
// time behaviour so only one pod/database pair occupies the cluster at once.
func (r *Runner) RunAll(ctx context.Context, apps []appspec.BenchmarkApp) error {
	if err := report.BuildResultHeader(); err != nil {
		return fmt.Errorf("building result header: %w", err)
	}

	for _, app := range apps {
		fmt.Printf("Testing %s\n", app.Name)
		if err := r.runOne(ctx, app); err != nil {
			return fmt.Errorf("running benchmark for %s: %w", app.Name, err)
		}
		if err := sleep(ctx, 10*time.Second); err != nil {
			return err
		}
	}

	return report.CopyResultToResults()
}

func (r *Runner) runOne(ctx context.Context, app appspec.BenchmarkApp) error {
	fmt.Println("   Creating database...")
	db, err := r.K8s.CreateDatabase(ctx, app.Slug, r.Cfg.PostgresUser, r.Cfg.PostgresPass, r.Cfg.PostgresDB)
	if err != nil {
		return fmt.Errorf("creating database: %w", err)
	}

	if err := dbwait.WaitAndCreateTable(ctx, dbwait.Conn{
		User:     r.Cfg.PostgresUser,
		Password: r.Cfg.PostgresPass,
		Host:     r.Cfg.ClusterIP,
		Port:     db.NodePort,
		Name:     r.Cfg.PostgresDB,
	}, r.Cfg.Debug); err != nil {
		return fmt.Errorf("waiting for database: %w", err)
	}

	fmt.Println("   Creating app...")
	appResult, err := r.K8s.CreateApp(ctx, k8sclient.AppSpec{
		Slug:         app.Slug,
		Image:        app.Image,
		DBHost:       db.ServiceName,
		PostgresUser: r.Cfg.PostgresUser,
		PostgresPass: r.Cfg.PostgresPass,
		PostgresDB:   r.Cfg.PostgresDB,
	})
	if err != nil {
		return fmt.Errorf("creating app: %w", err)
	}

	url := fmt.Sprintf("%s:%d", r.Cfg.ClusterHost, appResult.NodePort)
	if err := verifyAppReady(ctx, url, r.Cfg.Debug); err != nil {
		return fmt.Errorf("waiting for app: %w", err)
	}

	duration, err := parseTestDuration(r.Cfg.TestDuration)
	if err != nil {
		return err
	}

	if !r.Cfg.Test {
		fmt.Println("   Warming up...")
		if _, err := loadtest.Run(ctx, loadtest.Config{URL: url, Connections: loadTestConnections, Duration: duration}); err != nil {
			return fmt.Errorf("warm-up: %w", err)
		}
	}

	fmt.Println("   Running benchmark...")
	result, err := loadtest.Run(ctx, loadtest.Config{URL: url, Connections: loadTestConnections, Duration: duration})
	if err != nil {
		return fmt.Errorf("benchmark: %w", err)
	}

	if err := report.ProcessResult(app, result); err != nil {
		return fmt.Errorf("recording result: %w", err)
	}

	appErr := appResult.Delete(ctx)
	dbErr := db.Delete(ctx)
	if appErr != nil {
		return appErr
	}
	return dbErr
}

func verifyAppReady(ctx context.Context, url string, debug bool) error {
	for {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return err
		}
		resp, err := http.DefaultClient.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode < 500 {
				return nil
			}
		}
		if debug && err != nil {
			fmt.Printf("   -> %v\n", err)
		}
		if err := sleep(ctx, 500*time.Millisecond); err != nil {
			return err
		}
	}
}

// sleep waits for d, returning early with ctx.Err() if ctx is cancelled.
func sleep(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

func parseTestDuration(seconds string) (time.Duration, error) {
	n, err := strconv.Atoi(seconds)
	if err != nil {
		return 0, fmt.Errorf("invalid TEST_DURATION %q: %w", seconds, err)
	}
	return time.Duration(n) * time.Second, nil
}
