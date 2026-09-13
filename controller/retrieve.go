package controller

import (
	"MyMemory/logs"
	"MyMemory/model"
	"context"
	"sort"
)

func (c *Controller) Retrieve(ctx context.Context, req *model.Request) ([]*model.SearchResult, error) {
	search := &model.SearchRequest{
		UserID:   req.UID,
		Content:  req.Content,
		TopK:     10,
		KnnScore: req.KnnScore,
	}
	content, err := c.esRepo.BM25Search(ctx, search)
	if err != nil {
		logs.Error("get BM25 search err: %v", err)
		return nil, err
	}
	vec, err := c.esRepo.KNNSearch(ctx, search)
	if err != nil {
		logs.Error("get KNN search err: %v", err)
		return nil, err
	}
	return c.RrfSort(content, vec, req.TopK), nil
}

type rrf struct {
	r1     int
	r2     int
	score  float64
	doc    *model.Document
	source string
}

const rrfScore = 60

func (r *rrf) calScore() {
	score := 0.0
	if r.r1 < 1<<30 {
		score += 1.0 / float64(r.r1+rrfScore)
	}
	if r.r2 < 1<<30 {
		score += 1.0 / float64(r.r2+rrfScore)
	}
	r.score = score
}

func (c *Controller) RrfSort(contents, vecs []*model.SearchResult, top int) []*model.SearchResult {
	results := getResults(contents, vecs)
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	var docs []*model.SearchResult
	if top > len(results) {
		top = len(results)
	}
	for i := 0; i < top; i++ {
		docs = append(docs, &model.SearchResult{
			Score:  results[i].score,
			Doc:    results[i].doc,
			Source: results[i].source,
		})
	}
	return docs
}

func getResults(contents, vecs []*model.SearchResult) []*rrf {
	docMap := make(map[string]*rrf)
	for i, doc := range contents {
		docMap[doc.Doc.ID] = &rrf{
			r1:     i + 1,
			r2:     1 << 30, // 默认极大
			doc:    doc.Doc,
			source: doc.Source,
		}
	}

	var results []*rrf
	for i, vec := range vecs {
		doc, ok := docMap[vec.Doc.ID]
		if !ok {
			doc = &rrf{
				r1:     1 << 30, // 默认极大
				r2:     i + 1,
				doc:    vec.Doc,
				source: vec.Source,
			}
			docMap[vec.Doc.ID] = doc
		} else {
			doc.r2 = i + 1
			doc.source = doc.source + "/" + vec.Source
		}
	}

	for _, doc := range docMap {
		doc.calScore()
		results = append(results, doc)
	}
	return results
}
