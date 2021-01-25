package search

type ElasticResult struct {
	ScrollId string `json:"_scroll_id"`
	Took     int    `json:"took"`
	RootHits `json:"hits"`
}

type RootHits struct {
	Total    Total        `json:"total"`
	MaxScore float64      `json:"max_score"`
	Hits     []SourceHits `json:"hits"`
}

type Total struct {
	Value int `json:"value"`
}

type SourceHits struct {
	Index  string                 `json:"_index"`
	Type   string                 `json:"_type"`
	ID     string                 `json:"_id"`
	Score  float64                `json:"_score"`
	Source map[string]interface{} `json:"_source"`
}

type ElasticError struct {
	Error struct {
		RootCause []struct {
			Type   string `json:"type"`
			Reason string `json:"reason"`
		} `json:"root_cause"`
		Type     string `json:"type"`
		Reason   string `json:"reason"`
		CausedBy struct {
			Type   string `json:"type"`
			Reason string `json:"reason"`
		} `json:"caused_by"`
	} `json:"error"`
	Status int `json:"status"`
}
