# Kafka-AB
`Producer` helps produce a big number of system messages simultaneously to all partitions. Every message waits response with successful emitting to broker.

`Consumer` - read provided topic with unit group per thread with auto-commit option.

## usage
application follows the best CLI practises, so it contains help information about environment, supported, etc.

### go run / build / install
```bash
$ go run main.go c --topic my-topic --kafka-server 127.0.0.1:9094 --threads 10
```

### docker image
```bash
$ kubectl -n global run consumer --image d7561985/kafka-bench:v1.0.0 --env KAFKA_SERVER=my-cluster-kafka-external-bootstrap:9094 -it --rm  c
$ kubectl -n global run producer --image d7561985/kafka-bench:v1.0.0 --env KAFKA_SERVER=my-cluster-kafka-external-bootstrap:9094 -it --rm  p
```

### from what to start
assume that we have topic with 50 partitions and some replicas
```yaml
apiVersion: kafka.strimzi.io/v1beta1
kind: KafkaTopic
metadata:
  name: my-topic
  labels:
    strimzi.io/cluster: my-cluster
spec:
  partitions: 50
  replicas: 2
  config:
    retention.ms: 3600000
    segment.bytes: 1073741824
```
for test topic try to use small retention time as log is filled very significantly


### Use cases 
Application exits with status 1 when it can't meet some target. 

For example: consumer don't read any message or consumer not reads desired requests number which helps to perform some CI tests ;) 
#### Consumer
* `--timelimit` if any of consumer not read any message at the and of limit will exit with status 1
* `--requests` exit with status 0 when request meat desired number, but if timelimit comes - status will be 1  
* Concurrent (`--concurrency`) consumer with the force name (`--force-name`) - one group reads from different  partitions
* `--static` every run consumers will have the same name. This helps to check auto-commit option for example
* ` --earliest` option allow to deside do you need read from the start of just new income. Set false to read new income only. 