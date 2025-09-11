godeps:
	@go mod tidy

gobuild: godeps
	@go build -o bin/lgp_validator cmd/validator/main.go

gotest: godeps
	@go test -v --cover ./...
