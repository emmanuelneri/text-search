package search

type RequestQuery struct {
	Value            string
	MultiMatchType   string
	MultiMatchFields []string
	FieldValueFactor string
	Boost            string
}

func (s RequestQuery) Query() interface{} {
	fieldValueFactor := Function{FieldValueFactor{Field: s.FieldValueFactor}}
	if s.Value == "" {
		return s.matchAllQuery(fieldValueFactor)
	}

	return s.MultiMatchQuery(fieldValueFactor)
}

func (s RequestQuery) MultiMatchQuery(fieldValueFactor Function) SearchRequestQuery {
	return SearchRequestQuery{
		RootQuery{FunctionScore{
			Query: Query{
				MultiMatch: &MultiMatch{
					Query:  s.Value,
					Type:   s.MultiMatchType,
					Fields: s.MultiMatchFields,
				},
			},
			Boost:     s.Boost,
			Functions: []Function{fieldValueFactor},
		}},
	}
}

func (s RequestQuery) matchAllQuery(fieldValueFactor Function) SearchRequestQuery {
	return SearchRequestQuery{
		RootQuery{FunctionScore{
			Query: Query{
				MatchAll: &MatchAll{},
			},
			Boost:     s.Boost,
			Functions: []Function{fieldValueFactor},
		}},
	}
}

type SearchRequestQuery struct {
	RootQuery `json:"query"`
}

type RootQuery struct {
	FunctionScore `json:"function_score"`
}

type FunctionScore struct {
	Query     `json:"query"`
	Boost     string     `json:"boost"`
	Functions []Function `json:"functions"`
}

type Query struct {
	*MatchAll   `json:"match_all,omitempty"`
	*MultiMatch `json:"multi_match,omitempty"`
}

type MatchAll struct {
}

type MultiMatch struct {
	Query  string   `json:"query"`
	Type   string   `json:"type"`
	Fields []string `json:"fields"`
}

type Function struct {
	FieldValueFactor `json:"field_value_factor"`
}

type FieldValueFactor struct {
	Field string `json:"field"`
}
