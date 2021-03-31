package hzsqlcl

import (
	"fmt"
	"strings"
)

func CreateSQLForCreateMapping(keyValues map[string]interface{}) (string, error) {
	const field = "Field_"
	const option = "Option_"
	const intPrefix = "Int_"
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

func CreateSQLForJob(keyValues map[string]interface{}) (string, error) {
	const sinkField = "Sink_Field_"
	const sourceField = "Source_Field_"
	jobName := keyValues[JobName]
	sinkName := keyValues[SinkName]
	sourceName := keyValues[SourceName]
	sinkFields := []string{}
	sourceFields := []string{}
	for k, _ := range keyValues {
		if strings.HasPrefix(k, sinkField) {
			k = k[len(sinkField):]
			sinkFields = append(sinkFields, fmt.Sprintf("%s", k))
		} else if strings.HasPrefix(k, sourceField) {
			k = k[len(sourceField):]
			sourceFields = append(sourceFields, fmt.Sprintf("%s", k))
		}
	}
	return strings.TrimSpace(fmt.Sprintf(`
		CREATE JOB %s AS SINK INTO %s (%s) SELECT %s FROM %s;
	`, jobName, sinkName, strings.Join(sinkFields, ", "), strings.Join(sourceFields, ", "), sourceName)), nil
}
