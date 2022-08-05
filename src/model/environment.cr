require "json"

module Model
  class Environment
    include YAML::Serializable
  
    @[YAML::Field(key: "base_url")]
    property base_url : String
  end
end