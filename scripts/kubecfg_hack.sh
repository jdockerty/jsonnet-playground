#!/bin/sh

ACTUAL_KUBECFG_VERSION=$(go mod edit -json | jq -r '.Require[] | select(.Path | contains("kubecfg/kubecfg")).Version' | tr -d " ")
WRITTEN_KUBECFG_VERSION=$(grep "KubecfgVersion  " internal/server/backend_routes.go | cut -d "=" -f 2 | tr -d '" ')

if [ "${ACTUAL_KUBECFG_VERSION}" != "${WRITTEN_KUBECFG_VERSION}" ]; then
    echo "Written kubecfg version (${WRITTEN_KUBECFG_VERSION}) does not match go.mod file version (${ACTUAL_KUBECFG_VERSION})"
    exit 1
else
    echo "kubecfg versions match"
fi
