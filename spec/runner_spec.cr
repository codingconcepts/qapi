require "webmock"

require "./spec_helper"

describe Execution::Runner do
	it "makes a request with a body and extracts response body" do
		WebMock.stub(:post, "http://localhost:8080/test/endpoint").
			with(headers: {"Content-Type" => "application/json"}).
			to_return(status: 200, body: %({"result": "some_result"}))

		yaml = File.open("spec/test_config_body.yaml")
		runner = Execution::Runner.from_yaml(yaml)
		runner.start()

		runner.variables["result"].should eq "some_result"
	end

	it "makes a request with header variables" do
		WebMock.stub(:get, "http://localhost:8080/test/endpoint").
			with(headers: {"Authorization" => "Bearer some_token"}).
			to_return(status: 200)

		yaml = File.open("spec/test_config_headers.yaml")
		runner = Execution::Runner.from_yaml(yaml)	
		runner.start()
	end

	it "makes a request with url params" do
		WebMock.stub(:get, "http://localhost:8080/test/endpoint/some_id").
			to_return(status: 200)

		yaml = File.open("spec/test_config_params.yaml")
		runner = Execution::Runner.from_yaml(yaml)	
		runner.start()
	end
end