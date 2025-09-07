package application

import (
	"arch/internal/infrastructure"
)

type Application struct {
}

func New(post *infrastructure.Infrastructure) *Application {
	return &Application{}
}
