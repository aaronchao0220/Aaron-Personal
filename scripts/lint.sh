#!/usr/bin/env sh

VERSION="2.3.0"

go version

if ! "$(go env GOPATH)/bin/golangci-lint" --version 2>/dev/null | grep -q "version $VERSION"; then
	curl -sfL \
		"https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh" |
		sh -s -- -b "$(go env GOPATH)/bin" "v$VERSION"
fi

"$(go env GOPATH)/bin/golangci-lint" run
