TAG = $(shell git describe --tags --abbrev=0)
SEMVERS = $(shell echo $(TAG) | semver)

semver:
	git tag -m '$(TAG)' --sign $$(echo '$(SEMVERS)' | jq -r .major)
	git tag -m '$(TAG)' --sign $$(echo '$(SEMVERS)' | jq -r .majorminor)
	git push --tags

install:
	go install ./cmd/...
