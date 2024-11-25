package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Mohammadmohebi33/toll_calculator/types"
	"net/http"
)

type HttpClient struct {
	Endpoint string
}

func NewHttpClient(endpoint string) *HttpClient {
	return &HttpClient{
		Endpoint: endpoint,
	}
}

func (receiver *HttpClient) Aggregate(ctx context.Context, request *types.AggregateRequest) error {
	b, err := json.Marshal(request)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", receiver.Endpoint, bytes.NewBuffer(b))
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("invalid status code: %d", resp.StatusCode)
	}

	return nil
}
