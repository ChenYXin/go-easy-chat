#chatmodel 数据库
goctl model mongo --type chatLog --dir ./apps/im/immodels
# 会话模型
goctl model mongo --type conversations --dir ./apps/im/immodels
# 会话列表
goctl model mongo --type conversation --dir ./apps/im/immodels

#rpc im.proto 生成
goctl rpc protoc ./apps/im/rpc/im.proto --go_out=./apps/im/rpc --go-grpc_out=./apps/im/rpc --zrpc_out=./apps/im/rpc
