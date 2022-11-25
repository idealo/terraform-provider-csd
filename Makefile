TEST?=$$(go list ./...)
HOSTNAME=idealo.com
NAMESPACE=transport
NAME=idealo-tools
BINARY=terraform-provider-${NAME}
VERSION=0.0.1
OS_ARCH=$(shell go version | awk '{gsub("/","_");print $4}')

.PHONY: clean

default: help

clean:
	rm -fr ${BINARY}
	rm -fr csd/csd

build: clean
	go build -o ${BINARY}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

format:
	go fmt ./...

install: build
	install -Dm0755 ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

uninstall:
	rm -fr ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}

fake_api: clean
	cd csd; go build -o csd
	./csd/csd
