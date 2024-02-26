# protobuf

need to install protocol buffer compiler.  
You don't need to run it, but if you want to compile it yourself, refer to the following procedure.

protobufのコンパイラをインストールする必要があります。  
実行しなくても、コンパイル済みのファイルを含めていますが、自分でコンパイルする場合は以下の手順を参考にしてください。  

[Protocol Buffer Compiler Installation](https://grpc.io/docs/protoc-installation/)

```bash
$ protoc --go_out=../ --proto_path=. --go_opt=module=github.com/ytake/protoactor-go-cqrs-example *.proto 
```
