package ai

import (
	"arch/internal/domain/entity"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	methodPost = "POST"
)

var (
	timeout = time.Second * 100
)

type Ai struct {
	ai     entity.AiConfig
	client *http.Client
}

func New(cfg entity.AiConfig) *Ai {
	tr := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     100 * time.Second,
	}
	return &Ai{
		ai: cfg,
		client: &http.Client{
			Timeout:   0,
			Transport: tr,
		},
	}
}

func (q *Ai) Send(ctx context.Context, message string) ([]entity.AiPlace, error) {
	if q.client == nil {
		q.client = &http.Client{Timeout: timeout}
	}
	prompt, err := q.buildPayload(message)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, methodPost, "https://openrouter.ai/api/v1/chat/completions", strings.NewReader(prompt))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+q.ai.ApiKey)

	resp, err := q.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			fmt.Printf("close body error: %v\n", err)
		}
	}()

	if resp.StatusCode/100 != 2 {
		snippet, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		return nil, fmt.Errorf("bad status %s; body: %s", resp.Status, string(snippet))
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return nil, fmt.Errorf("499")
		}
		return nil, fmt.Errorf("read body: %w", err)
	}

	output, err := buildOutput(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	return output, nil
}
