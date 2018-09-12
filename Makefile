version := $(shell git tag --points-at $(git rev-parse HEAD) | grep "v\(.*\)" 2> /dev/null)
release_args := --user Wizcorp --repo terraform-provider-ncloud --tag $(version)

# Print out the list of known regions, zones, etc/
list-services:
	@go run \
		src/ncloud-products-list/main.go
.PHONY: generate-services

# Generate Services.md
generate-services:
	@go run \
		src/ncloud-products-list/main.go > Services.md
.PHONY: generate-services

# Build the provider 
# CGO_ENABLED=0 must be set to run on Alpine
# See https://stackoverflow.com/a/36308464/262831
build:
	CGO_ENABLED=0 go build \
		-o build/terraform-provider-ncloud \
		src/terraform-provider-ncloud/*.go
.PHONY: build

# Make a release for all supported platforms
release-all:
ifndef version 
	$(error usage: current commit is not tagged, please make sure to tag before releasing)
endif
	git push --tags upstream master
	github-release release $(release_args) \
		--name $(version) \
		--description $(version)

	mmake release target=linux
	mmake release target=darwin
	mmake release target=windows
.PHONY: release-all

zipfile := terraform-provider-ncloud-$(version)-$(target).zip

# Make a release for a specific target platform
release:
ifndef target
	$(error usage: mmake release target=(linux|darwin|windows))
endif
	echo $(version)
ifndef version 
	$(error usage: current commit is not tagged, please make sure to tag before releasing)
endif
	GOOS=$(target) mmake build
ifeq ($(OS),Windows_NT)
	Compress-Archive -Path ./build/terraform-provider-ncloud -CompressionLevel Fastest -DestinationPath $(zipfile)
else
	cd ./build && zip ../$(zipfile) ./terraform-provider-ncloud
endif

	github-release upload $(release_args) \
		--name $(zipfile) \
		--file $(zipfile)
	rm $(zipfile)
.PHONY: release