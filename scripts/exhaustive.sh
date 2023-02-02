#!/usr/bin/env bash

echo "==> Checking for switch statement exhaustiveness..."

go run github.com/nishanths/exhaustive/cmd/exhaustive ./...
