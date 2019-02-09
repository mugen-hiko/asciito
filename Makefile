srcs := $(shell \
  find . \
    -type d -path './.git' -prune -o \
    -type d -path './bin' -prune -o \
    -type d -path './vendor' -prune -o \
    -type f -name '*.go' -print \
)

.PHONY: init
init: bin/cobra

.PHONY: build
build: bin/asciito

bin/cobra:
	cd vendor/github.com/spf13/cobra/cobra \
	  && go build -o $(PWD)/$@ ./

bin/asciito: $(srcs)
	go build -o $@ ./
