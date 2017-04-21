PROJECT_NAME="siren"
TARGET_OSARCH="linux/amd64"

build: deps
	@which gox > /dev/null || go get github.com/mitchellh/gox
	@gox -output=./bin/$(PROJECT_NAME) -osarch="$(TARGET_OSARCH)" .

update:
	@godep save ./...

deps:
	@which godep > /dev/null || go get github.com/tools/godep
	@godep go install ./...

format:
	@which goimports > /dev/null || go get golang.org/x/tools/cmd/goimports
	@echo "--> Running goimports"
	@goimports -w .

clean:
	@rm -rf ./bin/* > /dev/null

