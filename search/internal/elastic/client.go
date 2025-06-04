package elastic

import (
	"bytes"
	"encoding/json"
	"log"
	"search/internal/models"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
)

var client *elasticsearch.Client

func Init() {
	cfg := elasticsearch.Config{
		Addresses: []string{"http://localhost:9200"},
	}
	c, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Elasticsearch init error: %v", err)
	}
	client = c
}

func GetClient() *elasticsearch.Client {
	return client
}

func BuildSearchQuery(req models.SearchRequest) *bytes.Reader {
	q := map[string]interface{}{
		"from": (req.Page - 1) * req.PageSize,
		"size": req.PageSize,
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"must":   []interface{}{},
				"filter": []interface{}{},
			},
		},
	}

	boolQuery := q["query"].(map[string]interface{})["bool"].(map[string]interface{})

	if req.Query != "" {
		boolQuery["must"] = append(boolQuery["must"].([]interface{}), map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  req.Query,
				"fields": []string{"name", "description"},
			},
		})
	}

	rangeFilter := map[string]interface{}{}
	if req.MinPrice > 0 {
		rangeFilter["gte"] = req.MinPrice
	}
	if req.MaxPrice > 0 {
		rangeFilter["lte"] = req.MaxPrice
	}
	if len(rangeFilter) > 0 {
		boolQuery["filter"] = append(boolQuery["filter"].([]interface{}), map[string]interface{}{
			"range": map[string]interface{}{
				"price": rangeFilter,
			},
		})
	}

	if req.SortBy != "" {
		q["sort"] = []interface{}{map[string]interface{}{req.SortBy: map[string]interface{}{"order": "asc"}}}
	}

	body, _ := json.Marshal(q)
	return bytes.NewReader(body)
}

func ParseSearchResult(res *esapi.Response, req models.SearchRequest) (*models.SearchResponse, error) {
	var r struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source models.Good `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	products := make([]*models.Good, 0, len(r.Hits.Hits))
	for _, hit := range r.Hits.Hits {
		p := hit.Source
		products = append(products, &p)
	}

	return &models.SearchResponse{
		Products: products,
		Total:    r.Hits.Total.Value,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
