<?php
    require 'vendor/autoload.php';
    const BROKER_URI = '127.0.0.1:9092';
    const TOPIC = 'postback';
    const VERSION = '3.2.1';
    const REFRESH_MS = 10000;
    
	$POST = json_decode(file_get_contents("php://input"), true);
    
    if(!isset($POST['endpoint']) || !isset($POST['data'])) {
        echo 'ERROR: Message body is missing "endpoint" or "data"';
        exit();
    }
    
    $endpoint = $POST['endpoint'];
    $data = $POST['data'];
	// var_dump($endpoint) . PHP_EOL;
	// var_dump($data) . PHP_EOL;
    
    $config = \Kafka\ProducerConfig::getInstance();
    $config->setMetadataRefreshIntervalMs(REFRESH_MS);
    $config->setMetadataBrokerList(BROKER_URI);
    $config->setRequiredAck(1);
    $config->setIsAsyn(false);
    $config->setProduceInterval(500);
    $producer = new \Kafka\Producer();
    
    foreach((array) $data as $value) {
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