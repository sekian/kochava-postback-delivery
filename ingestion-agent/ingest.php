<?php
    require 'vendor/autoload.php';
    const BROKER_URL = '127.0.0.1:9092'; // Kafka URL
    const TOPIC = 'postback'; // Kafka topic
    const VERSION = '3.2.1'; // Kafka version
    const REFRESH_MS = 10000;
    
	$POST = json_decode(file_get_contents("php://input"), true); // Read request body
    
    if(!isset($POST['endpoint']) || !isset($POST['data'])) { // Check endpoint and data present
        echo 'ERROR: Message body is missing "endpoint" or "data"';
        exit();
    }
    
    $endpoint = $POST['endpoint'];
    $data = $POST['data'];
	// var_dump($endpoint) . PHP_EOL;
	// var_dump($data) . PHP_EOL;
    
    $config = \Kafka\ProducerConfig::getInstance(); // Kafka producer setup
    $config->setMetadataRefreshIntervalMs(REFRESH_MS);
    $config->setMetadataBrokerList(BROKER_URL);
    $config->setRequiredAck(1);
    $config->setIsAsyn(false);
    $config->setProduceInterval(500);
    $producer = new \Kafka\Producer();
    
    foreach((array) $data as $value) { // Push postback to Kafka for each data object
        $value['method'] = $endpoint['method'] ?? null;
        $value['url'] = $endpoint['url'] ?? null;
        print_r($value) . PHP_EOL;
        $producer->send([
            [
                'topic' => TOPIC,
                'value' => json_encode($value),
                'key' => '',
            ],
        ]);
    }
?>