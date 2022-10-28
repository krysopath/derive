TAG = $(shell git tag | sort -r --version-sort | head -n1)
SEMVERS = $(shell echo $(TAG) | semver)

deps:
	@which semver || { go install github.com/krysopath/semver/cmd/semver@v1; which semver; }
	@which bats || { sudo apt install bats; which bats; }
	@which jq || { sudo apt install jq; which jq; }
major:
	@set -e;\
	export NEW=$$(echo $(TAG) | semver -release major| jq -r .canonical | tee /dev/stderr); \
	git tag -m "$$NEW major release" $$NEW
minor:
	@set -e;\
	export NEW=$$(echo $(TAG) | semver -release minor| jq -r .canonical | tee /dev/stderr); \
	git tag -m "$$NEW minor release" $$NEW 
patch:
	@set -e;\
	export NEW=$$(echo $(TAG) | semver -release patch| jq -r .canonical | tee /dev/stderr); \
	git tag -m "$$NEW patch release" --sign $$NEW
semver:
	@git tag -f -m '$(TAG)' "$$(echo '$(SEMVERS)' | jq -r .major | tee /dev/stderr)"
	@git tag -f -m '$(TAG)' "$$(echo '$(SEMVERS)' | jq -r .majorminor | tee /dev/stderr )"
release: semver
	git push; git push --tags -f
update:
	go list -m -u all \
	| awk -F" " '{ if ($$3 != "") print $$1 " " $$3; }' \
	| xargs -l bash -c 'VERSION=$$(grep -Po "(?<=\[).+(?=\])" <<<$$1); go get $$0@$$VERSION'
	go mod tidy
gotests:
	go test ./...
gobuild:
	go build ./cmd/...
tests: gotests gobuild
	@printf "BATS are testing it now\n"; \
	for bat in $$(find . -name '*.bats'); do \
		bats $$bat || { head -n 6 $$bat; echo err:$$bat; exit 1; }; \
	done;
install: tests
	go install -trimpath -ldflags='-extldflags=-static -w -s -X main.version=$(TAG)' ./cmd/...
	@bash -c "[ $$(git diff $(TAG) -- **/*.bats | wc -c) -eq 0 ]" || echo "!!! breaking changes"
