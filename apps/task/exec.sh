#指定主题发送消息
docker exec -it kafka kafka-console-producer.sh --broker-list 192.168.1.102:9092 --topic msgChatTransfer