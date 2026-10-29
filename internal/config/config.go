// Package config exposes the environment-derived settings the orchestrator
// needs, mirroring the previous src/constants.ts.
package config

import (
	"fmt"
	"os"
	"regexp"
)

var schemeRegexp = regexp.MustCompile(`^https?://`)

type Config struct {
	ClusterHost  string
	ClusterIP    string
	Debug        bool
	Test         bool
	TestDuration string
	PostgresUser string
	PostgresPass string
	PostgresDB   string
	K8sNamespace string
	ImagePrefix  string
}

func Load() (Config, error) {
	clusterHost := os.Getenv("CLUSTER_HOST")
	if clusterHost == "" {
		return Config{}, fmt.Errorf("CLUSTER_HOST not defined")
	}

	test := os.Getenv("TEST") == "1"
	testDuration := "40"
	if test {
		testDuration = "5"
	}

	return Config{
		ClusterHost:  clusterHost,
		ClusterIP:    schemeRegexp.ReplaceAllString(clusterHost, ""),
		Debug:        os.Getenv("DEBUG") == "1",
		Test:         test,
		TestDuration: testDuration,
		PostgresUser: "user",
		PostgresPass: "Password!234",
		PostgresDB:   "mydb",
		K8sNamespace: "default",
		ImagePrefix:  os.Getenv("IMAGE_PREFIX"),
	}, nil
}
