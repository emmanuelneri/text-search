package users

import "api/internal/search"

func toPaged(result *search.ElasticResult) *UserPaged {
	if result == nil {
		return &UserPaged{}
	}

	hits := result.Hits
	users := make([]User, 0, len(hits))
	for _, hit := range hits {
		users = append(users, User{
			ID:       hit.Source["ID"].(string),
			Name:     hit.Source["Name"].(string),
			Username: hit.Source["Username"].(string),
		})
	}

	return &UserPaged{
		ScrollId: result.ScrollId,
		Users:    users,
	}
}
