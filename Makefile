GO?=go

.PHONY: run
run:  ## run with test data
	$(GO) run . -f CPM.xlsx

.PHONY: fmt
fmt:  ## format Go code
	$(GO) fmt ./...
