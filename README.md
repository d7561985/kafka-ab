# Kafka benchmark
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
