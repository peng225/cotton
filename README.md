# Cotton

A simple in-memory blob server.

## How to build

```
make
```

## How to use

### Basic usage

1. Start the cotton server.
   ```
   $ ./cotton       
   server.go:37: Start server. port = 8080

   ```

2. Post data object.
   ```
   $ curl -i -X POST localhost:8080/test/data --data-binary @test.txt
   HTTP/1.1 201 Created
   Location: /test/data/3905f7d8-852f-4df3-bd8c-2fbe8e54c01a
   Date: Sun, 27 Nov 2022 13:01:27 GMT
   Content-Length: 0

   ```

3. Now you can get the posted data object.
   ```
   $ curl -i localhost:8080/test/data/3905f7d8-852f-4df3-bd8c-2fbe8e54c01a
   HTTP/1.1 200 OK
   Content-Length: 47
   Date: Sun, 27 Nov 2022 13:02:30 GMT
   Content-Type: text/plain; charset=utf-8
   
   This is a test file.
   1 + 1 = 2
   foo@example.com
   ```

4. You can delete the object as follows.
   ```
   $ curl -i -X DELETE localhost:8080/test/data/3905f7d8-852f-4df3-bd8c-2fbe8e54c01a
   HTTP/1.1 200 OK
   Date: Sun, 27 Nov 2022 13:05:14 GMT
   Content-Length: 0
   
   $ curl -i localhost:8080/test/data/3905f7d8-852f-4df3-bd8c-2fbe8e54c01a 
   HTTP/1.1 404 Not Found
   Date: Sun, 27 Nov 2022 13:05:22 GMT
   Content-Length: 0
   
   ```

### Encoding

Cotton supports only gzip for Content-Encoding. Note that the gzip encoding is not applicable for HEAD requests.

```
$ curl -i -X POST localhost:8080/test/data --data-binary @test.txt
HTTP/1.1 201 Created
Location: /test/data/bb8e3ecd-3cfb-452e-9b1b-0f7c27ac7cb4
Date: Thu, 01 Dec 2022 12:25:09 GMT
Content-Length: 0

$ curl -i --compressed localhost:8080/test/data/bb8e3ecd-3cfb-452e-9b1b-0f7c27ac7cb4 
HTTP/1.1 200 OK
Content-Encoding: gzip
Date: Thu, 01 Dec 2022 12:25:49 GMT
Transfer-Encoding: chunked

This is a test file.
1 + 1 = 2
foo@example.com
$ curl -i -I --compressed localhost:8080/test/data/bb8e3ecd-3cfb-452e-9b1b-0f7c27ac7cb4
HTTP/1.1 200 OK
Content-Length: 47
Date: Thu, 01 Dec 2022 12:26:31 GMT

```

## Spec

- Data objects are stored in memory. All objects are deleted when the cotton server stops.
- An unique UUID is assigned to each object.
- Maximum data size is 10MiB.
- Maximum path length in URL is 1024 including the UUID part.
- Supported methods are as follows:
  - GET
  - HEAD
  - POST
  - PUT
  - DELETE
