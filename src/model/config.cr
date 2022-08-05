require "yaml"

module Model
  class Config
    include YAML::Serializable
  
    @[YAML::Field(key: "environment")]
    property environment : Environment
  
    @[YAML::Field(key: "variables")]
    property variables : Hash(String, String)
  
    @[YAML::Field(key: "requests")]
    property requests : Array(Request)
  end
end