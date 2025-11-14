package models_openai

import "github.com/sashabaranov/go-openai"

func (oa *OpenaiClient) GetEmbedding(text string) ([]float32, error) {
	resp, err := oa.client.CreateEmbeddings(oa.ctx, openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2,
		Input: []string{text},
	})
	if err != nil {
		return nil, err
	}

	return resp.Data[0].Embedding, nil
}
