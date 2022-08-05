require "json"

module Model
    class Extractor
        include YAML::Serializable
      
        @[YAML::Field(key: "type")]
        property type : String
      
        @[YAML::Field(key: "selectors")]
        property selectors : Hash(String, String)
      end
end