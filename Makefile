HOSTNAME=idealo.com
NAMESPACE=transport
NAME=csd
BINARY=terraform-provider-${NAME}
VERSION=2.0.0
OS_ARCH=$(shell go version | awk '{gsub("/","_");print $$4}')
#AWS_PROFILE=sandbox

.PHONY: clean build format install uninstall test facade help

default: help

clean:  ## Cleanup sources
	rm -fr ${BINARY}

build: clean  ## Run go build
	go build -o ${BINARY}

format:  ## Run go fmt
	go fmt ./...

install: build  ## Install provider
	install -Dm0755 ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}_v${VERSION}

uninstall:  ## Remove provider
	rm -fr ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}

update:  ## Update dependencies
	go get -t -u ./...
	go mod tidy

test: install  ## Run example terraform
	rm -fr examples/.terraform.lock.hcl examples/.terraform
	cd examples; terraform init -upgrade
	cd examples; terraform apply --auto-approve

debug:
	dlv exec --accept-multiclient --continue --headless ./${BINARY} -- -debug

docs:  ## Generate documentation
	go generate ./...

help:  ## Display this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
