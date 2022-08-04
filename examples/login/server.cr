require "kemal"

post "/api/login" do |env|
    username = env.params.json["username"].as(String)
    password = env.params.json["password"].as(String)

    p username, password
    
    {
        "result": {
            "token": "MXU01u5KSMQ0SCNL4/6AFuP+DhZ7AoXWTIfmd7gl6Sp6vJQn0C2w6A/NsqZoBeGnZpw",
            "id": "4977feb8-fac2-4c2a-b608-771ae8b0f081"
        }
    }.to_json
end  

get "/api/get/:token" do |env|
    token = env.params.url["token"]

    p env.request.headers
    p token
end

Kemal.run 8080