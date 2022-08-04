# ![datagen logo](assets/cover.png)

# qapi
Run Quality Assurance tests against an API endpoint

## Installation

TODO

## Examples

#### Login

The login example is a simple test that targets a server with the following two endpoints:

* **/api/login** - A POST request endpoint that takes a username and password and returns a user id and token.
* **/api/get/:id** - A GET request endpoint that takes a user id as a URL param.

To execute the login example, run the following commands in different terminal windows:

```
$ make login_example_server
crystal examples/login/server.cr

[development] Kemal is ready to lead at http://0.0.0.0:8080
"un"
"pw"
2022-08-04 18:08:07 UTC 200 POST /api/login 162.82µs
HTTP::Headers{"Authorisation" => "Bearer MXU01u5KSMQ0SCNL4/6AFuP+DhZ7AoXWTIfmd7gl6Sp6vJQn0C2w6A/NsqZoBeGnZpw", "Connection" => "close", "Content-Length" => "0", "Host" => "localhost:8080", "User-Agent" => "Crystal", "Accept-Encoding" => "gzip, deflate"}
"4977feb8-fac2-4c2a-b608-771ae8b0f081"
2022-08-04 18:08:07 UTC 200 GET /api/get/4977feb8-fac2-4c2a-b608-771ae8b0f081 572.22µs
```

```
$ make login_example
crystal src/qapi.cr -c examples/login/config.yaml

2022-08-04T18:08:07.889369Z   INFO - request -- name: "login"
2022-08-04T18:08:07.893480Z   INFO - request -- name: "get"
```