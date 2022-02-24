SHELL := /bin/bash -o pipefail

UNAME_OS := $(shell uname -s)
UNAME_ARCH := $(shell uname -m)

.PHONY: proto
protogen:
	buf generate proto

lint:
	revive -config lintconfig.toml -formatter friendly chat/... echo/...

tidy:
	go mod tidy -go=1.17