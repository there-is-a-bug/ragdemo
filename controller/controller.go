package controller

import (
	"MyMemory/chunk"
	"MyMemory/es"
	"MyMemory/mysql"
	"sync"
)

type Controller struct {
	esRepo  *es.ESRepo
	sqlRepo *mysql.UserDocumentRepo
	chunk   chunk.Chunked
}

var (
	so         sync.Once
	controller *Controller
)

func init() {
	so.Do(func() {
		controller = &Controller{
			esRepo:  es.GetESRepo(),
			sqlRepo: mysql.GetUserDocumentRepo(),
			chunk:   chunk.NewTextChunk(),
		}
	})
}

func GetController() *Controller {
	return controller
}
