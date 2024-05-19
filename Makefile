SHELL := /bin/bash

PROJECT_NAME := "go-admin"
PKG := "$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v /api/)



.PHONY: ci-lint
# check the code specification against the rules in the .golangci.yml file
ci-lint:
	@gofmt -s -w .
	golangci-lint run ./...


.PHONY: test
# go test *_test.go files, the parameter -count=1 means that caching is disabled
test:
	go test -count=1 -short ${PKG_LIST}


.PHONY: cover
# generate test coverage
cover:
	go test -short -coverprofile=cover.out -covermode=atomic ${PKG_LIST}
	go tool cover -html=cover.out


.PHONY: graph
# generate interactive visual function dependency graphs
graph:
	@echo "generating graph ......"
	@cp -f cmd/admin/main.go .
	go-callvis -skipbrowser -format=svg -nostd -file=admin go-admin
	@rm -f main.go admin.gv


.PHONY: docs
# generate swagger docs, , only for ⓵ Web services created based on
docs:
	go mod tidy
	@gofmt -s -w .
	@bash scripts/swag-docs.sh $(HOST)


.PHONY: build
# build admin for linux amd64 binary
build:
	@echo "building 'admin', linux binary file will output to 'cmd/admin'"
	@cd cmd/admin && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build


.PHONY: run
# build and run service
run:
	@bash scripts/run.sh


.PHONY: run-nohup
# run service with nohup in local, if you want to stop the server, pass the parameter stop, e.g. make run-nohup CMD=stop
run-nohup:
	@bash scripts/run-nohup.sh $(CMD)


.PHONY: run-docker
# deploy service in local docker, if you want to update the service, run the make run-docker command again.
run-docker: image-build-local
	@bash scripts/deploy-docker.sh


.PHONY: binary-package
# packaged binary files
binary-package: build
	@bash scripts/binary-package.sh


.PHONY: deploy-binary
# deploy binary to remote linux server, e.g. make deploy-binary USER=root PWD=123456 IP=192.168.1.10
deploy-binary: binary-package
	@expect scripts/deploy-binary.sh $(USER) $(PWD) $(IP)


.PHONY: image-build-local
# build image for local docker, tag=latest, use binary files to build
image-build-local: build
	@bash scripts/image-build-local.sh


.PHONY: image-build
# build image for remote repositories, use binary files to build, e.g. make image-build REPO_HOST=addr TAG=latest
image-build:
	@bash scripts/image-build.sh $(REPO_HOST) $(TAG)


.PHONY: image-build2
# build image for remote repositories, phase II build, e.g. make image-build2 REPO_HOST=addr TAG=latest
image-build2:
	@bash scripts/image-build2.sh $(REPO_HOST) $(TAG)


.PHONY: image-push
# push docker image to remote repositories, e.g. make image-push REPO_HOST=addr TAG=latest
image-push:
	@bash scripts/image-push.sh $(REPO_HOST) $(TAG)


.PHONY: deploy-k8s
# deploy service to k8s
deploy-k8s:
	@bash scripts/deploy-k8s.sh


.PHONY: update-config
# update internal/config code base on yaml file
update-config:
	@sponge config --server-dir=.


.PHONY: clean
# clean binary file, cover.out, template file
clean:
	@rm -vrf cmd/admin/admin*
	@rm -vrf cover.out
	@rm -vrf main.go admin.gv
	@rm -vrf internal/ecode/*.go.gen*
	@rm -vrf internal/routers/*.go.gen*
	@rm -vrf internal/handler/*.go.gen*
	@rm -vrf internal/service/*.go.gen*
	@rm -rf admin-binary.tar.gz
	@echo "clean finished"


# show help
help:
	@echo ''
	@echo 'Usage:'
	@echo '  make <target>'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-_0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf "\033[1;36m  %-22s\033[0m %s\n", helpCommand,helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := all
