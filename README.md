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
# MySQL tables
1) Protection systems

| id  | name  | encryption_mode |
| --- | ----- |---------------- |
|  1  | AES 1 | AES + ECB       |
|  2  | AES 2 | AES + CBC       |

Was taken from exercise.

2) Devices

| id  | name  | protection_system_id |
| --- | ----- |--------------------- |
|  1  | Android | 1 |
|  2  | Samsung | 2 |
|  3  | iOS | 1 |
|  4  | LG | 2 |

Was taken from exercise.

3) Content

| id  | protection_system_id  | content_key | payload |
| --- | ----- |--------------------- | -------- |

Created by me.

1) id - auto_increment unique int;

2) protection_system_id - int key from "Protection systems" table;

3) content_key - VAR_CHAR encryption key to decrypt "payload";

4) payload - VAR_CHAR encrypted data, that can be decrypted by "content_key" in "encryption_mode", taked from "Protection systems" by ID.

# Api usage
1) GetContent (GET method) - endpoint to get all content from database.   

Logic:
- go to DB and get all content data. If it failed - get error;
- return HTTP 200 with all content data in JSON struct.
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

Logic:
- check "id" in URL. If not valid - get error;
- go to DB and get current content data. If it failed - get error;
- if response is null - get error.
- return HTTP 200 with content data in JSON struct.
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

Logic:
- decode request JSON. If not valid - get error;
- check struct as content data. If not valid - get error;
- check protection system. If we don't have it in DB - get error;
- check payload to decryption. If we can't decrypt it - get error;
- go to DB and add current content data. If it failed - get error;
- return HTTP 200.
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

Logic:
- check "id" in URL. If not valid - get error;
- decode request JSON. If not valid - get error;
- check struct as content data. If not valid - get error;
- check protection system if it set. If we don't have it in DB - get error;
- if new paiload is empty - get old data from DB. If it failed or null - get error;
- modern old params by new;
- check old payload to decryption. If we can't decrypt it - get error;
- check new payload to decryption. If we can't decrypt it - get error;
- go to DB and update current content data. If it failed - get error;
- return HTTP 200.
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

Logic:
- check "id" in URL. If not valid - get error;
- go to DB and delete current content data. If it failed - get error;
- return HTTP 200.
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

Logic:
- decode request JSON. If not valid - get error;
- check struct as view content data. If not valid - get error;
- check device if it set. If we don't have it in DB - get error;
- go to DB and get payload, content_key and encryption_mode from current content_id and device name. If it failed - get error;
- check that result from BD isn't null. If it's null - get error that user cant watch it;
- check that our ecryption module can work with this encryption mode. If it can't - get error;
- check payload to decryption. If we can't decrypt it - get error;
- return HTTP 200 with decrypted data.
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
# Easy demo
Here you can step-by-step try API functionality. Just copy and paste commands into terminal:

1) Build and run api-server in DEMO mode:
```
go get ./...
go build restapiserver.go
./restapiserver -demo
```
You need to little wait when docker with MySQL database will up. It's about 10 sec. After this you will see messages:
```
Checking needed utils...ON
Connecting to DataBase.............ON
```

2) Let's see what content we have in database:
```
curl -X GET "http://localhost:5000/content"
```
You will see list of all content in database:
```
[{"id":2,"protection_system_name":"AES 1","content_key":"mypassword","payload":"U2FsdGVkX1+lxfHPBsyNB+R1lJ2qOz/uA7NTprwWXhaMaQLNyhPRCyUq13VvkRDp"},{"id":1,"protection_system_name":"AES 2","content_key":"superpass","payload":"U2FsdGVkX190cOearjAhFozvAQFjW53OUhLQGKfTVZnj8iOwveiaZ8rqAPNBjeDB"},{"id":3,"protection_system_name":"AES 2","content_key":"badpass","payload":"U2FsdGVkX1+lxfJ2qOz/uA7NTprwWXhaMaQLNHPBsyNB+R1lyhPRCyUq13VvkRDp"}]
```
First and second - good data, but last one - damaged.

3) Let's add some content data:
```
curl -i -H "Content-Type: application/json" -X POST -d '{"protection_system_name":"AES 4","content_key":"popi","payload":"U2FsdGVkX18fO6a7VqCp2W2vcUGTbZqpzxJoHtR+80sy+ngb16+9OQBFPtH2aXxd"}' http://localhost:5000/content
```
We get:
```
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Thu, 30 May 2019 23:31:12 GMT
Content-Length: 38

no such protection system in database
```
That's right, because in our database we don't have such protection system with name "AES 4". Let's make it nice:
```
curl -i -H "Content-Type: application/json" -X POST -d '{"protection_system_name":"AES 2","content_key":"popi","payload":"U2FsdGVkX18fO6a7VqCp2W2vcUGTbZqpzxJoHtR+80sy+ngb16+9OQBFPtH2aXxd"}' http://localhost:5000/content
```
Result will be 200. Let's check that now we have 4 content in our database:
```
curl -X GET "http://localhost:5000/content"
```
Result:
```
[{"id":1,"protection_system_name":"AES 2","content_key":"superpass","payload":"U2FsdGVkX190cOearjAhFozvAQFjW53OUhLQGKfTVZnj8iOwveiaZ8rqAPNBjeDB"},{"id":2,"protection_system_name":"AES 1","content_key":"mypassword","payload":"U2FsdGVkX1+lxfHPBsyNB+R1lJ2qOz/uA7NTprwWXhaMaQLNyhPRCyUq13VvkRDp"},{"id":3,"protection_system_name":"AES 2","content_key":"badpass","payload":"U2FsdGVkX1+lxfJ2qOz/uA7NTprwWXhaMaQLNHPBsyNB+R1lyhPRCyUq13VvkRDp"},{"id":4,"protection_system_name":"AES 2","content_key":"popi","payload":"U2FsdGVkX18fO6a7VqCp2W2vcUGTbZqpzxJoHtR+80sy+ngb16+9OQBFPtH2aXxd"}]
```
Content with ID=4 is our new content. Let's see our new content closer:
```
curl -X GET "http://localhost:5000/content/4"
```
Result:
```
{"id":4,"protection_system_name":"AES 2","content_key":"popi","payload":"U2FsdGVkX18fO6a7VqCp2W2vcUGTbZqpzxJoHtR+80sy+ngb16+9OQBFPtH2aXxd"}
```

4) Let's try to watch some content. For example - with id=3 and on "Samsung":
```
curl -i -H "Content-Type: application/json" -X GET -d '{"content_id":3,"Device":"Samsung"}' http://localhost:5000/view
```
Result:
```
HTTP/1.1 417 Expectation Failed
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Thu, 30 May 2019 23:39:03 GMT
Content-Length: 28

invalid payload in database
```
Well, we have damaged content in our database with ID=3. Let's try to watch with content id = 4:
```
curl -i -H "Content-Type: application/json" -X GET -d '{"content_id":4,"Device":"Samsung"}' http://localhost:5000/view
```
Response:
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 30 May 2019 23:41:34 GMT
Content-Length: 20

"full_metal_jacket"
```
Okay, this is good content! Decrypted data is "full_metal_jacket" - movie from Stanley Kubrick.

5) Let's trying to fix content with id=3 - we replace old data with new, but let's use bad key:
```
curl -i -H "Content-Type: application/json" -X PUT -d '{"protection_system_name":"AES 1","content_key":"pass","payload":"U2FsdGVkX18PJILwscA+WPkF9jB+vtBMH4hjEVhQU1Wl+Zbi75xtwQuOhKVEuyEh"}' http://localhost:5000/content/3
```
Result:
```
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Thu, 30 May 2019 23:49:10 GMT
Content-Length: 16

invalid payload
```
Let's input right key:
```
curl -i -H "Content-Type: application/json" -X PUT -d '{"protection_system_name":"AES 1","content_key":"superpass","payload":"U2FsdGVkX18PJILwscA+WPkF9jB+vtBMH4hjEVhQU1Wl+Zbi75xtwQuOhKVEuyEh"}' http://localhost:5000/content/3
```
Result 200 OK. Okay, let's check our damaged content now:
```
curl -i -X GET http://localhost:5000/content/3
```
Result:
```
{"id":3,"protection_system_name":"AES 2","content_key":"superpass","payload":"U2FsdGVkX18PJILwscA+WPkF9jB+vtBMH4hjEVhQU1Wl+Zbi75xtwQuOhKVEuyEh"}
```
Looks like it's good. Let's check it on "iOS" device:
```
curl -i -H "Content-Type: application/json" -X GET -d '{"content_id":3,"Device":"iOS"}' http://localhost:5000/view
```
Response:
```
HTTP/1.1 451 Unavailable For Legal Reasons
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Thu, 30 May 2019 23:54:44 GMT
Content-Length: 23

content can't be shown
```
Well, it's because our device can't work with "AES 2" format. Let's try this content on "LG":
```
curl -i -H "Content-Type: application/json" -X GET -d '{"content_id":3,"Device":"LG"}' http://localhost:5000/view
```
And response is:
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 30 May 2019 23:54:33 GMT
Content-Length: 24

"wow_it_super_password"
```
Now we have fixed content with id=3.

6) Let's see some content data with realy big content id:
```
curl -i -X GET http://localhost:5000/content/777
```
Result:
```
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8
X-Content-Type-Options: nosniff
Date: Thu, 30 May 2019 23:59:26 GMT
Content-Length: 28

no such content in database
```
Well, this is right, because we don't have content with id = 777.

7) We've trying to get content with id=3 on our "iOS" device, that didn't work with "AES 2". Let's make it "AES 1":
```
- curl -i -H "Content-Type: application/json" -X PUT -d '{"protection_system_name":"AES 1"}' http://localhost:5000/content/3
```
Response: 200 OK. It means that now we can watch our content with id=3 on "iOS" device. Let's try:
```
curl -i -H "Content-Type: application/json" -X GET -d '{"content_id":3,"Device":"iOS"}' http://localhost:5000/view
```
Response:
```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Fri, 31 May 2019 00:05:07 GMT
Content-Length: 24

"wow_it_super_password"
```

8) We've seen our content with id=3 on different devices. Let's delete this content:
```
curl -i -X DELETE http://localhost:5000/content/3
```
Response: 200 OK. So let's see what we have now:
```
curl -i -X GET http://localhost:5000/content
```
Response:
```
[{"id":1,"protection_system_name":"AES 2","content_key":"superpass","payload":"U2FsdGVkX190cOearjAhFozvAQFjW53OUhLQGKfTVZnj8iOwveiaZ8rqAPNBjeDB"},{"id":2,"protection_system_name":"AES 1","content_key":"mypassword","payload":"U2FsdGVkX1+lxfHPBsyNB+R1lJ2qOz/uA7NTprwWXhaMaQLNyhPRCyUq13VvkRDp"},{"id":4,"protection_system_name":"AES 2","content_key":"popi","payload":"U2FsdGVkX18fO6a7VqCp2W2vcUGTbZqpzxJoHtR+80sy+ngb16+9OQBFPtH2aXxd"}]
```
So, content with id=3 was deleted. And demo is over;)

9) Press CTRL+C to determinate application. Before determinate, it will stop and remove docker container with demo MySQL database.

# Tests
To run unit-tests you need to go to folder where you can find some 'XXX_test.go' file and run:
```
go test
```

# PS
Application is working now with ALREADY ENCRYPTED data, but it's easy to working with "clean" data. Also, application working with ENCRYPTED STRINGS (was selected as most basic data for show that idea of encryption is working), but can be upgraded to working with different binary data.
Application was decomposet to modules (Database, Encryption, RESTApi, logger, config/error-file and others) to be more flexible in modernisation.
I can make some mistakes with HTTP-codes, but i think that my variant is not so bad)
