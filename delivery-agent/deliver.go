package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	topic            = "postback" // Kafka topic
	groupID          = "postback"
	kafkaURL         = "localhost:9092"
	logToFile        = true
	logFilename      = "logfile.log"
	unmatchedDefault = "" // default value for unmatched parameters
	maxRetryAttempt  = 3
	retryWaitTime    = 1 // seconds
)

func makeRequest(method string, url string) (int, string) {
	log.Printf("Make request %s %s\n", method, url)
	fmt.Printf("Make request %s %s\n", method, url)

	client := http.Client{}
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		log.Println("Request error: ", err)
		return -1, ""
	}

	response, err := client.Do(request)
	if err != nil {
		log.Println("Request error: ", err)
		return -1, ""
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("Request error: ", err)
		return -1, ""
	}

	return response.StatusCode, string(body)
}

func replaceUnmatchedByDefault(endpoint_url string) string {
	regexBrackets := regexp.MustCompile(`{.*}`) // Match substrings between curly brackets
	endpoint_url = regexBrackets.ReplaceAllString(endpoint_url, unmatchedDefault)
	return endpoint_url
}

func replaceUrlParameters(encoded_postback []uint8) (string, string) {
	var postback map[string]interface{}
	err := json.Unmarshal(encoded_postback, &postback)
	if err != nil {
		log.Println("failed to unmarshal", err)
	}
	endpoint_method := postback["method"].(string)
	endpoint_url := postback["url"].(string)
	for key, value := range postback { // Replace url parameters with postback value
		old_value := "{" + key + "}"
		new_value := url.QueryEscape(value.(string))
		num_repl := -1 // Unlimited replacements
		endpoint_url = strings.Replace(endpoint_url, old_value, new_value, num_repl)
	}
	return endpoint_method, endpoint_url
}

func logf(msg string, a ...interface{}) {
	fmt.Printf(msg, a...)
	fmt.Println()
}

func getKafkaReader(kafkaURL, topic, groupID string) *kafka.Reader {
	brokers := strings.Split(kafkaURL, ",")
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  groupID,
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
		// Logger:      kafka.LoggerFunc(logf),
		ErrorLogger: kafka.LoggerFunc(logf),
	})
}

func consume() {
	reader := getKafkaReader(kafkaURL, topic, groupID)
	defer reader.Close()
	for { // Consumer infinite loop
		startDelivery := time.Now()
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Println(err)
		}
		fmt.Printf("Message at topic:%v partition:%v offset:%v key:%s value:%s\n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))

		method, url := replaceUrlParameters(msg.Value)
		url = replaceUnmatchedByDefault(url)

		startRequest := time.Now()
		statusCode, body := makeRequest(method, url)
		for i := 0; (statusCode != http.StatusOK) && (i < maxRetryAttempt); i++ {
			statusCode, body = makeRequest(method, url)
			time.Sleep(retryWaitTime * time.Second)
		}

		responseTime := time.Since(startRequest)
		deliveryTime := time.Since(startDelivery)
		log.Printf("Delivery time: %6fs| Response code: %v| Response time: %6fs| Response body: %s\n",
			deliveryTime.Seconds(), statusCode, responseTime.Seconds(), body)
	}
}

func setLogToFile() {
	file, err := os.OpenFile(logFilename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	log.SetOutput(file)
}

func main() {
	if logToFile {
		setLogToFile()
	} else {
		log.SetOutput(io.Discard)
	}
	consume()
}
