# Copyright 2020 American Express Travel Related Services Company, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except
# in compliance with the License. You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software distributed under the License
# is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express
# or implied. See the License for the specific language governing permissions and limitations under
# the License.

.PHONY: all clean tests lint build format benchmarks benchmarks-race coverage

all: build tests lint

tests:
	@echo "Running Tests with Coverage Report"
	go test -race ./...

benchmarks:
	@echo "Running Benchmarks"
	go test -run=^$$ -bench=. -benchmem ./...

benchmarks-race:
	@echo "Running Benchmarks with Race Detection"
	go test -race -run=^$$ -bench=. -benchmem ./...

build:
	@echo "Building package"
	go build ./...

format:
	@echo "Formatting code"
	gofmt -s -w .
	@if command -v golines >/dev/null 2>&1; then \
		golines -w .; \
	else \
		echo "golines not installed, skipping line wrapping"; \
	fi

lint:
	@echo "Linting code"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run --enable-only=misspell --enable-only=revive --enable-only=errname ./...; \
	else \
		echo "golangci-lint not installed, skipping lint"; \
	fi

coverage:
	@echo "Running coverage"
	go test -race -covermode=atomic -coverprofile=coverage.out ./...
	go tool cover -html=./coverage.out -o ./coverage.html

clean:
	@echo "Cleaning build artifacts"
	@find . -type f -name "*.test" -delete
	@find . -type f -name "coverage.out" -delete
	@find . -type f -name "coverage.html" -delete
	@rm -rf bin/
