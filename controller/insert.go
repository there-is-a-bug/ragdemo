package controller

import (
	"MyMemory/embed"
	"MyMemory/logs"
	"MyMemory/model"
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

func (c *Controller) Insert(ctx context.Context, req *model.Request) error {
	err := c.insertMysql(ctx, req)
	if err != nil {
		logs.Error("insert mysql fail", err)
		return err
	}
	contents := c.chunk.Chunk(req.Content, req.Chunk, 0)
	worker := make([][]string, 10)
	for i, content := range contents {
		worker[i%10] = append(worker[i%10], content)
	}
	wg := &sync.WaitGroup{}
	var esErr atomic.Value
	wg.Add(len(worker))
	for i, w := range worker {
		offset := i
		go func() {
			defer wg.Done()
			for i, content := range w {
				err1 := c.insertEs(ctx, &model.Request{
					Content: content,
					ID:      req.ID + fmt.Sprintf("-%d-%d", offset, i),
					UID:     req.UID,
					Title:   req.Title,
				})
				if err1 != nil {
					logs.Error("chunk[%d]insert es fail", offset+i*10, err1)
					esErr.Store(err1)
				}
			}
		}()
	}
	wg.Wait()
	if e := esErr.Load(); e != nil {
		logs.Error("some chunk insert es partial fail")
		return fmt.Errorf("some chunk insert es partial fail")
	}
	return nil
}

func (c *Controller) insertMysql(ctx context.Context, req *model.Request) error {
	db := c.sqlRepo
	err := db.Create(ctx, &model.UserDocument{
		DocID:   req.ID,
		UserID:  req.UID,
		Content: req.Content,
		Title:   req.Title,
		DocType: "content",
	})
	if err != nil {
		logs.Error("mysql: insert document error:", err)
		return fmt.Errorf("mysql: insert document error: %s", err.Error())
	}
	return nil
}

func (c *Controller) insertEs(ctx context.Context, req *model.Request) error {
	embedding, err := embed.Embedding(req.Content)
	if err != nil {
		logs.Error("embedding error:", err)
		return fmt.Errorf("embedding error: %s", err.Error())
	}
	repo := c.esRepo
	err = repo.Insert(ctx, &model.Document{
		ID:      req.ID,
		Content: req.Content,
		UserID:  req.UID,
		Vec:     embedding,
		Title:   req.Title,
	})
	if err != nil {
		logs.Error("es: insert document error:", err)
		return fmt.Errorf("es: insert document error: %s", err.Error())
	}
	return nil
}
