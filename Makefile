export GO111MODULE=on

TARGET := fluxslack
.DEFAULT_GOAL: $(TARGET)

# test runner (can be overriden by CI)
GOTEST ?= go test

# go source files, ignore vendor directory
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: all build clean install uninstall fmt simplify check deps generate run

all: check build

$(TARGET): $(SRC)
	@go build -o build/$(TARGET) ./cmd/$(TARGET)

build: $(TARGET)
	@true

clean:
	@go clean
	@rm -rf build/

install: build
	@cp build/$(TARGET) $(GOPATH)/bin/

uninstall: clean
	@rm -f $$(which ${TARGET})

fmt:
	@gofmt -l -w $(SRC)

simplify:
	@gofmt -s -l -w $(SRC)

check:
	@go vet ./...

generate:
	$(eval GENERATE_SRCS := $(shell git --no-pager grep -wl "go:generate" -- "*.go" ":(exclude)vendor"))
	@for GENERATE_SRC in $(GENERATE_SRCS); do \
		go generate $$GENERATE_SRC; \
	done

run: build
	@./build/$(TARGET)

test: test-all

test-all:
	@$(GOTEST) -v `go list ./... | grep -v test/e2e`

test-e2e:
	@$(GOTEST) -v `go list ./test/e2e` $(FLAGS)
