all:
	go build -o bin/stratus cmd/stratus/*.go

test:
	go test ./... -v

mocks:
	mockery --name=StateManager --dir internal/state --output internal/state/mocks
	mockery --name=FileSystem --structname FileSystemMock --dir internal/state --output internal/state/mocks