# go-mail-sender-kafka-example
Example of mail sender.

Includes:    
- logrus;  
- Apache Kafka;  
- Docker;
- Docker Compose.

## Description

Sample project for educational purposes.  
There are 2 ideas:  
- checking and glueing together technologies mentioned above;
- sharing my own experience for the ones who want to glue the same technologies.

### Running

1) Use `docker-compose` to run the containers with Kafka, Zookeeper and Maildev:  
- `docker-compose build`  
- `docker-compose up`  
  
2) Use `go run .` from the folder that contains `main.go` to run the example.