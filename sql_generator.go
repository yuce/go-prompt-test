package hzsqlcl

import "strings"

func CreateSQLForCreateMapping(keyValues map[string]interface{}) (string, error) {
	return strings.TrimSpace(`
		CREATE MAPPING myJsonTopic (key INT, name VARCHAR, age INT) TYPE Kafka OPTIONS ('valueFormat' = 'json', 'bootstrap.servers' = '127.0.0.1:9092');
	`), nil
}
