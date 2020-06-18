#!/bin/bash

shopt -s globstar
shopt -s nullglob

echo "Running gofmt"
FMT_DIFF=`gofmt -s -d **/*.go`
if [ "$FMT_DIFF" != "" ]
then
    echo "Go fmt found error(s)"
    printf '%s\n' "$FMT_DIFF"
    exit 1
fi

echo "Running goimports"
IMPORTS_DIFF=`goimports -d **/*.go`
if [ "$IMPORTS_DIFF" != "" ]
then
    echo "Goimports found error(s)"
    printf '%s\n' "$IMPORTS_DIFF"
    exit 1
fi

echo "Running go vet"
go vet ./... || exit 1

# go install golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow
echo "Running shadow"
go vet -vettool=$(which shadow) ./... || exit 1

# go get -u github.com/mgechev/revive
echo "Running revive"
revive -config ./revive.toml -formatter stylish ./... || exit 1

# go get -v github.com/go-lintpack/lintpack/...
# go get -v github.com/go-critic/go-critic/...
# cd $(go env GOPATH)/src/github.com/go-critic/go-critic && make gocritic
echo "Running gocritic"
gocritic check ./... || exit 1

# https://staticcheck.io/docs/
echo "Running staticcheck"
staticcheck ./... || exit 1

# https://github.com/kisielk/errcheck
echo "Running errcheck"
errcheck -asserts ./... || exit 1

# https://github.com/mdempsky/unconvert
echo "Running unconvert"
unconvert -v ./... || exit 1

# https://github.com/client9/misspell
misspell -error **/*.{go,md} || exit 1