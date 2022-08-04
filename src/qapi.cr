require "http/client"
require "json"
require "log"
require "option_parser"
require "uri"
require "yaml"

file_path = ""

OptionParser.parse do |parser|
  parser.on "-h", "--help", "display help information" do
    puts parser
    exit
  end
	
	parser.on "-c", "--config=<path>", "path to a config file (required)" { |x| file_path = x }
	parser.on "-d", "--debug", "show debug output" { Log.setup(:debug) }
end

if file_path == ""
  puts "missing -c, --config agument"
  exit
end

yaml = File.open(file_path)
config = Config.from_yaml(yaml)

make_requests(config)

def make_requests(c : Config)
  c.requests.each do |r|
    Log.info &.emit("request", name: r.name)

		elapsed = Time.measure do
			make_request(c, r)
		end

		Log.debug &.emit("request", elapsed_time: "#{elapsed.total_milliseconds}ms")
  end
end

def make_request(c : Config, r : Request)
  url = add_variables(c, r.path)
  body = add_variables(c, r.body || "")
  headers = attach_headers(c, r)

	uri = URI.parse Path.new(c.environment.base_url, url).to_s
  resp = HTTP::Client.exec(r.method, uri, headers, body)
  raise "#{resp.body}" if resp.status_code != 200

  run_extractors(c, r, resp)
end

def add_variables(c : Config, s : String) : String
  matches = s.scan(/{{\w+}}/).map(&.[0])
  matches.each do |m|
    s = s.gsub(m, c.variables[m.strip("{}")])
  end

  s
end

def attach_headers(c : Config, r : Request) : HTTP::Headers
  headers = HTTP::Headers.new
  return headers if !r.headers

  h = r.headers.not_nil!
  h.each do |k, v|
    h[k] = add_variables(c, v)
  end

  headers.merge! h

  headers
end

def run_extractors(c : Config, r : Request, resp : HTTP::Client::Response)
  return if !r.extractors

  r.extractors.not_nil!.each do |e|
    case e.type
    when "json"
      run_json_extractor(c, e.selectors, resp)
    end
  end
end

def run_json_extractor(c : Config, selectors : Hash(String, String), resp : HTTP::Client::Response)
  body = JSON.parse(resp.body)

  selectors.each do |sk, dot_path|
    b = body
    dot_path.split(".").each do |part|
      b = b[part]
    end

    c.variables[sk] = b.to_s
  end
end

class Config
  include YAML::Serializable

  @[YAML::Field(key: "environment")]
  property environment : Environment

  @[YAML::Field(key: "variables")]
  property variables : Hash(String, String)

  @[YAML::Field(key: "requests")]
  property requests : Array(Request)
end

class Environment
  include YAML::Serializable

  @[YAML::Field(key: "base_url")]
  property base_url : String
end

class Request
  include YAML::Serializable

  @[YAML::Field(key: "name")]
  property name : String

  @[YAML::Field(key: "headers")]
  property headers : Hash(String, String)?

  @[YAML::Field(key: "path")]
  property path : String

  @[YAML::Field(key: "method")]
  property method : String

  @[YAML::Field(key: "body")]
  property body : String?

  @[YAML::Field(key: "extractors")]
  property extractors : Array(Extractor)?
end

class Extractor
  include YAML::Serializable

  @[YAML::Field(key: "type")]
  property type : String

  @[YAML::Field(key: "selectors")]
  property selectors : Hash(String, String)
end
