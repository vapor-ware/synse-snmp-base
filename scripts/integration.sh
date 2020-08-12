#!/usr/bin/env bash
#
# integration.sh
#
# Integration test runner script. This should be run from the project root
# (./scripts/integration.sh) or via Make target (make integration-test).
#
# This script starts an emulator container and runs the integration tests.
# Once the tests complete, it will clean up the container(s).
#

docker-compose -f configs/integration-test.yaml up -d
sleep 2

go test -run Integration -coverprofile=coverage.out -covermode=atomic ./...
rc=$?

if [[ ${rc} -ne 0 ]]; then
    echo "-------------------------------------------------------------------"
    echo "Error: Integration tests failed. The emulator container is cleaned"
    echo "  up by default. If you wish to inspect the logs, disable container"
    echo "  removal in ./scripts/integration.sh"
    echo "-------------------------------------------------------------------"
fi

docker-compose -f configs/integration-test.yaml rm --force --stop

exit ${rc}
