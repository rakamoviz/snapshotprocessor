package auth

import "context"

type ApiClient struct {
	Name string
}

type Service interface {
	ValidateApiKey(ctx context.Context, apiKey string) (ApiClient, bool, error)
}

type memoryBasedService struct {
	apiClients map[string]ApiClient
}

func NewMemoryBasedService(apiClients map[string]ApiClient) *memoryBasedService {
	return &memoryBasedService{
		apiClients: apiClients,
	}
}

func (s *memoryBasedService) ValidateApiKey(ctx context.Context, apiKey string) (ApiClient, bool, error) {
	val, ok := s.apiClients[apiKey]

	return val, ok, nil
}
