<?php
	$POST = json_decode(file_get_contents("php://input"), true);
    $endpoint = $POST["endpoint"];
    $data = $POST["data"];

	echo print_r($endpoint) . PHP_EOL;
	echo print_r($data) . PHP_EOL;
?>