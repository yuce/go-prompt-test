package hzsqlcl

import (
	"fmt"
	"hzsqlcl/components"
	"sort"
	"strings"
)

func CreateSQLForCreateMapping(fieldPrefix string, keyValues map[string]interface{}) (string, error) {
	field := fmt.Sprintf("%sField_", fieldPrefix)
	const option = "Option_"
	const intPrefix = "Int_"
	mappingName := keyValues[components.MappingName]
	mappingType := keyValues[components.MappingType]
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
		CREATE MAPPING %s (%s) TYPE %s OPTIONS (%s)
	`, mappingName, strings.Join(fields, ", "), mappingType, strings.Join(options, ", "))), nil
}

func CreateSQLForJob(keyValues map[string]interface{}) (string, error) {
	const sinkField = "Sink_Field_"
	const sourceField = "Source_Field_"
	jobName := keyValues[components.JobName]
	sinkName := keyValues[components.SinkName]
	sourceName := keyValues[components.SourceName]
	sinkFields := fieldStrings{}
	sourceFields := fieldStrings{}
	for k, _ := range keyValues {
		if strings.HasPrefix(k, sinkField) {
			key := k[len(sinkField):]
			sinkFields = append(sinkFields, key)
		} else if strings.HasPrefix(k, sourceField) {
			key := k[len(sourceField):]
			sourceFields = append(sourceFields, key)
		}
	}
	sort.Sort(sinkFields)
	sort.Sort(sourceFields)
	return strings.TrimSpace(fmt.Sprintf(`
			CREATE JOB %s AS SINK INTO %s (%s) SELECT %s FROM %s
		`, jobName, sinkName, strings.Join(sinkFields, ", "), strings.Join(sourceFields, ", "), sourceName)), nil
}

type fieldStrings []string

func (fs fieldStrings) Len() int {
	return len(fs)
}

func (fs fieldStrings) Less(i, j int) bool {
	return strings.TrimPrefix(fs[i], "-") < strings.TrimPrefix(fs[j], "-")
}

func (fs fieldStrings) Swap(i, j int) {
	fs[i], fs[j] = fs[j], fs[i]
}
