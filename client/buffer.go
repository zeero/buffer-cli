package client

import (
	"context"

	"github.com/machinebox/graphql"
)

type BufferClient struct {
	gql   *graphql.Client
	token string
}

func NewWithEndpoint(token, endpoint string) *BufferClient {
	return &BufferClient{
		gql:   graphql.NewClient(endpoint),
		token: token,
	}
}

func New(token string) *BufferClient {
	return NewWithEndpoint(token, "https://api.buffer.com")
}

func (c *BufferClient) Run(ctx context.Context, query string, vars map[string]interface{}, resp interface{}) error {
	req := graphql.NewRequest(query)
	req.Header.Set("Authorization", "Bearer "+c.token)

	for k, v := range vars {
		req.Var(k, v)
	}

	return c.gql.Run(ctx, req, resp)
}
