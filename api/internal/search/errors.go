package search

import "github.com/pkg/errors"

var ErrSearchServiceUnavailable = errors.New("Search Service Unavailable")

type InvalidQuery struct {
	message      string
	error        error
	elasticError *ElasticError
}

func (iq InvalidQuery) Error() string {
	return iq.message
}

func (iq InvalidQuery) OriginError() error {
	return iq.error
}

func (iq InvalidQuery) ElasticError() *ElasticError {
	return iq.elasticError
}
