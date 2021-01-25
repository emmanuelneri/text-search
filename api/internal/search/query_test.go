package search

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestRequestQuery_Query(t *testing.T) {
	type fields struct {
		value            string
		multiMatchType   string
		multiMatchFields []string
		fieldValueFactor string
		boost            string
	}
	tests := []struct {
		name         string
		fields       fields
		expectedFile string
	}{
		{
			name: "when query with value should use multi match query with relevance",
			fields: fields{
				value:            "Jadir",
				multiMatchType:   "phrase_prefix",
				multiMatchFields: []string{"Name", "Username"},
				fieldValueFactor: "relevance",
				boost:            "5",
			},
			expectedFile: "testdata/multi-match-query.json",
		},
		{
			name: "when query with empty value should use match all query with relevance",
			fields: fields{
				value:            "",
				multiMatchType:   "phrase_prefix",
				multiMatchFields: []string{"Name", "Username"},
				fieldValueFactor: "relevance",
				boost:            "5",
			},
			expectedFile: "testdata/match-all-query.json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := RequestQuery{
				Value:            tt.fields.value,
				MultiMatchType:   tt.fields.multiMatchType,
				MultiMatchFields: tt.fields.multiMatchFields,
				FieldValueFactor: tt.fields.fieldValueFactor,
				Boost:            tt.fields.boost,
			}

			file, err := ioutil.ReadFile(tt.expectedFile)
			assert.Nil(t, err)

			fileBuffer := new(bytes.Buffer)
			if err := json.Compact(fileBuffer, file); err != nil {
				assert.Nil(t, err)
			}

			marshaledQuery, err := json.Marshal(s.Query())
			assert.Nil(t, err)

			assert.Equal(t, fileBuffer.String(), string(marshaledQuery))
		})
	}
}
