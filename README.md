# HTTP Server from Scratch in Go

## Features

- Accepting incoming TCP packets
- Reading the request string line by line
- Parsing the request string into headers and body
- Passing the request to a handler function
  - **/echo/{foo}**: prints the request back to the client
    - if contains body and valid Accept-Encoding header, compresses the body
  - **/files/{foo}**: serves files from the filesystem
    - if POST, saves the body to a file
- Handles Concurrent Connections

## Usage

`cd app && go run .`

Use a http client such as *curl* to interact with the server.

## Testing

- `cd app && go test -v .`

## Further Improvements

- Add SSL support (HTTPS)
- Add CRUD capabilities (PUT, DELETE)
- Add query parameter support (?foo=bar)
- Add JSON support
- Add logging
