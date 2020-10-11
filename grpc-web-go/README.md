Article [here](https://medium.com/@aravindhanjay/a-todo-app-using-grpc-web-and-vue-js-4e0c18461a3e)


## Generate Stub from protos
```sh
  go get -u github.com/golang/protobuf/{proto,protoc-gen-go}
  export PATH=$PATH:$GOPATH/bin
  protoc -I server/ server/todo.proto --go_out=plugins=grpc:server
```


## Use local copy of package
Init go module

```sh
   go mod github.com/oneeyedsunday/go_playground/grpc-web-go
   go run server.go
```

