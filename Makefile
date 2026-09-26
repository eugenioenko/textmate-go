.PHONY: test bench bench-compare lint

test:
	CGO_ENABLED=0 go test ./...

bench:
	CGO_ENABLED=0 go test -run '^$$' -bench . -benchmem ./...

# Compares textmate-go with Chroma and vscode-textmate; see benchmarks/README.md.
bench-compare:
	node benchmarks/compare.mjs

lint:
	test -z "$$(gofmt -l $$(rg --files -g '*.go' -g '!conformance/corpus/**'))"
	CGO_ENABLED=0 go vet ./...
	golangci-lint run ./...
