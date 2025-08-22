#!/usr/bin/env bash

# Purpose: This script runs the component tests.

set -eu

TEST_OUTPUT_DIR=${TEST_OUTPUT_DIR:-""}
QLTY_COVERAGE_TOKEN=${QLTY_COVERAGE_TOKEN:-""} # Used for selection report type

GO_PACKAGES=$(go list ./test/component)

# Will display some colours in the output if ´gotest´ is installed
GOTEST=$(which gotest > /dev/null 2>&1 && echo "gotest" || echo "go test")

if [ "$TEST_OUTPUT_DIR" == "" ]
then
    # This block handles the default local component test running

    # The idiomatic way to disable test caching explicitly is to use -count=1
    $GOTEST -count=1 -race -cover $GO_PACKAGES

elif [ "$QLTY_COVERAGE_TOKEN" == "" ]
then
    # This block handles the local component tests if junit output is desired

    mkdir -p "$TEST_OUTPUT_DIR/component"
    gotestsum --junitfile "$TEST_OUTPUT_DIR/component/report.xml" -- -count=1 -race -cover $GO_PACKAGES
else
    # This block handles the component tests in CI. It does junit output. It also doesn't stop after the first failure.
    # It'll continue so that it can report a full list of failures in circle

    test_exit_code=0
    test_fail=0

    for pkg in $GO_PACKAGES; do \
        PKG_BASENAME=$(basename -- "$pkg")
        mkdir -p "$TEST_OUTPUT_DIR/$PKG_BASENAME"
        gotestsum --format testname --junitfile "$TEST_OUTPUT_DIR/component/$PKG_BASENAME.xml" -- -v -race -coverpkg="github.com/qlik-trial/usage-telemetry-publisher/internal/..." -json -coverprofile=$(echo $pkg | tr / -).out $pkg || test_exit_code=$? ; \
        if [ $test_exit_code -ne 0 ] ; then
            test_fail=1
        fi
    done

    if [ $test_fail -ne 0 ] ; then
        echo "Component tests failed "
        exit 1
    fi
fi
