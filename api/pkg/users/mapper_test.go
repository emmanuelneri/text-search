package users

import (
	"api/internal/search"
	"math/rand"
	"reflect"
	"testing"
)

func Test_toPaged(t *testing.T) {
	type args struct {
		result   *search.ElasticResult
		pageSize int
	}
	tests := []struct {
		name string
		args args
		want *UserPaged
	}{
		{
			name: "when elasticResult is null should return empty userPaged",
			args: args{
				result:   nil,
				pageSize: 15,
			},
			want: &UserPaged{},
		},
		{
			name: "when elasticResult has two sources should return two users on userPaged",
			args: args{
				result: buildElasticResult([]User{
					{string(rune(rand.Int())), "Estefane Roxinol", "estefane.roxinol"},
					{string(rune(rand.Int())), "stefane Aureniza", "estefaneaureniza"},
				}, ""),
				pageSize: 15,
			},
			want: &UserPaged{
				Users: []User{
					{string(rune(rand.Int())), "Estefane Roxinol", "estefane.roxinol"},
					{string(rune(rand.Int())), "stefane Aureniza", "estefaneaureniza"},
				},
			},
		},
		{
			name: "when elasticResult has two sources should return two users on userPaged with scroll",
			args: args{
				result: buildElasticResult([]User{
					{string(rune(rand.Int())), "Estefane Roxinol", "estefane.roxinol"},
					{string(rune(rand.Int())), "stefane Aureniza", "estefaneaureniza"},
				}, "FGluY2x1ZGVfY29udGV4dF91dWlkDXF1ZXJ5QW5kRmV0Y2gBFkJySEhCMnR5U1ZpQS16STBpRmFHMXcAAAAAAAAAMhY5QW81aWx1clE4dWxnMzlFZDI2NEln"),
				pageSize: 15,
			},
			want: &UserPaged{
				ScrollId: "FGluY2x1ZGVfY29udGV4dF91dWlkDXF1ZXJ5QW5kRmV0Y2gBFkJySEhCMnR5U1ZpQS16STBpRmFHMXcAAAAAAAAAMhY5QW81aWx1clE4dWxnMzlFZDI2NEln",
				Users: []User{
					{string(rune(rand.Int())), "Estefane Roxinol", "estefane.roxinol"},
					{string(rune(rand.Int())), "stefane Aureniza", "estefaneaureniza"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toPaged(tt.args.result); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("toPaged() = %v, want %v", got, tt.want)
			}
		})
	}
}

func buildElasticResult(users []User, scrollId string) *search.ElasticResult {
	hits := make([]search.SourceHits, 0, len(users))

	for _, u := range users {
		hits = append(hits, search.SourceHits{
			Source: map[string]interface{}{
				"ID":       u.ID,
				"Name":     u.Name,
				"Username": u.Username},
		})
	}
	return &search.ElasticResult{
		ScrollId: scrollId,
		RootHits: search.RootHits{
			Hits: hits,
		},
	}
}
