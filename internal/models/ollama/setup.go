package models_ollama

import "context"

type OllamaClient struct {
	ctx context.Context
}

var SetupOllama = func(ctx context.Context) *OllamaClient {
	return &OllamaClient{
		ctx: ctx,
	}
}
