#!/bin/sh
if [ "$IS_CONTAINER" != "" ]; then
  for TARGET in "${@}"; do
    find "${TARGET}" -name '*.go' ! -path '*/vendor/*' ! -path '*/.build/*' -exec gofmt -s -w {} \+
  done
  git diff --exit-code
else
  docker run -it --rm \
    --env IS_CONTAINER=TRUE \
    --volume "${PWD}:/go/src/github.com/openshift/cluster-api-provider-libvirt:z" \
    --workdir /go/src/github.com/openshift/cluster-api-provider-libvirt \
    openshift/origin-release:golang-1.10 \
    ./hack/go-fmt.sh "${@}"
fi
