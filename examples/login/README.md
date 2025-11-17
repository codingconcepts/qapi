Start server

```sh
go run examples/login/server.go
```

Start qapi

```sh
go run qapi.go \
--config examples/login/config.yaml \
--vus 1 \
--duration 10s
```