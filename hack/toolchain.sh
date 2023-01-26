#!/usr/bin/env bash
set -euo pipefail

go install github.com/google/go-licenses@v1.5.0
go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.1
go install golang.org/x/vuln/cmd/govulncheck@v0.0.0-20221215205010-9bf256343acc
