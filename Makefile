.PHONY: build

validate_version:
ifndef VERSION
	$(error VERSION is undefined)
endif

login_example_server:
	go run examples/login/server.go

login_example:
	go run qapi.go -config examples/login/config.yaml

release: validate_version
	GOOS=linux go build -ldflags "-X main.version=${VERSION}" -o build/qapi
	(cd build && tar -zcvf qapi_${VERSION}_linux.tar.gz ./qapi)

	GOOS=darwin go build -ldflags "-X main.version=${VERSION}" -o build/qapi
	(cd build && tar -zcvf qapi_${VERSION}_macos.tar.gz ./qapi)

	GOOS=windows go build -ldflags "-X main.version=${VERSION}" -o build/qapi
	(cd build && tar -zcvf qapi_${VERSION}_windows.tar.gz ./qapi)

test:
	go test ./... -v -cover

coverage:
	go test ./... --coverprofile=coverage.out
	go tool cover --html=coverage.out