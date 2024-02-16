package controller

import (
	"context"
	"github.com/samsarahq/thunder/reactive"
	"kroseida.org/slixx/internal/common"
	"time"
)

type EnvironmentDto struct {
	Version string `json:"version" graphql:"version"`
}

func Environment(ctx context.Context) (*EnvironmentDto, error) {
	reactive.InvalidateAfter(ctx, 5*time.Second)
	return &EnvironmentDto{
		Version: common.CurrentVersion,
	}, nil
}
