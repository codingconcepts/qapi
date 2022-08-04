login_example_server:
	crystal examples/login/server.cr

login_example:
	crystal src/qapi.cr -c examples/login/config.yaml
