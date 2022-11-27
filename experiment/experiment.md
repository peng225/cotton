# Behavior investigation of curl's data-specifying options

## -d, --data, --data-ascii

Newline characters are removed.

```
$ curl -XPOST -v localhost:8080/test/data -d @test.txt
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /test/data HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> Content-Length: 44
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 44 out of 44 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 201 Created
< Location: /test/data/dbb6e5b0-cc0c-401a-8dee-6f24b2c5156b
< Date: Sun, 27 Nov 2022 09:25:22 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
$ curl -v localhost:8080/test/data/dbb6e5b0-cc0c-401a-8dee-6f24b2c5156b
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /test/data/dbb6e5b0-cc0c-401a-8dee-6f24b2c5156b HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Length: 44
< Date: Sun, 27 Nov 2022 09:25:42 GMT
< Content-Type: text/plain; charset=utf-8
< 
* Connection #0 to host localhost left intact
This is a test file.1 + 1 = 2foo@example.com%                                                                                                      
$ curl localhost:8080/test/data/dbb6e5b0-cc0c-401a-8dee-6f24b2c5156b | hexdump -C
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    44  100    44    0     0  14666      0 --:--:-- --:--:-- --:--:-- 14666
00000000  54 68 69 73 20 69 73 20  61 20 74 65 73 74 20 66  |This is a test f|
00000010  69 6c 65 2e 31 20 2b 20  31 20 3d 20 32 66 6f 6f  |ile.1 + 1 = 2foo|
00000020  40 65 78 61 6d 70 6c 65  2e 63 6f 6d              |@example.com|
0000002
```

## --data-raw

The special character `@` is ignored.

```
$ curl -XPOST -v localhost:8080/test/data --data-raw @test.txt                 
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /test/data HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> Content-Length: 9
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 9 out of 9 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 201 Created
< Location: /test/data/92e90e26-24a5-4bc9-a5cb-3436d61c38b9
< Date: Sun, 27 Nov 2022 09:29:03 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
$ curl -v localhost:8080/test/data/92e90e26-24a5-4bc9-a5cb-3436d61c38b9          
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /test/data/92e90e26-24a5-4bc9-a5cb-3436d61c38b9 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Length: 9
< Date: Sun, 27 Nov 2022 09:29:21 GMT
< Content-Type: text/plain; charset=utf-8
< 
* Connection #0 to host localhost left intact
@test.txt%                                                                                                                                         
$ curl localhost:8080/test/data/92e90e26-24a5-4bc9-a5cb-3436d61c38b9 | hexdump -C
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100     9  100     9    0     0   1800      0 --:--:-- --:--:-- --:--:--  2250
00000000  40 74 65 73 74 2e 74 78  74                       |@test.txt|
00000009
```

## --data-binary

File contents are sent as binary data (without any changes).

```
% curl -XPOST -v localhost:8080/test/data --data-binary @test.txt
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /test/data HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> Content-Length: 47
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 47 out of 47 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 201 Created
< Location: /test/data/69b8ae0c-97da-4b9c-b94e-04ed52bd3f6d
< Date: Sun, 27 Nov 2022 10:23:00 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
$ curl -v localhost:8080/test/data/69b8ae0c-97da-4b9c-b94e-04ed52bd3f6d
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /test/data/69b8ae0c-97da-4b9c-b94e-04ed52bd3f6d HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Length: 47
< Date: Sun, 27 Nov 2022 10:23:16 GMT
< Content-Type: text/plain; charset=utf-8
< 
This is a test file.
1 + 1 = 2
foo@example.com
* Connection #0 to host localhost left intact
$ curl localhost:8080/test/data/69b8ae0c-97da-4b9c-b94e-04ed52bd3f6d | hexdump -C
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    47  100    47    0     0  23500      0 --:--:-- --:--:-- --:--:-- 23500
00000000  54 68 69 73 20 69 73 20  61 20 74 65 73 74 20 66  |This is a test f|
00000010  69 6c 65 2e 0a 31 20 2b  20 31 20 3d 20 32 0a 66  |ile..1 + 1 = 2.f|
00000020  6f 6f 40 65 78 61 6d 70  6c 65 2e 63 6f 6d 0a     |oo@example.com.|
0000002f
```

## --data-urlencode

File contents are converted into the URL-encoded form.

```
$ curl -XPOST -v localhost:8080/test/data --data-urlencode @test.txt
Note: Unnecessary use of -X or --request, POST is already inferred.
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /test/data HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> Content-Length: 75
> Content-Type: application/x-www-form-urlencoded
> 
* upload completely sent off: 75 out of 75 bytes
* Mark bundle as not supporting multiuse
< HTTP/1.1 201 Created
< Location: /test/data/46a05c08-8951-4e28-a4f3-7d7e0f6d19df
< Date: Sun, 27 Nov 2022 10:24:11 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
$ curl -v localhost:8080/test/data/46a05c08-8951-4e28-a4f3-7d7e0f6d19df
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /test/data/46a05c08-8951-4e28-a4f3-7d7e0f6d19df HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Length: 75
< Date: Sun, 27 Nov 2022 10:24:27 GMT
< Content-Type: text/plain; charset=utf-8
< 
* Connection #0 to host localhost left intact
This%20is%20a%20test%20file.%0A1%20%2B%201%20%3D%202%0Afoo%40example.com%0A%                                                                       
$ curl localhost:8080/test/data/46a05c08-8951-4e28-a4f3-7d7e0f6d19df | hexdump -C
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100    75  100    75    0     0   9375      0 --:--:-- --:--:-- --:--:-- 10714
00000000  54 68 69 73 25 32 30 69  73 25 32 30 61 25 32 30  |This%20is%20a%20|
00000010  74 65 73 74 25 32 30 66  69 6c 65 2e 25 30 41 31  |test%20file.%0A1|
00000020  25 32 30 25 32 42 25 32  30 31 25 32 30 25 33 44  |%20%2B%201%20%3D|
00000030  25 32 30 32 25 30 41 66  6f 6f 25 34 30 65 78 61  |%202%0Afoo%40exa|
00000040  6d 70 6c 65 2e 63 6f 6d  25 30 41                 |mple.com%0A|
0000004b
```

## -F

This option emulates posting form data on browsers.

Because cotton does not support `Content-Type: multipart/form-data`,
the body of posted data is stored as it is without interpreting the boundary.

```
$ curl -v -F "file=@test.txt;filename=test.txt" -F "file2=@test2.txt" http://localhost:8080/test/data/ 
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> POST /test/data/ HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> Content-Length: 423
> Content-Type: multipart/form-data; boundary=------------------------9003121934b793e7
> 
* We are completely uploaded and fine
* Mark bundle as not supporting multiuse
< HTTP/1.1 201 Created
< Location: /test/data/daacecac-2f11-439e-a73a-bb9beed770d9
< Date: Sun, 27 Nov 2022 10:26:10 GMT
< Content-Length: 0
< 
* Connection #0 to host localhost left intact
$ curl -v http://localhost:8080/test/data/daacecac-2f11-439e-a73a-bb9beed770d9
*   Trying 127.0.0.1:8080...
* TCP_NODELAY set
* Connected to localhost (127.0.0.1) port 8080 (#0)
> GET /test/data/daacecac-2f11-439e-a73a-bb9beed770d9 HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/7.68.0
> Accept: */*
> 
* Mark bundle as not supporting multiuse
< HTTP/1.1 200 OK
< Content-Length: 423
< Date: Sun, 27 Nov 2022 10:26:26 GMT
< Content-Type: text/plain; charset=utf-8
< 
--------------------------9003121934b793e7
Content-Disposition: form-data; name="file"; filename="test.txt"
Content-Type: text/plain

This is a test file.
1 + 1 = 2
foo@example.com

--------------------------9003121934b793e7
Content-Disposition: form-data; name="file2"; filename="test2.txt"
Content-Type: text/plain

This is a test2 file.
1 + 2 = 3
bar@example.com

--------------------------9003121934b793e7--
* Connection #0 to host localhost left intact
```

## -T

This option is used to upload multiple files in a single command. If this option is used, PUT request is issued.

Because cotton does not support PUT method, I just show the command example here.

```
$ curl -v -T "{test.txt,test2.txt}" http://localhost:8080/test/data
```