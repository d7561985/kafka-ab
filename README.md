# Kafka benchmark
`Producer` helps produce a big number of system messages simultaneously to all partitions. Every message waits response with successful emitting to broker.

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
