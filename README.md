# RESTApi
Test RESTApi application

# Current version
1.02

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
- Demo-mode with full tables

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
1) GetContent (GET method) - endpoint to get all content from database.     
Usage: 
```
  host_address:port/content
```
Return:
```
  - HTTP 200 and all content data with params {id,protection_system_name,content_key,payload} if everything is good;
  - HTTP 502 and error="database problem", if you have some problem with database.
```
Example:
```
Usage:  
  curl -X GET "http://localhost:5000/content"
Return: 
    [{"id":1,"protection_system_name":"AES 2", 
      "content_key":"superpass", 
      "payload":"U2FsdGVkX190cOearjAhFozvAQFjW53OUhLQGKfTVZnj8iOwveiaZ8rqAPNBjeDB"},
    {"id":2,"protection_system_name":"AES 1", 
      "content_key":"mypassword", 
      "payload":"U2FsdGVkX1+lxfHPBsyNB+R1lJ2qOz/uA7NTprwWXhaMaQLNyhPRCyUq13VvkRDp"},
    {"id":3,"protection_system_name":"AES 1", 
      "content_key":"mypassword", 
      "payload":"U2FsdGVkX1+lxfHPBsyNB+R1lJ2qOz/uA7NTprwWXhaMaQLNyhPRCyUq13VvkRDp"}]
```
2) GetContentById (GET method) - endpoint to get content with selected ID (int).
Usage:
```
  host_address:port/content/<id>
```
Return:
```
  - HTTP 200 and selected content data with params {id,protection_system_name,content_key,payload} if everything is good;
  - HTTP 400 and error="invalid content ID", if you input some non-int <id> (content id MUST be INT);
  - HTTP 404 and error="no such content in database", if content with selected ID isn't in database;
  - HTTP 502 and error="database problem", if you have some problem with database.
```
Example:
```
  Usage:  
    curl -X GET "http://localhost:5000/content/1"
  Return: 
    {"id":1,"protection_system_name":"AES 2", 
      "content_key":"superpass", 
      "payload":"U2FsdGVkX190cOearjAhFozvAQFjW53OUhLQGKfTVZnj8iOwveiaZ8rqAPNBjeDB"}
____________________________________________________________________________________
  Usage with error:
    curl -i -X GET http://localhost:5000/content/xxx
  Result error:
    400 Bad Request
    invalid content ID
```
3) AddContent (POST method) - endpoint to add new content.    
Usage:
```
  host_address:port/content
```
JSON format:
```
  '{"protection_system_name":<string>,
    "content_key":<string>,
    "payload":<string>}'
```
where:
* protection_system_name - name of protection system in database (in demo ["AES 1","AES 2"]);
* content_key - encryption key to decrypt encrypted <payload>;
* payload - data, that was encrypted by symmetric <content_key>.
  Warning - every parameter MUST be in JSON structure.  

Return:
```
  - HTTP 200 if everything is good and new content was added;
  - HTTP 400 and error="invalid json", if you input bad JSON struct (see JSON struct above);
  - HTTP 400 and error="invalid content data", if you input not all of required content params (see WARNING above);
  - HTTP 400 and error="no such protection system in database", if selected protection system name isn't in database;
  - HTTP 400 and error="invalid payload", if you input bad data into JSON struct, that can't be decrypted (for example, if key isn't work with this encrypted data and aes-module get error);
  - HTTP 502 and error="database problem", if you have some problem with database.
```
Example:
```
  Usage:  
    curl -i -H "Content-Type: application/json" -X POST -d '{"protection_system_name":"AES 1", "content_key":"mypassword","payload":"U2FsdGVkX1+lxfHPBsyNB+R1lJ2qOz/uA7NTprwWXhaMaQLNyhPRCyUq13VvkRDp"}' http://localhost:5000/content
  Return: 
    200 OK
____________________________________________________________________________________  
  Usage with error:
    curl -i -H "Content-Type: application/json" -X POST -d '{"content_key":"mypassword","payload":"U2FsdGVkX1+lxfHPBsyNB+R1lJ2qOz/uA7NTprwWXhaMaQLNyhPRCyUq13VvkRDp"}' http://localhost:5000/content
  Result error:
    400 Bad Request
    invalid content ID
____________________________________________________________________________________
  Usage with error:
    curl -i -H "Content-Type: application/json" -X POST -d '{"protection_system_name":"AES BAD","content_key":"mypassword","payload":"U2FsdGVkX1+lxfHPBsyNB+R1lJ2qOz/uA7NTprwWXhaMaQLNyhPRCyUq13VvkRDp"}' http://localhost:5000/content
  Result error:
    400 Bad Request
    no such protection system in database
```
4) UpdateContent (PUT method) - endpoint to update content with selected ID (int).  
Usage:
```
  host_address:port/content/<id>
```
JSON format:
```
  '{"protection_system_name":<string>,
    "content_key":<string>,
    "payload":<string>}'
```
where:
* protection_system_name - name of protection system in database (in demo ["AES 1","AES 2"]);
* content_key - encryption key to decrypt encrypted <payload>;
* payload - data, that was encrypted by symmetric <content_key>
  Warning - all parameter are optional, but some of them needed to be included into JSON structure. For example: if you put only <payload> parameter into JSON - you will get error because encrypted data without <content_key> and <protection_system_name> can't be decrypted. Otherhand, you can input <protection_system_name> and <content_key> to update old payload. Server logic will decrypt old payload with old params from database and, if payload was not damaged, encrypt with new params.

Return:
```
  - HTTP 200 if everything is good and content was successfully updated;
  - HTTP 400 and error="invalid content ID", if you input some non-int <id> (content id MUST be INT);
  - HTTP 400 and error="invalid json", if you input bad JSON struct (see JSON struct above);
  - HTTP 400 and error="invalid content data", if you input not all of required content params (see WARNING above);
  - HTTP 400 and error="no such protection system in database", if selected protection system name isn't in database;
  - HTTP 404 and error="no such content in database", if content with selected ID isn't in database;
  - HTTP 417 and error="invalid payload in database", if somehow we put bad data in database before (for example, if old content_key and payload can't be usable);
  - HTTP 400 and error="invalid payload", if you input bad data into JSON struct, that can't be decrypted (for example, if key isn't work with this encrypted data and aes-module get error);
  - HTTP 502 and error="database problem", if you have some problem with database.
```
Example:
```
  Usage:  
    curl -i -H "Content-Type: application/json" -X PUT -d '{"content_key":"TestKey123"}' http://localhost:5000/content/3
    curl -i -H "Content-Type: application/json" -X PUT -d '{"content_key":"TestKey987","protection_system_name":"AES 2"}' http://localhost:5000/content/3
    curl -i -H "Content-Type: application/json" -X PUT -d '{"content_key":"popi","payload":"U2FsdGVkX18fO6a7VqCp2W2vcUGTbZqpzxJoHtR+80sy+ngb16+9OQBFPtH2aXxd"}' http://localhost:5000/content/3
    curl -i -H "Content-Type: application/json" -X PUT -d '{"protection_system_name":"AES 1","content_key":"mypassword","payload":"U2FsdGVkX1+lxfHPBsyNB+R1lJ2qOz/uA7NTprwWXhaMaQLNyhPRCyUq13VvkRDp"}' http://localhost:5000/content/3
  Return: 
    200 OK  
____________________________________________________________________________________
  Usage with error:
    curl -i -H "Content-Type: application/json" -X PUT -d '{"content_key":"TestKey23","payload":"U2FsdGVkX1+lxfHPBsyNB+R1lJ2qOz/uA7NTprwWXhaMaQLNyhPRCyUq13VvkRDp"}' http://localhost:5000/content/3
  Result error:
    400 Bad Request
    invalid payload
____________________________________________________________________________________
  Usage with error:
    curl -i -H "Content-Type: application/json" -X POST -d '{"protection_system_name":"AES BAD","content_key":"mypassword","payload":"U2FsdGVkX1+lxfHPBsyNB+R1lJ2qOz/uA7NTprwWXhaMaQLNyhPRCyUq13VvkRDp"}' http://localhost:5000/content
  Result error:
    400 Bad Request
    no such protection system in database
```
5) DeleteContent (DELETE method) - endpoint to delete content with selected ID (int).  
Usage:
```
  host_address:port/content/<id>
```
Return:
```
  - HTTP 200 if everything is good and content was successfully deleted;
  - HTTP 400 and error="invalid content ID", if you input some non-int <id> (content id MUST be INT);
  - HTTP 502 and error="database problem", if you have some problem with database.
```
Example:
```
  Usage:  
    curl -X DELETE "http://localhost:5000/content/1"
  Result:
    200 OK
```
6) CheckView (GET method) - endpoint to get decrypted payload data from content with selected ID (int) with selected device name (string). 
Usage: 
```
  host_address:port/view
```
JSON format:
```
  '{"content_id":<int>,
    "device":<string>}'
```
where:
* content_id - id of selected content;
* device - name of device (in demo ["Android","Samsung","iOS","LG"]).
  Warning - every parameter MUST be in JSON structure.

Return:
```
  - HTTP 200 and decrypted data of selected content if everything is good and user can get this data;
  - HTTP 451 and error="content can't be shown", if selected device can't work with selected encrypted data(for example - if content was encrypted by AES-CBC and device can work with AES-ECB)
  - HTTP 400 and error="invalid json", if you input bad JSON struct (see JSON struct above);
  - HTTP 400 and error="invalid view content data", if you input not all of required view content params (see WARNING above);
  - HTTP 400 and error="no such device in database", if selected device name isn't in database;
  - HTTP 400 and error="no such protection system in encryption module", if protection system of selected device can't be decrypet with current encryption module  (for example, if we somehow add some new device with protection scheme,that absense in current version of encryption module);
  - HTTP 417 and error="invalid payload in database", if somehow we put bad data in database before (for example, if somehow someone put content_key and payload in database, that can't be usable);
  - HTTP 502 and error="database problem", if you have some problem with database.
```
Example:
```
  Usage:  
    curl -i -H "Content-Type: application/json" -X GET -d '{"content_id":3,"Device":"iOS"}' http://localhost:5000/view
  Return: 
    200 OK
    "precious-content"
____________________________________________________________________________________
  Usage with error:
    curl -i -H "Content-Type: application/json" -X GET -d '{"content_id":3,"Device":"OS"}' http://localhost:5000/view
  Return: 
    400 Bad Request
    no such device in database
____________________________________________________________________________________
  Usage with error:
     curl -i -H "Content-Type: application/json" -X GET -d '{"content_id":3,"Device":"LG"}' http://localhost:5000/view
  Return: 
    451 Unavailable For Legal Reasons
    content can't be shown
```
# Requests examples
1) curl -X GET "http://localhost:5000/content"
2) curl -X GET "http://localhost:5000/content/12"
3) curl -i -H "Content-Type: application/json" -X POST -d '{"protection_system_name":"testAes","content_key":"testKey","payload":"testPayload"}' http://localhost:5000/content
4) curl -i -H "Content-Type: application/json" -X PUT -d '{"protection_system_name":"testAes","content_key":"testKey","payload":"testPayload"}' http://localhost:5000/content/1
5) curl -i -H "Content-Type: application/json" -X DELETE http://localhost:5000/content/1
6) curl -i -H "Content-Type: application/json" -X POST -d '{"content_id":1,"Device":"LG"}' http://localhost:5000/content/view

# Tests
To run unit-tests you need to go to folder where you can find some 'XXX_test.go' file and run:
```
go test
```
