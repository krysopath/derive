TAG = $(shell git tag | sort -r --version-sort | head -n1)
SEMVERS = $(shell echo $(TAG) | semver)

major:
	@set -e;\
	export NEW=$$(echo $(TAG) | semver -release major| jq .canonical); \
	git tag -m "$$NEW major release" $$NEW
minor:
	@set -e;\
	export NEW=$$(echo $(TAG) | semver -release minor| jq .canonical); \
	echo git tag -m "$$NEW minor release" $$NEW
patch:
	@set -e;\
	export NEW=$$(echo $(TAG) | semver -release patch| jq .canonical); \
	echo git tag -m "$$NEW patch release" $$NEW

semver:
	git tag -f -m '$(TAG)' "$$(echo '$(SEMVERS)' | jq -r .major)"
	git tag -f -m '$(TAG)' "$$(echo '$(SEMVERS)' | jq -r .majorminor)"

release:
	git push
	git push --tags -f

gotests:
	go test ./...

gobuild:
	go build ./cmd/...; \

SHELL = /bin/bash
tests: gotests gobuild
	@shopt -s extglob; \
	echo -e "BATS are testing it now\n"; \
	for bat in ./cmd/**/*bats; do bats $$bat||{ head -n 6 $$bat; echo err:$$bat; exit 1; }; done; \

install: tests
	go install -trimpath -ldflags='-extldflags=-static -w -s -X main.version=$(TAG)' ./cmd/...
