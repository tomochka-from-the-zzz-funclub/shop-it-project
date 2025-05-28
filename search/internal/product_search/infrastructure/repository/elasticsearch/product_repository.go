package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	entityF "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/facet/entity"
	entityFB "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/facet_bucket/entity"
	entitySP "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/search_params/entity"
	entitySR "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/search_result/entity"

	"github.com/elastic/go-elasticsearch/v8/esapi"
	entityCH "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/cluster_health/entity"
	entityP "gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/domain/product/entity"
	"gitlab.mai.ru/4-bogatyra/backend/search/internal/product_search/infrastructure/config/elasticsearch"
)

type ProductRepository struct {
	es    *elasticsearch.Client
	index string
}

func NewProductRepository(client *elasticsearch.Client, index string) *ProductRepository {
	log.Printf("[ProductRepository] Initialized with index=%q", index)
	return &ProductRepository{
		es:    client,
		index: index,
	}
}

func (r *ProductRepository) Save(p *entityP.Product) error {
	log.Printf("[ProductRepository] Saving product ID=%s", p.ID)

	body, err := json.Marshal(p)
	if err != nil {
		log.Printf("[ProductRepository][ERROR] Failed to marshal product ID=%s: %v", p.ID, err)
		return fmt.Errorf("marshal product: %w", err)
	}

	req := esapi.IndexRequest{
		Index:      r.index,
		DocumentID: p.ID,
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}
	res, err := req.Do(context.Background(), r.es.Connection)
	if err != nil {
		log.Printf("[ProductRepository][ERROR] Index request failed for ID=%s: %v", p.ID, err)
		return fmt.Errorf("index request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[ProductRepository][ERROR] Index returned error for ID=%s: %s", p.ID, res.String())
		return fmt.Errorf("index error: %s", res.String())
	}

	log.Printf("[ProductRepository] Product indexed successfully ID=%s", p.ID)
	return nil
}

func (r *ProductRepository) BulkSave(products []*entityP.Product) error {
	log.Printf("[ProductRepository] Bulk saving %d products", len(products))

	var buf bytes.Buffer
	for _, p := range products {
		meta := map[string]map[string]string{"index": {"_index": r.index, "_id": p.ID}}
		metaBytes, _ := json.Marshal(meta)
		buf.Write(metaBytes)
		buf.WriteByte('\n')

		srcBytes, _ := json.Marshal(p)
		buf.Write(srcBytes)
		buf.WriteByte('\n')
	}

	res, err := r.es.Connection.Bulk(
		bytes.NewReader(buf.Bytes()),
		r.es.Connection.Bulk.WithRefresh("true"),
	)
	if err != nil {
		log.Printf("[ProductRepository][ERROR] Bulk request failed: %v", err)
		return fmt.Errorf("bulk request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[ProductRepository][ERROR] Bulk returned error: %s", res.String())
		return fmt.Errorf("bulk error: %s", res.String())
	}

	log.Printf("[ProductRepository] Bulk save succeeded for %d products", len(products))
	return nil
}

func (r *ProductRepository) Delete(id string) error {
	log.Printf("[ProductRepository] Deleting product ID=%s", id)

	res, err := r.es.Connection.Delete(
		r.index,
		id,
		r.es.Connection.Delete.WithRefresh("true"),
	)
	if err != nil {
		log.Printf("[ProductRepository][ERROR] Delete request failed for ID=%s: %v", id, err)
		return fmt.Errorf("delete request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[ProductRepository][ERROR] Delete returned error for ID=%s: %s", id, res.String())
		return fmt.Errorf("delete error: %s", res.String())
	}

	log.Printf("[ProductRepository] Product deleted successfully ID=%s", id)
	return nil
}

//func (r *ProductRepository) Search(params *entitySP.SearchParams) (*entitySR.SearchResult, error) {
//	log.Printf("[ProductRepository] Searching products with params: %+v", params)
//
//	query := make(map[string]interface{})
//	boolQ := make(map[string]interface{})
//
//	if params.Query != "" {
//		boolQ["must"] = map[string]interface{}{
//			"multi_match": map[string]interface{}{
//				"query":  params.Query,
//				"fields": []string{"name^2", "description"},
//			},
//		}
//	}
//
//	var filters []map[string]interface{}
//	if len(params.Categories) > 0 {
//		filters = append(filters, map[string]interface{}{"terms": map[string]interface{}{"category": params.Categories}})
//	}
//	if len(params.Brand) > 0 {
//		filters = append(filters, map[string]interface{}{"terms": map[string]interface{}{"brand": params.Brand}})
//	}
//	if params.MinPrice > 0 || params.MaxPrice > 0 {
//		rangeQ := map[string]interface{}{}
//		if params.MinPrice > 0 {
//			rangeQ["gte"] = params.MinPrice
//		}
//		if params.MaxPrice > 0 {
//			rangeQ["lte"] = params.MaxPrice
//		}
//		filters = append(filters, map[string]interface{}{"range": map[string]interface{}{"price": rangeQ}})
//	}
//	if len(filters) > 0 {
//		boolQ["filter"] = filters
//	}
//
//	query["query"] = map[string]interface{}{"bool": boolQ}
//
//	if params.SortBy != "" {
//		order := params.SortOrder
//		if order != "asc" && order != "desc" {
//			order = "desc"
//		}
//		var sortField string
//		switch params.SortBy {
//		case "price":
//			sortField = "price"
//		case "popularity":
//			sortField = "popularity"
//		default:
//			sortField = "_score"
//		}
//		query["sort"] = []map[string]interface{}{
//			{sortField: map[string]interface{}{"order": order}},
//		}
//	}
//
//	from := (params.Page - 1) * params.PageSize
//	query["from"] = from
//	query["size"] = params.PageSize
//
//	if len(params.HighlightFields) > 0 {
//		hl := map[string]interface{}{"fields": map[string]interface{}{}}
//		for _, f := range params.HighlightFields {
//			hl["fields"].(map[string]interface{})[f] = map[string]interface{}{}
//		}
//		query["highlight"] = hl
//	}
//
//	aggs := map[string]interface{}{}
//	for _, field := range []string{"category", "brand"} {
//		aggs[field] = map[string]interface{}{
//			"terms": map[string]interface{}{"field": field, "size": 10},
//		}
//	}
//	query["aggs"] = aggs
//
//	body, _ := json.Marshal(query)
//	res, err := r.es.Connection.Search(
//		r.es.Connection.Search.WithContext(context.Background()),
//		r.es.Connection.Search.WithIndex(r.index),
//		r.es.Connection.Search.WithBody(bytes.NewReader(body)),
//	)
//	if err != nil {
//		log.Printf("[ProductRepository][ERROR] Search request failed: %v", err)
//		return nil, fmt.Errorf("search request: %w", err)
//	}
//	defer res.Body.Close()
//
//	if res.IsError() {
//		log.Printf("[ProductRepository][ERROR] Search returned error: %s", res.String())
//		return nil, fmt.Errorf("search error: %s", res.String())
//	}
//
//	var resp struct {
//		Hits struct {
//			Total struct {
//				Value int64 `json:"value"`
//			} `json:"total"`
//			Hits []struct {
//				ID        string              `json:"_id"`
//				Source    entityP.Product     `json:"_source"`
//				Highlight map[string][]string `json:"highlight"`
//			} `json:"hits"`
//		} `json:"hits"`
//		Aggregations map[string]struct {
//			Buckets []struct {
//				Key      interface{} `json:"key"`
//				DocCount int64       `json:"doc_count"`
//			} `json:"buckets"`
//		} `json:"aggregations"`
//	}
//	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
//		log.Printf("[ProductRepository][ERROR] Failed to decode search response: %v", err)
//		return nil, fmt.Errorf("decode response: %w", err)
//	}
//
//	result := &entitySR.SearchResult{
//		Total:      resp.Hits.Total.Value,
//		Page:       params.Page,
//		PageSize:   params.PageSize,
//		Products:   make([]*entityP.Product, 0, len(resp.Hits.Hits)),
//		Highlights: make(map[string][]string),
//		Facets:     make(map[string]entityF.Facet),
//	}
//
//	for _, hit := range resp.Hits.Hits {
//		prod := hit.Source
//		prod.ID = hit.ID
//		result.Products = append(result.Products, &prod)
//	}
//
//	for field, agg := range resp.Aggregations {
//		facet := entityF.Facet{
//			Field:   field,
//			Buckets: make([]entityFB.FacetBucket, len(agg.Buckets)),
//		}
//		for i, b := range agg.Buckets {
//			facet.Buckets[i] = entityFB.FacetBucket{
//				Key:   fmt.Sprint(b.Key),
//				Count: b.DocCount,
//			}
//		}
//		result.Facets[field] = facet
//	}
//
//	log.Printf("[ProductRepository] Search completed: found %d products", result.Total)
//	return result, nil
//}

func (r *ProductRepository) Search(params *entitySP.SearchParams) (*entitySR.SearchResult, error) {
	log.Printf("[ProductRepository] Starting product search with params: %+v", params)

	query := make(map[string]interface{})
	boolQ := make(map[string]interface{})

	if params.Query != "" {
		boolQ["must"] = map[string]interface{}{
			"multi_match": map[string]interface{}{
				"query":  params.Query,
				"fields": []string{"Name^2", "Description", "Category", "Brand"},
			},
		}
		log.Printf("[ProductRepository] Applied full-text query: %s", params.Query)
	}

	var filters []map[string]interface{}
	if len(params.Categories) > 0 {
		log.Printf("[ProductRepository] Filtering by categories: %+v", params.Categories)
		filters = append(filters, map[string]interface{}{
			"terms": map[string]interface{}{"Category": params.Categories},
		})
	}
	if len(params.Brand) > 0 {
		log.Printf("[ProductRepository] Filtering by brands: %+v", params.Brand)
		filters = append(filters, map[string]interface{}{
			"terms": map[string]interface{}{"Brand": params.Brand},
		})
	}
	if params.MinPrice > 0 || params.MaxPrice > 0 {
		rangeQ := map[string]interface{}{}
		if params.MinPrice > 0 {
			rangeQ["gte"] = params.MinPrice
		}
		if params.MaxPrice > 0 {
			rangeQ["lte"] = params.MaxPrice
		}
		log.Printf("[ProductRepository] Applying price range filter: %+v", rangeQ)
		filters = append(filters, map[string]interface{}{
			"range": map[string]interface{}{"Price": rangeQ},
		})
	}
	if len(filters) > 0 {
		boolQ["filter"] = filters
		log.Printf("[ProductRepository] Final filters: %+v", filters)
	}

	query["query"] = map[string]interface{}{"bool": boolQ}

	if params.SortBy != "" {
		order := params.SortOrder
		if order != "asc" && order != "desc" {
			order = "desc"
		}
		var sortField string
		switch params.SortBy {
		case "price":
			sortField = "Price"
		case "popularity":
			sortField = "Popularity"
		default:
			sortField = "_score"
		}
		query["sort"] = []map[string]interface{}{
			{sortField: map[string]interface{}{"order": order}},
		}
		log.Printf("[ProductRepository] Sorting by %s %s", sortField, order)
	}

	from := (params.Page - 1) * params.PageSize
	query["from"] = from
	query["size"] = params.PageSize
	log.Printf("[ProductRepository] Pagination - from: %d, size: %d", from, params.PageSize)

	if len(params.HighlightFields) > 0 {
		hl := map[string]interface{}{"fields": map[string]interface{}{}}
		for _, f := range params.HighlightFields {
			hl["fields"].(map[string]interface{})[f] = map[string]interface{}{}
		}
		query["highlight"] = hl
		log.Printf("[ProductRepository] Highlighting fields: %+v", params.HighlightFields)
	}

	aggs := map[string]interface{}{}
	for _, field := range []string{"Category", "Brand"} {
		aggs[field] = map[string]interface{}{
			"terms": map[string]interface{}{"field": field, "size": 10},
		}
	}
	query["aggs"] = aggs
	log.Printf("[ProductRepository] Aggregations configured for fields: category, brand")

	body, _ := json.Marshal(query)
	log.Printf("[ProductRepository] Final ES query JSON: %s", string(body))
	res, err := r.es.Connection.Search(
		r.es.Connection.Search.WithContext(context.Background()),
		r.es.Connection.Search.WithIndex(r.index),
		r.es.Connection.Search.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		log.Printf("[ProductRepository][ERROR] Search request failed: %v", err)
		return nil, fmt.Errorf("search request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[ProductRepository][ERROR] Search returned error: %s", res.String())
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var resp struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				ID        string              `json:"_id"`
				Source    entityP.Product     `json:"_source"`
				Highlight map[string][]string `json:"highlight"`
			} `json:"hits"`
		} `json:"hits"`
		Aggregations map[string]struct {
			Buckets []struct {
				Key      interface{} `json:"key"`
				DocCount int64       `json:"doc_count"`
			} `json:"buckets"`
		} `json:"aggregations"`
	}
	if err = json.NewDecoder(res.Body).Decode(&resp); err != nil {
		log.Printf("[ProductRepository][ERROR] Failed to decode search response: %v", err)
		return nil, fmt.Errorf("decode response: %w", err)
	}

	result := &entitySR.SearchResult{
		Total:      resp.Hits.Total.Value,
		Page:       params.Page,
		PageSize:   params.PageSize,
		Products:   make([]*entityP.Product, 0, len(resp.Hits.Hits)),
		Highlights: make(map[string][]string),
		Facets:     make(map[string]entityF.Facet),
	}

	log.Printf("[ProductRepository] Search returned %d hits", len(resp.Hits.Hits))

	for _, hit := range resp.Hits.Hits {
		prod := hit.Source
		prod.ID = hit.ID
		result.Products = append(result.Products, &prod)
		log.Printf("[ProductRepository] Found product: ID=%s, Name=%s", prod.ID, prod.Name)
		if len(hit.Highlight) > 0 {
			log.Printf("[ProductRepository] Highlight for product ID=%s: %+v", prod.ID, hit.Highlight)
			result.Highlights[prod.ID] = []string{}
			for _, snippets := range hit.Highlight {
				result.Highlights[prod.ID] = append(result.Highlights[prod.ID], snippets...)
			}
		}
	}

	for field, agg := range resp.Aggregations {
		facet := entityF.Facet{
			Field:   field,
			Buckets: make([]entityFB.FacetBucket, len(agg.Buckets)),
		}
		for i, b := range agg.Buckets {
			facet.Buckets[i] = entityFB.FacetBucket{
				Key:   fmt.Sprint(b.Key),
				Count: b.DocCount,
			}
		}
		result.Facets[field] = facet
		log.Printf("[ProductRepository] Facet %s: %+v", field, facet.Buckets)
	}

	log.Printf("[ProductRepository] Search completed successfully. Total products found: %d", result.Total)
	return result, nil
}

func (r *ProductRepository) Health() (*entityCH.ClusterHealth, error) {
	log.Println("[ProductRepository] Checking cluster health")

	res, err := r.es.Connection.Cluster.Health(
		r.es.Connection.Cluster.Health.WithTimeout(5 * time.Second),
	)
	if err != nil {
		log.Printf("[ProductRepository][ERROR] Cluster health request failed: %v", err)
		return nil, fmt.Errorf("cluster health request: %w", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		log.Printf("[ProductRepository][ERROR] Cluster health returned error: %s", res.String())
		return nil, fmt.Errorf("cluster health error: %s", res.String())
	}

	var ch struct {
		Status           string `json:"status"`
		ActiveShards     int64  `json:"active_shards"`
		RelocatingShards int64  `json:"relocating_shards"`
		UnassignedShards int64  `json:"unassigned_shards"`
		TimedOut         bool   `json:"timed_out"`
	}
	if err = json.NewDecoder(res.Body).Decode(&ch); err != nil {
		log.Printf("[ProductRepository][ERROR] Failed to decode cluster health response: %v", err)
		return nil, fmt.Errorf("decode cluster health: %w", err)
	}

	clusterHealth := &entityCH.ClusterHealth{
		Status:           ch.Status,
		ActiveShards:     ch.ActiveShards,
		RelocatingShards: ch.RelocatingShards,
		UnassignedShards: ch.UnassignedShards,
		TimedOut:         ch.TimedOut,
	}
	log.Printf("[ProductRepository] Cluster health status: %s", clusterHealth.Status)
	return clusterHealth, nil
}

func (r *ProductRepository) IndicesExistsAndCreateIfMissing() error {
	ctx := context.Background()

	res, err := r.es.Connection.Indices.Exists([]string{r.index})
	if err != nil {
		return fmt.Errorf("check index existence: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {
		log.Printf("[ProductRepository] Index '%s' already exists", r.index)
		return nil
	}

	if res.StatusCode != 404 {
		return fmt.Errorf("unexpected status when checking index existence: %d", res.StatusCode)
	}

	log.Printf("[ProductRepository] Index '%s' not found. Creating...", r.index)

	mapping := map[string]interface{}{
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"Name": map[string]interface{}{
					"type": "text",
					"fields": map[string]interface{}{
						"keyword": map[string]interface{}{
							"type": "keyword",
						},
					},
				},
				"Description": map[string]interface{}{
					"type": "text",
				},
				"Category": map[string]interface{}{
					"type": "keyword",
				},
				"Brand": map[string]interface{}{
					"type": "keyword",
				},
				"Price": map[string]interface{}{
					"type": "float",
				},
				"Popularity": map[string]interface{}{
					"type": "integer",
				},
			},
		},
	}

	body, err := json.Marshal(mapping)
	if err != nil {
		return fmt.Errorf("marshal mapping: %w", err)
	}

	// Создание индекса
	createRes, err := r.es.Connection.Indices.Create(
		r.index,
		r.es.Connection.Indices.Create.WithContext(ctx),
		r.es.Connection.Indices.Create.WithBody(bytes.NewReader(body)),
	)
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}
	defer createRes.Body.Close()

	if createRes.IsError() {
		buf := new(bytes.Buffer)
		buf.ReadFrom(createRes.Body)
		return fmt.Errorf("failed to create index: %s", strings.TrimSpace(buf.String()))
	}

	log.Printf("[ProductRepository] Index '%s' created successfully", r.index)
	return nil
}
