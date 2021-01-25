package users

import (
	"api/internal/search"
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"github.com/pkg/errors"
	"time"
)

type Service interface {
	Search(ctx context.Context, query string) (*UserPaged, error)
	Scroll(id string) (*UserPaged, error)
}

type ServiceImpl struct {
	elasticsearchClient *elasticsearch.Client
	index               string
	pageSize            int
	scrollDuration      time.Duration
}

func NewUserService(elasticsearchClient *elasticsearch.Client, index string, pageSize int, scrollDuration time.Duration) Service {
	return &ServiceImpl{
		elasticsearchClient: elasticsearchClient,
		index:               index,
		pageSize:            pageSize,
		scrollDuration:      scrollDuration,
	}
}

func (s ServiceImpl) Search(ctx context.Context, query string) (*UserPaged, error) {
	requestQuery := search.RequestQuery{
		Value:            query,
		MultiMatchType:   "phrase_prefix",
		MultiMatchFields: []string{"Name", "Username"},
		FieldValueFactor: "relevance",
		Boost:            "5",
	}.Query()

	var searchRequest bytes.Buffer
	if err := json.NewEncoder(&searchRequest).Encode(requestQuery); err != nil {
		return nil, errors.Wrap(err, "fail to encode request")
	}

	res, err := s.elasticsearchClient.Search(
		s.elasticsearchClient.Search.WithContext(ctx),
		s.elasticsearchClient.Search.WithIndex(s.index),
		s.elasticsearchClient.Search.WithBody(&searchRequest),
		s.elasticsearchClient.Search.WithSize(s.pageSize),
		s.elasticsearchClient.Search.WithScroll(s.scrollDuration),
	)

	return handleResponse(err, res)
}

func (s ServiceImpl) Scroll(id string) (*UserPaged, error) {
	res, err := s.elasticsearchClient.Scroll(
		s.elasticsearchClient.Scroll.WithScrollID(id),
		s.elasticsearchClient.Scroll.WithScroll(s.scrollDuration),
	)

	if err != nil {
		return nil, errors.Wrap(err, "fail to search")
	}

	return handleResponse(err, res)
}

func handleResponse(err error, res *esapi.Response) (*UserPaged, error) {
	if err != nil {
		return nil, errors.Wrap(search.ErrSearchServiceUnavailable, err.Error())
	}

	defer res.Body.Close()
	elasticResult, err := search.HandleResponse(res, err)
	if err != nil {
		return nil, err
	}

	usersPaged := toPaged(elasticResult)
	return usersPaged, nil
}
