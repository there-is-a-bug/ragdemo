package embed

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"
)

type EmbResp struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
	} `json:"data"`
	Model string `json:"model"`
}

const url = "http://192.168.1.212:8080/embed"

func Embedding(input string) ([]float32, error) {
	buf, _ := json.Marshal(map[string]any{
		"inputs": input,
	})
	client := http.Client{Timeout: 10 * time.Second}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(buf))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	raw, _ := io.ReadAll(resp.Body)
	var data [][]float32
	_ = json.Unmarshal(raw, &data)
	return data[0], nil
}
