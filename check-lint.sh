#!/bin/bash

SCRIPTDIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

cd ${SCRIPTDIR}/..

go get -u github.com/alecthomas/gometalinter
"${GOPATH}/bin/gometalinter" --install
"${GOPATH}/bin/gometalinter" --vendor --disable-all --enable=vet --enable=goimports --enable=vetshadow --enable=golint --enable=ineffassign --enable=goconst --tests ./...
