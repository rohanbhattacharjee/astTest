GOFMT_FILES?=$$(find . -name '*.go')
PACKAGE=main

default: fmt build

build: ;go build -o astTest

fmt:
	goimports -w -local ${PACKAGE} $(GOFMT_FILES)
