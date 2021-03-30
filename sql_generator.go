package hzsqlcl

import (
	"fmt"
	"strings"
)

func CreateSQLForCreateMapping(keyValues map[string]interface{}) (string, error) {
	mappingName := keyValues[MappingName]
	mappingType := keyValues[MappingType]
	return strings.TrimSpace(fmt.Sprintf(`
		CREATE MAPPING %s (key INT, name VARCHAR, age INT) TYPE %s OPTIONS ('valueFormat' = 'json', 'bootstrap.servers' = '127.0.0.1:9092');
	`, mappingName, mappingType)), nil
}
