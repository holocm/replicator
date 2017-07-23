PKG    = github.com/holocm/replicator
PREFIX := /usr

all: build/replicator

GO            := GOPATH=$(CURDIR)/.gopath GOBIN=$(CURDIR)/build go
GO_BUILDFLAGS :=
GO_LDFLAGS    := -s -w

# This target uses the incremental rebuild capabilities of the Go compiler to speed things up.
# If no source files have changed, `go install` exits quickly without doing anything.
build/replicator: FORCE
	$(GO) install $(GO_BUILDFLAGS) -ldflags '$(GO_LDFLAGS)' '$(PKG)'

install: FORCE all
	install -D -m 0755 build/replicator "$(DESTDIR)$(PREFIX)/bin/replicator"

vendor: FORCE
	@# vendoring by https://github.com/holocm/golangvend
	golangvend

.PHONY: FORCE
