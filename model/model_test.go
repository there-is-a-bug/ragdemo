package model

import (
	"fmt"
	"testing"
)

func TestModel(t *testing.T) {
	p := &Properties{
		UserID: &Property{
			T: long,
		},
		ID: &Property{
			T: keyword,
		},
		Title: &Property{
			T: text,
		},
		Content: &Property{
			T: text,
		},
		Vec: &Vector{
			T:       vector,
			Dims:    1024,
			Index:   true,
			Similar: "cosine",
		},
	}
	fmt.Println(p.Mappings())
}
