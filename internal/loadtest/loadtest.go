// Package loadtest is a minimal replacement for the previous autocannon
// (Node) dependency: a goroutine worker pool hammers a URL for a fixed
// duration and reports the same shape of numbers result.ts already expects
// (mean/p50 latency, mean/p50 requests-per-second, success rate).
//
// Note: because the underlying load generator changed language/runtime,
// absolute numbers from this won't be bit-for-bit comparable to older
// autocannon-generated results.
package loadtest

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type Config struct {
	URL         string
	Connections int
	Duration    time.Duration
}

type Result struct {
	AverageLatencyMs      float64
	MedianLatencyMs       float64
	AverageRequestsPerSec float64
	MedianRequestsPerSec  float64
	TotalRequests         int64
	Non2xx                int64
}

// Run drives Config.Connections concurrent workers against Config.URL for
// Config.Duration and reports aggregate latency/throughput/success stats.
func Run(ctx context.Context, cfg Config) (Result, error) {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: cfg.Connections,
		},
		Timeout: 10 * time.Second,
	}

	start := time.Now()
	deadline := start.Add(cfg.Duration)
	secondBuckets := int(cfg.Duration/time.Second) + 2

	var (
		mu        sync.Mutex
		latencies = make([]time.Duration, 0, 1024)
		perSecond = make([]int64, secondBuckets)
		totalReqs int64
		non2xx    int64
	)

	var wg sync.WaitGroup
	for i := 0; i < cfg.Connections; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				if ctx.Err() != nil || time.Now().After(deadline) {
					return
				}

				reqStart := time.Now()
				resp, err := client.Get(cfg.URL)
				elapsed := time.Since(reqStart)

				atomic.AddInt64(&totalReqs, 1)
				bucket := int(time.Since(start) / time.Second)
				if bucket >= 0 && bucket < secondBuckets {
					atomic.AddInt64(&perSecond[bucket], 1)
				}

				if err != nil || resp.StatusCode < 200 || resp.StatusCode >= 300 {
					atomic.AddInt64(&non2xx, 1)
				}
				if err == nil {
					resp.Body.Close()
				}

				mu.Lock()
				latencies = append(latencies, elapsed)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	if totalReqs == 0 {
		return Result{}, fmt.Errorf("no requests completed against %s", cfg.URL)
	}

	avgLatency, medianLatency := latencyStats(latencies)
	avgRPS, medianRPS := throughputStats(perSecond, cfg.Duration)

	return Result{
		AverageLatencyMs:      avgLatency,
		MedianLatencyMs:       medianLatency,
		AverageRequestsPerSec: avgRPS,
		MedianRequestsPerSec:  medianRPS,
		TotalRequests:         totalReqs,
		Non2xx:                non2xx,
	}, nil
}

func latencyStats(latencies []time.Duration) (avgMs, medianMs float64) {
	if len(latencies) == 0 {
		return 0, 0
	}

	var sum time.Duration
	for _, l := range latencies {
		sum += l
	}
	avgMs = float64(sum.Microseconds()) / 1000 / float64(len(latencies))

	sorted := append([]time.Duration(nil), latencies...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	medianMs = float64(sorted[len(sorted)/2].Microseconds()) / 1000

	return avgMs, medianMs
}

func throughputStats(perSecond []int64, duration time.Duration) (avg, median float64) {
	fullSeconds := int(duration / time.Second)
	if fullSeconds == 0 || fullSeconds > len(perSecond) {
		fullSeconds = len(perSecond)
	}
	counts := perSecond[:fullSeconds]
	if len(counts) == 0 {
		return 0, 0
	}

	var sum int64
	for _, c := range counts {
		sum += c
	}
	avg = float64(sum) / float64(len(counts))

	sorted := append([]int64(nil), counts...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	median = float64(sorted[len(sorted)/2])

	return avg, median
}
