#rpc user.proto 生成
goctl rpc protoc ./apps/user/rpc/user.proto --go_out=./apps/user/rpc --go-grpc_out=./apps/user/rpc --zrpc_out=./apps/user/rpc

#user 数据库
goctl model mysql ddl -src="./deploy/sql/user.sql" -dir="./apps/user/models" -c

#api user.api 生成
goctl api go -api apps/user/api/user.api -dir apps/user/api -style gozero