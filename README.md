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

Cotton supports only gzip for Content-Encoding.

```
$ curl -i -X POST localhost:8080/test/data --data-binary @test.txt
HTTP/1.1 201 Created
Location: /test/data/67b72b01-0b2a-4f4b-b4e9-e06c90903ae3
Date: Sun, 27 Nov 2022 13:09:48 GMT
Content-Length: 0

$ curl -i --compressed localhost:8080/test/data/67b72b01-0b2a-4f4b-b4e9-e06c90903ae3
HTTP/1.1 200 OK
Content-Encoding: gzip
Content-Length: 71
Date: Sun, 27 Nov 2022 13:10:49 GMT

This is a test file.
1 + 1 = 2
foo@example.com
$ curl -i -I --compressed localhost:8080/test/data/67b72b01-0b2a-4f4b-b4e9-e06c90903ae3
HTTP/1.1 200 OK
Content-Encoding: gzip
Content-Length: 71
Date: Sun, 27 Nov 2022 13:14:31 GMT

```

## Spec

- Data objects are stored in memory. All objects are deleted when the cotton server stops.
- An unique UUID is assigned to each object.
- Maximum data size is 10MiB.
- Maximum path length in URL is 1024 including the UUID part.

## Limitations

- PUT method is not supported, which means that you cannot update the posted data.
- TLS is not supported.