package model

import "encoding/json"

type Document struct {
	ID      string    `json:"id"`
	UserID  int64     `json:"user_id"`
	Content string    `json:"content"`
	Vec     []float32 `json:"vec"`
	Title   string    `json:"title"`
}

type EsType string

const (
	keyword EsType = "keyword"
	text    EsType = "text"
	vector  EsType = "dense_vector"
	long    EsType = "long"
)

type Properties struct {
	UserID  *Property `json:"user_id"`
	ID      *Property `json:"id"`
	Title   *Property `json:"title"`
	Content *Property `json:"content"`
	Vec     *Vector   `json:"vec"`
}

type Vector struct {
	T       EsType `json:"type"`
	Dims    int    `json:"dims"`       // 维度
	Index   bool   `json:"index"`      // 构建 HNSW 索引，支持 knn 快速近似检索（生产必选）
	Similar string `json:"similarity"` // 相似度算法
}

type Property struct {
	T EsType `json:"type"`
}

func (p *Properties) Mappings() string {
	m := map[string]any{
		"mappings": map[string]any{
			"properties": p,
		},
	}
	v, _ := json.Marshal(m)
	return string(v)
}

type SearchRequest struct {
	UserID   int64
	Content  string
	TopK     int
	KnnScore int
}

type SearchResult struct {
	Doc    *Document
	Score  float64
	Source string
}

type Shard struct {
	Total      int `json:"total"`
	Successful int `json:"successful"`
	Skipped    int `json:"skipped"`
	Failed     int `json:"failed"`
}

type Total struct {
	Value    int    `json:"value"`
	Relation string `json:"relation"`
}

// EsHit 单条文档hit
type EsHit struct {
	Index  string    `json:"_index"`
	ID     string    `json:"_id"`
	Score  float64   `json:"_score"`
	Source *Document `json:"_source"` // 业务文档内容
}

type Hit struct {
	Total    *Total   `json:"total"`
	MaxScore float64  `json:"max_score"`
	Hits     []*EsHit `json:"hits"`
}

// EsSearchResp ES搜索完整返回体
type EsSearchResp struct {
	Took     int    `json:"took"`
	TimedOut bool   `json:"timed_out"`
	Shards   *Shard `json:"_shards"`
	Hits     *Hit   `json:"hits"`
}
