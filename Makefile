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
	GOOS=linux go build -ldflags "-X main.version=${VERSION}" -o build/aqpi
	(cd build && tar -zcvf aqpi_${VERSION}_linux.tar.gz ./aqpi)

	GOOS=darwin go build -ldflags "-X main.version=${VERSION}" -o build/aqpi
	(cd build && tar -zcvf aqpi_${VERSION}_macos.tar.gz ./aqpi)

	GOOS=windows go build -ldflags "-X main.version=${VERSION}" -o build/aqpi
	(cd build && tar -zcvf aqpi_${VERSION}_windows.tar.gz ./aqpi)

test:
	go test ./... -v -cover