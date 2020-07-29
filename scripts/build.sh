#!/bin/bash

set -e

VERSION=$(cat VERSION)

echo "Building version $VERSION"
go build -ldflags "-X github.com/olegsu/kubectl-fetch-yaml/cmd.version=$VERSION" -o kubectl-fetch-yaml *.go