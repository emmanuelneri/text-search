package search

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestDecodeResult(t *testing.T) {
	file, err := ioutil.ReadFile("testdata/result-query.json")
	assert.Nil(t, err)

	result := ElasticResult{}
	err = json.Unmarshal(file, &result)
	assert.Nil(t, err)

	assert.Equal(t, "FGluY2x1ZGVfY29udGV4dF91dWlkDXF1ZXJ5QW5kRmV0Y2gBFkJySEhCMnR5U1ZpQS16STBpRmFHMXcAAAAAAAAAIhY5QW81aWx1clE4dWxnMzlFZDI2NEln", result.ScrollId)
	assert.Equal(t, 4, result.Took)

	hits := result.Hits
	assert.NotNil(t, hits)

	total := result.Total
	assert.NotNil(t, total)
	assert.Equal(t, 10000, total.Value)

	subHits := result.Hits
	assert.NotNil(t, subHits)

	expectedSource := []map[string]string{
		{"ID": "36e9f58f-c4e0-4c41-84a4-81ac5e608553", "Name": "Tamiris Kesheh", "Username": "tamiris.kesheh"},
		{"ID": "e6da5140-fae5-435e-868d-9c4e4e7a465e", "Name": "Elizabeth Casali", "Username": "elizabethcasali"},
		{"ID": "1daac1c8-6393-4649-a115-31bf7c4b8e9d", "Name": "Eberson Belmiro Raimundo", "Username": "ebersonbelmiroraimundo"},
		{"ID": "9ec3783c-80ac-4a0c-8226-c5b15f4490d3", "Name": "Dieison Sikvs", "Username": "dieisonsikvs"},
		{"ID": "fd383b52-8e01-4b1b-88cd-c534db56152a", "Name": "Manuela Agnoletto Cortepace", "Username": "manuela.agnoletto.cortepace"},
		{"ID": "87edc52d-6243-463a-a47f-0d998d0970f2", "Name": "Lys Dall", "Username": "lys.dall"},
		{"ID": "66e00f3a-56fd-40d1-9805-04a3239b803c", "Name": "mariana Mattos", "Username": "mariana.mattos"},
		{"ID": "dc917876-fc07-4479-9849-875b42c70fd3", "Name": "James Baby Bioto", "Username": "james.baby.bioto"},
		{"ID": "c3745bd2-cd97-413d-b01a-a7a6d47f896e", "Name": "Deysiane Barbato", "Username": "deysiane.barbato"},
		{"ID": "697e6415-dbb5-4e46-ac9d-a3d381d04de6", "Name": "juh Deo", "Username": "juh.deo"},
	}

	for i := 0; i < 10; i++ {
		assert.NotNil(t, subHits)
		assert.Equal(t, "index", subHits[i].Index)
		assert.Equal(t, "_doc", subHits[i].Type)
		assert.Equal(t, 1.0, subHits[i].Score)

		user := subHits[i].Source
		assert.NotNil(t, user)
		assert.Equal(t, expectedSource[i]["ID"], user["ID"])
		assert.Equal(t, expectedSource[i]["Name"], user["Name"])
		assert.Equal(t, expectedSource[i]["Username"], user["Username"])
	}
}
