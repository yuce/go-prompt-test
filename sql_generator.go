package hzsqlcl

import (
	"fmt"
	"strings"
)

func CreateSQLForCreateMapping(keyValues map[string]interface{}) (string, error) {
	const field = "Field_"
	mappingName := keyValues[MappingName]
	mappingType := keyValues[MappingType]
	fields := []string{}
	for k, v := range keyValues {
		if strings.HasPrefix(k, field) {
			k = k[len(field):]
			fields = append(fields, fmt.Sprintf("%s %s", k, v))
		}
	}
	return strings.TrimSpace(fmt.Sprintf(`
		CREATE MAPPING %s (%s) TYPE %s OPTIONS ('valueFormat' = 'json', 'bootstrap.servers' = '127.0.0.1:9092');
	`, mappingName, strings.Join(fields, ", "), mappingType)), nil
}
