GIT_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
BUILD_DATE := $(shell date "+%Y-%m-%d:%H:%M")
GIT_TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
ifeq ($(GIT_TAG),)
VERSION := ${GIT_COMMIT}-${BUILD_DATE}
else
VERSION := $(GIT_TAG:v%=%)-${GIT_COMMIT}-${BUILD_DATE}
endif

build:
	go mod tidy
	go build -ldflags "-X main.Version=${VERSION}" -o generator

run:
	go run generator.go -name=Foo -type=string > string.go

clean:
	rm generator

test:
	go test -v ./...


