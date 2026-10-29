#!/bin/bash

set -e

eval $(minikube docker-env)
go build -o bin/orchestrator ./cmd/orchestrator
exec ./bin/orchestrator
