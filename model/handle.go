package model

type Request struct {
	ID       string `json:"id"`
	UID      int64  `json:"user_id"`
	Content  string `json:"content"`
	Title    string `json:"title"`
	Chunk    int    `json:"chunk"`
	KnnScore int    `json:"knn_score"`
	TopK     int    `json:"top_k"`
}
