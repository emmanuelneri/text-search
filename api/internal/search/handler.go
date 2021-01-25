package search

import (
	"api/internal/logs"
	"encoding/json"
	"fmt"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"go.uber.org/zap"
)

func HandleResponse(res *esapi.Response, err error) (*ElasticResult, error) {
	if res.IsError() {
		elasticError := &ElasticError{}
		if err := json.NewDecoder(res.Body).Decode(&elasticError); err != nil {
			return nil, &InvalidQuery{message: fmt.Sprintf("%s response error. fail decode response body", res.Status()), error: err}
		}

		if elasticError.Status == 0 {
			return nil, &InvalidQuery{message: fmt.Sprintf("%s response error. empty response body", res.Status()), error: err}
		}

		message := fmt.Sprintf("status %s - type :%s - reason: %s ", res.Status(),
			elasticError.Error.Type, elasticError.Error.Reason)

		return nil, &InvalidQuery{message: message, error: err, elasticError: elasticError}
	}

	result := &ElasticResult{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, &InvalidQuery{message: "decode response body error", error: err}
	}

	if res.HasWarnings() {
		logs.Logger.Info("Warnings...", zap.Strings("warnings", res.Warnings()))
	}

	return result, nil
}
