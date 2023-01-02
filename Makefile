HOSTNAME=idealo.com
NAMESPACE=transport
NAME=csd
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=$(shell go version | awk '{gsub("/","_");print $$4}')

.PHONY: clean build format install uninstall tf fake_api help

default: help

clean:  ## Cleanup sources
	rm -fr ${BINARY}

build: clean  ## Run go build
	go build -o ${BINARY}

format:  ## Run go fmt
	go fmt ./...

install: build  ## Install provider
	#install -Dm0755 ${BINARY} ~/.terraform.d/plugins/idealo.com/transport/csd/0.0.1/linux_amd64/terraform-provider-csd_v0.0.1
	install -Dm0755 ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}_v${VERSION}

uninstall:  ## Remove provider
	rm -fr ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}

tf: install  ## Run example terraform
	rm -fr examples/.terraform.lock.hcl examples/.terraform
	cd examples; terraform init
	#cd examples; env AWS_ACCESS_KEY_ID=aaa AWS_SECRET_ACCESS_KEY=aaa AWS_SESSION_TOKEN=aaa terraform apply --auto-approve
	cd examples; terraform apply --auto-approve

fake_api:  ## Build fake API
	rm -fr csd/csd
	cd csd; go build -o csd
	./csd/csd

help:  ## Display this help
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
