package ai

import (
	"arch/internal/domain/entity"
	"encoding/json"
	"fmt"
	"io"
)

type chatCompletionResponse struct {
	Choices []struct {
		Message entity.AiMessage `json:"message"`
	} `json:"choices"`
}

func (q *Ai) buildPayload(message string) (string, error) {
	payload := entity.AiSend{
		Model: q.ai.Model,
		Messages: []entity.AiMessage{
			{
				Role:    "assistant",
				Content: sendPrompt,
			},
			{
				Role:    entity.QwQRole,
				Content: message,
			},
		},
	}
	prompt, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(prompt), nil
}

func buildOutput(reader io.Reader) ([]entity.AiPlace, error) {
	var cc chatCompletionResponse
	if err := json.NewDecoder(reader).Decode(&cc); err != nil {
		return nil, err
	}
	if len(cc.Choices) == 0 {
		return nil, fmt.Errorf("empty choices")
	}
	content := cc.Choices[0].Message.Content
	if content == "" {
		return nil, fmt.Errorf("empty content")
	}

	var output []entity.AiPlace

	if err := json.Unmarshal([]byte(content), &output); err != nil {
		return nil, err
	}

	return output, nil
}
