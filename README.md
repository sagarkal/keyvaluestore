# Simple key value store implemented using Standard library in Golang

# Prerequisites
Golang version >= 1.15.1 needs to be set up on your machine

## To run Server

```shell
cd server
go run server.go
```

## To run client
```shell
cd client
go run client.go
```

## Useful Info
- Key and Value can be of any primitive type
- Unsupported commands will receive an error response
- Multiple clients can simultaneously access the server, potentially the same key 


## Example commands
- SET B 8
- GET B
- DELETE B

