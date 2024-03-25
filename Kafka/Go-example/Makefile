.PHONY: producer consumer

PUB = Producer
SUB = Consumer

producer:
	go run $(PUB)/util.go $(PUB)/producer.go getting-started.properties

consumer:
	go run $(SUB)/util.go $(SUB)/consumer.go getting-started.properties
