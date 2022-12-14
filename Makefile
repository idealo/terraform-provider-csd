TEST?=$$(go list ./...)
HOSTNAME=idealo.com
NAMESPACE=transport
NAME=csd
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=$(shell go version | awk '{gsub("/","_");print $$4}')

.PHONY: clean

default: help

clean:
	rm -fr ${BINARY}

build: clean
	go build -o ${BINARY}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

format:
	go fmt ./...

install: build
	#install -Dm0755 ${BINARY} ~/.terraform.d/plugins/idealo.com/transport/csd/0.0.1/linux_amd64/terraform-provider-csd_v0.0.1
	install -Dm0755 ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}/${BINARY}_v${VERSION}

uninstall:
	rm -fr ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}

tf: install
	rm -fr examples/.terraform.lock.hcl examples/.terraform
	cd examples; terraform init
	cd examples; env AWS_ACCESS_KEY_ID="" AWS_SECRET_ACCESS_KEY=aaa AWS_SESSION_TOKEN=aaa terraform apply --auto-approve

fake_api:
	rm -fr csd/csd
	cd csd; go build -o csd
	./csd/csd
