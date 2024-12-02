VERSION := $(shell echo $(shell git describe --tags) | sed 's/^v//')
FRONT_DIR = pkg/topology/embed/frontend
LDFLAGS = -X main.version=${VERSION}

build: build-front build-go

build-front:
	pnpm -C $(FRONT_DIR) install
	pnpm -C $(FRONT_DIR) run build

build-go:
	go build -ldflags '$(LDFLAGS)' cmd/tmtop.go

install:
	go install -ldflags '$(LDFLAGS)' cmd/tmtop.go

lint:
	golangci-lint run --fix ./...

test:
	go test -coverprofile cover.out -v ./...
