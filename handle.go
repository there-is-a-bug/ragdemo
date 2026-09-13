package main

import (
	ctl "MyMemory/controller"
	"MyMemory/logs"
	"MyMemory/model"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func getReq(c *gin.Context) (*model.Request, error) {
	req := &model.Request{}
	err := c.BindJSON(req)
	if err != nil {
		return nil, err
	}
	if req.UID == 0 {
		return nil, fmt.Errorf("user_id should not be empty")
	}
	if req.ID == "" {
		req.ID = time.Now().Format("20060102150405")
	}
	if req.Title == "" {
		req.Title = "doc" + time.Now().Format("20060102150405")
	}
	content := req.Content
	if len(content) == 0 {
		return nil, fmt.Errorf("content should not be empty")
	}
	chunk := req.Chunk
	if chunk > 4096 || chunk <= 0 {
		logs.Warn("chunk should > 0 or less", chunk)
		req.Chunk = 4096
	}
	if req.KnnScore >= 100 {
		logs.Warn("knn_score should less 100", req.KnnScore)
		req.KnnScore = 80
	}
	logs.Info("Get request: %+v", req)
	return req, nil

}

func Insert(c *gin.Context) {
	req, err := getReq(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	ctx := context.Background()
	err = ctl.GetController().Insert(ctx, req)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	return
}

func Retrieve(c *gin.Context) {
	req, err := getReq(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
		return
	}
	ctx := context.Background()
	lists, err := ctl.GetController().Retrieve(ctx, req)
	if err != nil {
		c.JSON(500, gin.H{"code": 500, "msg": err.Error()})
		return
	}
	type resp struct {
		Content string  `json:"content"`
		Score   float64 `json:"score"`
		Source  string  `json:"source"`
	}
	var resps []*resp
	for _, list := range lists {
		resps = append(resps, &resp{
			Content: list.Doc.Content,
			Score:   list.Score,
			Source:  list.Source,
		})
	}
	//by, _ := json.Marshal(resps)
	c.JSON(http.StatusOK, gin.H{"code": 200, "data": resps})
}

//func Retrieve(c *gin.Context) {
//	req, err := getReq(c)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": err.Error()})
//		return
//	}
//
//
//}
