# Dungeonforge

## Compile Proto

Go to `/src/lobby/proto` and run `protoc --go_out=.. --go_opt=paths=source_relative --go-grpc_out=.. --go-grpc_opt=paths=source_relative lob
by.proto --experimental_allow_proto3_optional`