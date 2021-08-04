# Makefile is used to drive testing and build

tests:
	@echo "Running Tests with Coverage Report"
	go test -v -covermode=count -coverprofile=coverage.out ./...
	go tool cover -html=./coverage.out -o ./coverage.html

benchmarks:
	@echo "Running Benchmarks"
	go test -run=Bench -count=10 -bench . ./...
