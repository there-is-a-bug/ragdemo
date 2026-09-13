package embed

import (
	"fmt"
	"testing"
)

func TestEmbed(t *testing.T) {
	b, _ := Embedding("test world")
	fmt.Println(b)
}
