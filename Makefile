.PHONY: build

validate_version:
ifndef VERSION
	$(error VERSION is undefined)
endif

login_example_server:
	crystal examples/login/server.cr

login_example:
	crystal src/qapi.cr -c examples/login/config.yaml

build: validate_version
	crystal build src/qapi.cr --release -o ./build/qapi
	(cd build && tar -zcvf qapi_${VERSION}_macos.tar.gz ./qapi)

test:
	crystal spec -p