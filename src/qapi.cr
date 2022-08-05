require "http/client"
require "json"
require "log"
require "option_parser"
require "uri"
require "yaml"

require "../src/model/config"
require "../src/execution/runner"

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

runner = Execution::Runner.from_yaml(yaml)
runner.start()