# grpc-playground
Playground for GRPC server and client capabilities

# generate proto files

Before running any of the examples, we need to generate the proto stubs.
We can use buf.build for this. 

## installing buf.build:

Install as instructed for your platform: [Buf.Build installation](https://docs.buf.build/installation)

## Generating files:

```
make protogen
```
This will generate the files in `api` folder.

# Running go gRPC server for error details:

```
go run echo/servererrors/servererrors.go
```