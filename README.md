# RESTApi
Test RESTApi application

# Current version
1.01

# REQUIREMENT
- Go 1.9.2
- openssl
- MySql server
- Docker (MySql tests)

# Features
- REST Api
- MySql module
- Encryption module, based on 'openssl'
- Logger

# Pre-build
To install all needed external libs:
```
go get ./...
```
# Build
To build executable file:
```
go builg restapiserver.go
```
# Usage 
To launch application in demo-mode:
```
./restapiserver -demo
```
To launch application in normal mode:
```
./restapiserver -port=<SERVER_PORT> -db-host=<HOST> -db-port=<PORT> -db-name=<DATABASE_NAME> -db-user=<USERNAME> -db-pass=<PASSWORD>
```
Help output:
```
  -db-host string
        MySql host address (default "localhost")
  -db-name string
        MySql database name (default "rakuten")
  -db-pass string
        MySql password (default "root")
  -db-port int
        MySql port index (default 8889)
  -db-user string
        MySql user  (default "root")
  -demo
        flag to launch RESTApi server in demo mode
  -port int
        port for RESTApi server (default 5000)
  -version
        get current version.
```
# Api usage
1) GetContent      - GET     host_address:port/content
2) GetContentById  - GET     host_address:port/content/<id>
3) AddContent      - POST    host_address:port/content
4) UpdateContent   - PUT     host_address:port/content/<id>
5) DeleteContent   - DELETE  host_address:port/content/<id>
6) CheckView       - POST    host_address:port/content/view

# Requests examples
1) curl -X GET "http://localhost:5000/content"
2) curl -X GET "http://localhost:5000/content/12"
3) curl -i -H "Content-Type: application/json" -X POST -d '{"protection_system_name":"testAes","content_key":"testKey","payload":"testPayload"}' http://localhost:5000/content
4) curl -i -H "Content-Type: application/json" -X PUT -d '{"protection_system_name":"testAes","content_key":"testKey","payload":"testPayload"}' http://localhost:5000/content/1
5) curl -i -H "Content-Type: application/json" -X DELETE http://localhost:5000/content/1
6) curl -i -H "Content-Type: application/json" -X POST -d '{"content_id":1,"Device":"LG"}' http://localhost:5000/content/view

# Tests
Yes
