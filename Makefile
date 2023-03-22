GIT_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
BUILD_DATE := $(shell date "+%Y-%m-%d:%H:%M")
GIT_TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
ifeq ($(GIT_TAG),)
VERSION := ${GIT_COMMIT}-${BUILD_DATE}
else
VERSION := $(GIT_TAG:v%=%)-${GIT_COMMIT}-${BUILD_DATE}
endif

build: core/primitives.auto.go
	#go mod tidy
	go build -ldflags "-X main.Version=${VERSION}"

install: build
	go install -ldflags "-X main.Version=${VERSION}" .

test:
	go test -v -coverprofile coverage.out ./...

test-stream:
	go test -v ./core --run TestMachine

coverage: test
	go tool cover -html=coverage.out

gen:
	go run _generator/generator.go -name=Foo -type=string > string.go


release:
	goreleaser release --rm-dist

core/primitives.auto.go: core/primitives.go.tmpl
	go run _helpers/generator.go --template $< --out $@

