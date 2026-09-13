package chunk

import "strings"

type TextChunk struct {
}

var separator = []string{
	"\n\n",
	".", "?", "!", "。", "？", "！",
	",", ";", "，", "；",
	"\n",
}
var (
	paragraph = "\n\n"
)

func NewTextChunk() *TextChunk {
	return &TextChunk{}
}

func (t *TextChunk) Chunk(text string, chunkSize, overlap int) []string {
	return chunkBySeparators(text, chunkSize, separator)
}

func chunkBySeparators(text string, chunkSize int, sp []string) []string {
	if len([]rune(text)) <= chunkSize {
		return []string{text}
	}

	if len(sp) == 0 {
		return splitBySize(text, chunkSize)
	}

	var result, next []string
	sts := strings.Split(text, sp[0])
	//
	if len(sts) == 0 {
		next = chunkBySeparators(text, chunkSize, sp[1:])
		result = append(result, next...)
		return result
	}
	//
	result, next = extract(sts, chunkSize, sp[0])
	//
	for _, st := range next {
		sts = chunkBySeparators(st, chunkSize, sp[1:])
		result = append(result, sts...)
	}
	return cleanStr(result)
}

func extract(texts []string, chunkSize int, sp string) ([]string, []string) {
	var result, next []string
	if sp == "," || sp == "，" {
		var curr []string
		remain := chunkSize
		for _, st := range texts {
			if len([]rune(st)) > chunkSize {
				next = append(next, st)
				result = append(result, strings.Join(curr, sp))
				curr = []string{}
				remain = chunkSize
				continue
			}

			remain -= len([]rune(st))
			if remain < 0 {
				result = append(result, strings.Join(curr, sp))
				curr, remain = refresh(curr, remain, chunkSize)
			}
			curr = append(curr, st)
		}
		if len(curr) > 0 {
			result = append(result, strings.Join(curr, sp))
		}
		return result, next
	}

	for _, st := range texts {
		if len([]rune(st)) <= chunkSize {
			result = append(result, st)
			continue
		}
		next = append(next, st)
	}
	return result, next

}

func refresh(texts []string, remain, chunkSize int) ([]string, int) {
	for i, st := range texts {
		remain += len([]rune(st))
		if remain < 0 {
			continue
		}
		if i+1 >= len(texts) {
			return []string{}, remain
		}
		return texts[i+1:], chunkSize
	}
	return []string{}, chunkSize
}

func splitBySize(text string, chunkSize int) []string {
	var result []string
	runes := []rune(text)
	end := 0
	for start := 0; start < len(runes); {
		end = start + chunkSize
		if end > len(runes) {
			end = len(runes)
		}
		result = append(result, string(runes[start:end]))
		start = end
	}
	return result
}

func cleanStr(texts []string) []string {
	var result []string
	for _, text := range texts {
		t := strings.TrimSpace(text)
		if len(t) > 0 {
			result = append(result, t)
		}
	}
	return result
}
