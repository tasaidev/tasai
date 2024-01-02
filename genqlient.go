package tasai

import (
	"net/http"

	"github.com/Khan/genqlient/graphql"
)

//go:generate go run github.com/Khan/genqlient genqlient.yaml

type authedTransport struct {
	key     *string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.key == nil {
		return t.wrapped.RoundTrip(req)
	}
	req.Header.Set("Authorization", "Bearer "+*t.key)
	return t.wrapped.RoundTrip(req)
}

func login(key string) graphql.Client {
	return getGenQlient(key)
}

func getGenQlient(token ...string) graphql.Client {
	transport := &authedTransport{
		wrapped: http.DefaultTransport,
	}
	if len(token) > 0 {
		transport.key = &token[0]
	}
	httpClient := http.Client{
		Transport: transport,
	}
	graphqlClient := graphql.NewClient("http://localhost:8080/graphql", &httpClient)
	return graphqlClient
}
