package client

import (
	"context"
	"github.com/Mohammadmohebi33/toll_calculator/types"
)

type Client interface {
	Aggregate(ctx context.Context, request *types.AggregateRequest) error
}
