VERSION=0.0.1
NAME=idealo-tools
BINARY=terraform-provider-${NAME}

.PHONY: clean

clean:
	rm -fr csd/csd

build:
	go build -o ${BINARY}

fake_api:
	cd csd; go build -o csd
	./csd/csd
