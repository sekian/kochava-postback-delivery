# Kochava Project :: Postback Delivery

## Description
A service to function as a small scale simulation of how Kochava distributes data to third parties in real time. 

## App Operation (System Design)
The project consists of a php application to ingest http requests and a go application to deliver http responses with a kafka job queue between them.

The interaction flow between the components is the following:

**Web request &rarr; Ingestion Agent (php) &rarr; Delivery Queue (kafka) &rarr; Delivery Agent (go) &rarr; Web response**

The detailed app operation is:

1) Ingestion agent accepts incoming http request
2) Ingestion agent pushes a postback object to the kafka  delivery queue for each "data" object contained in accepted request.
3) Delivery agent continuously pulls postback objects from kafka delivery queue.
4) Delivery agent delivers each postback object to the http endpoint after having filled the query parameters with the data values.
5) Delivery agent logs delivery time, response code, response time, and response body.

## Configuration
The configuration commands are for Windows 10. The commands for a Linux machine will be similar. 
The configuration of constants will be on the header of the relevant files.

### Ingestion Agent (php):

Run `composer update` on the ingestion-agent folder to retrieve all the dependencies.

To setup the PHP server two options were tested:

1. Download and unzip [PHP](https://www.php.net/downloads). It may be useful to add it to the PATH environtment variable. You can then startup the server with:

`php -S localhost:80 -t "$PWD"`

2. Download and install [XAMPP](https://www.apachefriends.org/download.html). Run the XAMPP Control Panel and start the Apache module. Using XAMPP will support logging and other useful features out of the box. By default the PHP code will be on "C:\xampp\htdocs". If you wish to change the default directory on the XAMPP app go to `Config->httpd.conf` and change the DocumentRoot and Directory:

```
DocumentRoot "C:/xampp/htdocs"
<Directory "C:/xampp/htdocs">
```

The PHP versions used are 8.1.6 with XAMPP v3.3.0 and 8.1.10 with the regular PHP.

### Delivery Queue (kafka):
Download and unzip [Apache Kafka](https://kafka.apache.org/downloads). Go to the unzipped Kafka folder and run the following commands:

`bin/windows/zookeeper-server-start.bat config/zookeeper.properties`

`bin/windows/kafka-server-start.bat config/server.properties`

`bin/windows/kafka-topics.bat --create --bootstrap-server localhost:9092 --replication-factor 1 --partitions 1 --topic postback`

`bin/windows/kafka-topics.bat --list --bootstrap-server localhost:9092`

The Apache Kafka version used is kafka_2.13-3.2.1

### Delivery Agent (go):
Download and install [Golang](https://go.dev/doc/install)

Run `go get -u github.com/segmentio/kafka-go` on the delivery-agent folder to install the needed dependency

To startup the delivery agent run `go run deliver.go`. By default the service will log to `logfile.log`

The used Go version is 1.19

## Using the Postback service

After setting up everything you can send a POST request to the php ingest agent:

`curl "localhost:80/ingest.php" -H 'Content-Type: application/json' -X POST -d '@request_weather.json'`

You can also use Postman or any other software to send the request to the php ingest agent.
