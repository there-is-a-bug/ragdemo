package chunk

import (
	"fmt"
	"testing"
)

func TestTextChunk(t *testing.T) {
	var a = `
Embedding 本质上是一种把文本等离散对象映射到连续向量空间的表示方法。在 RAG 中，通常会使用专门的 Embedding Model 将文档切分后的 Chunk 转换成固定维度的向量，然后存储到 Elasticsearch 这类支持向量检索的数据库中。

用户查询时，会使用同一个 Embedding Model 将 Query 转成向量，然后通过 Cosine Similarity、Dot Product 等相似度计算方式进行向量检索，找到语义上最相关的 Top K 文档。

Embedding Model 本身不是简单的规则函数，而是通过训练学习得到的语义表示空间。常见的训练思路包括对比学习，让 Query 和正样本的向量距离更近，让负样本距离更远。

需要注意的是，Embedding Model 和最终负责生成答案的 LLM 可以是不同模型，它们承担的任务不同：Embedding 主要负责语义表示和检索，LLM 负责基于检索到的上下文进行理解和生成。Query 和文档需要使用兼容的 Embedding 模型，否则向量不处于同一个语义空间，无法直接比较。

在实际 RAG 工程中，Embedding 只是检索链路的一部分，还需要考虑 Chunk 切分、向量维度、索引方式、Top K、Hybrid Search、Rerank 以及召回效果评估等问题。
`
	text := &TextChunk{}
	ts := text.Chunk(a, 50, 0)
	for i, t := range ts {
		fmt.Println(i, len([]rune(t)), t)
	}

}
