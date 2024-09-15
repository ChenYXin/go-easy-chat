#
# 运行镜像容器
#
read -r -p "请输入容器名称: " container_name
read -r -p "请输入容器映射端口: " container_port
read -r -p "请输入镜像ID: " image_id
docker run -itd --name ${container_name} \
-p ${container_port}:8080 \
-v /Users/chenyuexin/Desktop/go/demo/gindemo/config.yaml:/app/config.yaml \
${image_id}

#停止并删除容器
#docker stop redis | xargs docker rm
