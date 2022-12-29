GO?=go

GOOS?=$(shell $(GO) env GOOS)
GOARCH?=$(shell $(GO) env GOARCH)
export GOOS
export GOARCH

ifeq ($(GOOS),windows)
GOFLAGS+=-ldflags -H=windowsgui
endif

.PHONY: build
build:  ## build an executable
	$(GO) build $(GOFLAGS) .

.PHONY: run2
run2:  ## run with test data
	$(GO) run . -f CPM2.xlsx

.PHONY: run
run:  ## run with test data
	$(GO) run . -f CPM.xlsx

.PHONY: fmt
fmt:  ## format Go code
	$(GO) fmt ./...
