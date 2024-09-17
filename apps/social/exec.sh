#rpc social.proto 生成
goctl rpc protoc ./apps/social/rpc/social.proto --go_out=./apps/social/rpc --go-grpc_out=./apps/social/rpc --zrpc_out=./apps/social/rpc

#social 数据库
goctl model mysql ddl -src="./deploy/sql/social.sql" -dir="./apps/social/socialmodels" -c

#api social.api 生成
goctl api go -api apps/social/api/social.api -dir apps/social/api -style gozero