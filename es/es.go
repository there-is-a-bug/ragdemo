package es

import (
	"MyMemory/embed"
	"MyMemory/logs"
	"MyMemory/model"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sync"

	es9 "github.com/elastic/go-elasticsearch/v9"
)

type ESRepo struct {
	es    *es9.Client
	index string
}

type filter struct {
	knnScore int
	topK     int
	source   string
}

var (
	so   sync.Once
	Repo *ESRepo
)

func init() {
	so.Do(func() {
		cli, err := es9.New(es9.WithAddresses(model.EsAddress))
		if err != nil {
			panic(err)
		}
		Repo = &ESRepo{
			es:    cli,
			index: model.EsIndex,
		}
	})
}

func GetESRepo() *ESRepo {
	return Repo
}

func (e *ESRepo) Insert(ctx context.Context, doc *model.Document) error {
	data, err := json.Marshal(doc)
	if err != nil {
		logs.Error("error marshalling document: %v", err)
		return err
	}
	resp, err := e.es.Index(e.index, bytes.NewReader(data), e.es.Index.WithDocumentID(doc.ID))
	if err != nil {
		logs.Error("insert es fail", err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logs.Error("error closing body: %v", err)
		}
		var by []byte
		resp.Body.Read(by)
		fmt.Println(string(by))
	}(resp.Body)

	return nil
}

func (e *ESRepo) BM25Search(ctx context.Context, request *model.SearchRequest) ([]*model.SearchResult, error) {
	body := map[string]any{
		"size": request.TopK * 2,
		"query": map[string]any{
			"bool": map[string]any{
				"must": []any{
					map[string]any{
						"match": map[string]any{
							"content": request.Content,
						},
					},
				},
				"filter": []any{
					map[string]any{
						"term": map[string]any{
							"user_id": request.UserID,
						},
					},
				},
			},
		},
	}
	return e.search(ctx, body, &filter{
		source: "bm25",
	})
}

func (e *ESRepo) KNNSearch(ctx context.Context, request *model.SearchRequest) ([]*model.SearchResult, error) {
	vec, _ := embed.Embedding(request.Content)
	body := map[string]any{
		"size": request.TopK,
		"knn": map[string]any{
			"field":          "vec",
			"query_vector":   vec,
			"k":              request.TopK * 2,
			"num_candidates": 50,
			"filter": map[string]any{
				"term": map[string]any{
					"user_id": request.UserID,
				},
			},
		},
	}
	return e.search(ctx, body, &filter{
		knnScore: request.KnnScore,
		source:   "knn",
	})
}

func (e *ESRepo) search(ctx context.Context, body map[string]any, f *filter) ([]*model.SearchResult, error) {
	by, _ := json.Marshal(body)
	es := e.es
	resp, err := es.Search(
		es.Search.WithIndex(e.index),
		es.Search.WithBody(bytes.NewReader(by)),
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取返回原始json
	var reader bytes.Buffer
	_, _ = reader.ReadFrom(resp.Body)
	var result model.EsSearchResp
	_ = json.Unmarshal(reader.Bytes(), &result)
	if result.Hits == nil {
		return nil, fmt.Errorf("none result")
	}
	fmt.Println(f.source, "result:", string(reader.Bytes()))

	var results []*model.SearchResult
	for _, hit := range result.Hits.Hits {
		if f.source == "knn" && int(hit.Score*100) < f.knnScore {
			continue
		}
		results = append(results, &model.SearchResult{
			Doc:    hit.Source,
			Score:  hit.Score,
			Source: f.source,
		})
	}
	return results, nil
}
