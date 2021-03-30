package hzsqlcl

import (
	"fmt"
	"strings"
)

func CreateSQLForCreateMapping(keyValues map[string]interface{}) (string, error) {
	const field = "Field_"
	const option = "Option_"
	mappingName := keyValues[MappingName]
	mappingType := keyValues[MappingType]
	fields := []string{}
	options := []string{}
	for k, v := range keyValues {
		if strings.HasPrefix(k, field) {
			k = k[len(field):]
			fields = append(fields, fmt.Sprintf("%s %s", k, v))
		} else if strings.HasPrefix(k, option) {
			k = k[len(option):]
			options = append(options, fmt.Sprintf("'%s' = '%s'", k, v))
		}
	}
	return strings.TrimSpace(fmt.Sprintf(`
		CREATE MAPPING %s (%s) TYPE %s OPTIONS (%s);
	`, mappingName, strings.Join(fields, ", "), mappingType, strings.Join(options, ", "))), nil
}
