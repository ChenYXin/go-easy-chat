#
# user rpc
# 构建镜像并推送到阿里云
#

#可选择的操作系统
#linux/arm64, linux/amd64, linux/amd64/v2, linux/riscv64, linux/ppc64le, linux/s390x, linux/386, linux/mips64le, linux/mips64, linux/arm/v7, linux/arm/v6
#我的mac电脑linux/arm64
#阿里云服务linux/amd64
buildplatform="linux/arm64"
imgname="easyChatUserRpc"
imgversion="v1"
Dockerfile="deploy/dockerfile/Dockerfile_user_rpc_dev"

#删除镜像
rm -f ${imgname}_${imgversion}_${buildplatform:6}.tar

docker buildx build --platform ${buildplatform} -f ${Dockerfile}  -t ${imgname}:${imgversion} . --load
#-t entrycentos
#docker run entrycentos
#docker run entrycentos -l
docker save -o  ./deploy/tar/${imgname}_${imgversion}_${buildplatform:6}.tar ${imgname}:${imgversion}
echo "构建完成"
echo "镜像操作系统：${buildplatform}"
echo "指定镜像名称：${imgname}:${imgversion}"
echo "本地tar名称：${imgname:7}_${imgversion}_${buildplatform:6}.tar"

#printf "是否推送到阿里云(默认不推送)\n0.不推送\n1.推送\n: "
#read -r -n 1 push_aliyun
#if [[ ${push_aliyun} -eq 1 ]];then
#   docker login --username=xx --password=yy registry.cn-shenzhen.aliyuncs.com
#   docker tag ${imgname}:${imgversion} registry.cn-shenzhen.aliyuncs.com/donkor/aa:${imgversion}
#   docker push registry.cn-shenzhen.aliyuncs.com/donkor/aa:${imgversion}
#fi