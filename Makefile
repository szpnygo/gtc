cat := $(if $(filter $(OS),Windows_NT),type,cat)
version := $(shell $(cat) version)

.PHONY: vet
vet:
	@go vet $(go list ./...)

.PHONY: lint
lint:
	@golangci-lint run

build:
	go build -ldflags="-s -w" -ldflags="-X 'main.BuildTime=$(version)'" -o gtc main.go
	$(if $(shell command -v upx), upx gtc)