package application

import (
	"arch/internal/domain/ai"
	"arch/internal/domain/entity"
	"context"
)

func (a *Application) SendAi(ctx context.Context, input entity.Send) ([]entity.AiPlace, error) {
	q := ai.New(*a.qwqConfig)
	output, err := q.Send(ctx, input.Message)
	if err != nil {
		return nil, err
	}

	return output, nil
}
