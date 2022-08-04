# ![datagen logo](assets/cover.png)

# qapi
Run Quality Assurance tests against an API endpoint

## Installation

TODO

## Configuration

To execute aqpi against an endpoint, run the following command:

```
$ qapi -c PATH_TO_A_CONFIG.yaml
```

The config file provides qapi with a set of instructions. These instructions are simply a list of requests to make and some initial configuration to use when making those requests. Here's an annotated example config file.

``` yaml
# Configuration variables that are known to qapi and provide special meaning.
environment:
	# The URL from which all other URLs will be built.
  base_url: http://localhost:8080/api

# A collection of runtime variables that will be used to further configure
# qapi. Any variables obtained from extractors will be added to this collection.
variables:
  username: un
  password: pw

# A collection of requests to perform agains the API.
requests:
  # A request to perform. The header, path, and body parts can be parameterised
  # with variables.
  - name: login
    headers:
      Content-Type: application/json
    method: POST
    path: "/login"
    body: |
      {
        "username": "{{username}}",
        "password": "{{password}}"
      }
    # A collection of extractors, used to extract values from API responses.
    extractors:
      # A JSON extractor, that selects id, and token values from a top-level
      # "result" object and applies them to the runtime variables collection.
      - type: json
        selectors:
          id: result.id
          token: result.token
  
  - name: get
    method: GET
    path: "/get/{{id}}"
    headers:
      Authorisation: Bearer {{token}}
```

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