buildplatform="linux/arm64"
imgname="easy-chat-user-rpc"
imgversion="v1"
Dockerfile="deploy/dockerfile/Dockerfile_user_rpc_dev"

#删除镜像
rm -f ${imgname}_${imgversion}_${buildplatform:6}.tar

docker buildx build --platform ${buildplatform} -f ${Dockerfile}  -t ${imgname}:${imgversion} . --load
#-t entrycentos
#docker run entrycentos
#docker run entrycentos -l
docker save -o  ./deploy/tar/${imgname}_${imgversion}_${buildplatform:6}.tar ${imgname}:${imgversion}