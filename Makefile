GO?=go

.PHONY: build
build:  ## build an executable
	$(GO) build .

.PHONY: run2
run2:  ## run with test data
	$(GO) run . -f CPM2.xlsx

.PHONY: run
run:  ## run with test data
	$(GO) run . -f CPM.xlsx

.PHONY: fmt
fmt:  ## format Go code
	$(GO) fmt ./...
