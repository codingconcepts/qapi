require "json"

module Model
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
end