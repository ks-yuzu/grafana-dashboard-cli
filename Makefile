default: run

repo := $(shell basename $(shell pwd))

# == binary ========================================
BIN := bin/$(repo)
SRC := $(shell find . -type f -name '*.go')

LDFLAGS    := -s -w
BULID_OPTS := -ldflags="$(LDFLAGS)" -trimpath

init:
	go mod download

$(BIN): build
build: $(SRC) version.go
	CGO_ENABLED=0 go build $(BULID_OPTS) -o $(BIN) cmd/main.go
# go test -c ./cmd/main_test.go -o $(BIN)

run: $(BIN)
	$(BIN)

watch-and-build: $(SRC) version.go
	inotifywait -e CLOSE_WRITE -m ./cmd/main.go -m ./pkg/* -m ./pkg/*/* --format '%w%f' | while read file; do \
	echo "detect update of ${file}"; \
	make build && echo done.; \
  echo; \
	done

clean:
	$(RM) bin/*


# == image =========================================
IMAGE_NAME := kustomize-snapshot-test
VERSION    := $(shell grep VERSION version.go | cut -d'"' -f2)
IMAGE_TAG  := $(IMAGE_NAME):$(VERSION)

image:
	docker build . -t $(IMAGE_TAG)

push:
	docker push $(IMAGE_TAG)
