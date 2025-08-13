package tmdb

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Service struct {
	ApiUrl         string
	ApiKey         string
	ContextTimeout time.Duration
	httpClient     *http.Client
}

func NewTmdbService(apiUrl string, apiKey string, contextTimeout int) *Service {
	return &Service{
		ApiUrl:         apiUrl,
		ApiKey:         apiKey,
		ContextTimeout: time.Duration(contextTimeout) * time.Second,
		httpClient:     &http.Client{Timeout: time.Duration(contextTimeout) * time.Second},
	}
}

func (service *Service) getBaseApiEndpoint(endpoint string, queryParams ...map[string]string) string {
	baseURL := strings.TrimRight(service.ApiUrl, "/")
	endpoint = strings.TrimLeft(endpoint, "/")

	parsedURL, err := url.Parse(baseURL + "/" + endpoint)
	if err != nil {
		return fmt.Sprintf("%s/%s?api_key=%s", baseURL, endpoint, service.ApiKey)
	}

	query := parsedURL.Query()

	query.Set("api_key", service.ApiKey)

	for _, params := range queryParams {
		for key, value := range params {
			if value != "" {
				query.Set(key, value)
			}
		}
	}

	parsedURL.RawQuery = query.Encode()

	return parsedURL.String()
}
