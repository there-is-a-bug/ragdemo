package chunk

type Chunked interface {
	Chunk(text string, chunkSize, overlap int) []string
}
