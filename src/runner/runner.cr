require "../model/*"

module Runner
    def self.make_requests(c : Model::Config)
        c.requests.each do |r|
          Log.info &.emit("request", name: r.name)
      
              elapsed = Time.measure do
                  make_request(c, r)
              end
      
              Log.debug &.emit("request", elapsed_time: "#{elapsed.total_milliseconds}ms")
        end
    end

    def self.make_request(c : Model::Config, r : Model::Request)
        url = add_variables(c, r.path)
        body = add_variables(c, r.body || "")
        headers = attach_headers(c, r)
      
          uri = URI.parse Path.new(c.environment.base_url, url).to_s
        resp = HTTP::Client.exec(r.method, uri, headers, body)
        raise "#{resp.body}" if resp.status_code != 200
      
        run_extractors(c, r, resp)
      end
      
      def self.add_variables(c : Model::Config, s : String) : String
        matches = s.scan(/{{\w+}}/).map(&.[0])
        matches.each do |m|
          s = s.gsub(m, c.variables[m.strip("{}")])
        end
      
        s
      end
      
      def self.attach_headers(c : Model::Config, r : Model::Request) : HTTP::Headers
        headers = HTTP::Headers.new
        return headers if !r.headers
      
        h = r.headers.not_nil!
        h.each do |k, v|
          h[k] = add_variables(c, v)
        end
      
        headers.merge! h
      
        headers
      end
      
      def self.run_extractors(c : Model::Config, r : Model::Request, resp : HTTP::Client::Response)
        return if !r.extractors
      
        r.extractors.not_nil!.each do |e|
          case e.type
          when "json"
            run_json_extractor(c, e.selectors, resp)
          end
        end
      end
      
      def self.run_json_extractor(c : Model::Config, selectors : Hash(String, String), resp : HTTP::Client::Response)
        body = JSON.parse(resp.body)
      
        selectors.each do |sk, dot_path|
          b = body
          dot_path.split(".").each do |part|
            b = b[part]
          end
      
          c.variables[sk] = b.to_s
        end
      end      
end