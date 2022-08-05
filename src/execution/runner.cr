require "yaml"

require "../model/*"

module Execution
  class Runner
    include YAML::Serializable
  
    @[YAML::Field(key: "environment")]
    property environment : Model::Environment
  
    @[YAML::Field(key: "variables")]
    property variables : Hash(String, String)
  
    @[YAML::Field(key: "requests")]
    property requests : Array(Model::Request)

    def start()
      requests.each do |r|
        Log.info &.emit("request", name: r.name)
    
        elapsed = Time.measure do
          make_request(r)
        end
  
        Log.debug &.emit("request", elapsed_time: "#{elapsed.total_milliseconds}ms")
      end
    end

    def make_request(r : Model::Request)
      url = add_variables(r.path)
      body = add_variables(r.body || "")
      headers = attach_headers(r)
    
      uri = URI.parse Path.new(environment.base_url, url).to_s
      resp = HTTP::Client.exec(r.method, uri, headers, body)
      raise "#{resp.body}" if resp.status_code != 200
    
      run_extractors(r, resp)
    end
    
    def add_variables(s : String) : String
      matches = s.scan(/{{\w+}}/).map(&.[0])
      matches.each do |m|
        s = s.gsub(m, variables[m.strip("{}")])
      end
    
      s
    end
    
    def attach_headers(r : Model::Request) : HTTP::Headers
      headers = HTTP::Headers.new
      return headers if !r.headers
    
      h = r.headers.not_nil!
      h.each do |k, v|
        h[k] = add_variables(v)
      end
    
      headers.merge! h
    
      headers
    end
    
    def run_extractors(r : Model::Request, resp : HTTP::Client::Response)
      return if !r.extractors
    
      r.extractors.not_nil!.each do |e|
        case e.type
        when "json"
          run_json_extractor(e.selectors, resp)
        end
      end
    end
    
    def run_json_extractor(selectors : Hash(String, String), resp : HTTP::Client::Response)
      body = JSON.parse(resp.body)
    
      selectors.each do |sk, dot_path|
        b = body
        dot_path.split(".").each do |part|
          b = b[part]
        end
    
        variables[sk] = b.to_s
      end
    end
  end
end