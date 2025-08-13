package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (service *Service) GetMovieById(id string) (*MovieData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), service.ContextTimeout)
	defer cancel()

	url := service.getBaseApiEndpoint(fmt.Sprintf("movie/%s", id))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for movie ID %s: %v.", id, err)
		return nil, err
	}

	resp, err := service.httpClient.Do(req)
	if err != nil {
		log.Printf("Error fetching movie ID %s: %v.", id, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned status %d for movie ID %s.", resp.StatusCode, id)
		return nil, err
	}

	var movieData MovieData
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&movieData)
	if err != nil {
		log.Printf("Error decoding response for movie ID %s: %v.", id, err)
		return nil, err
	}

	return &movieData, nil
}

func (service *Service) SearchForMovie(query string, page int) (*SearchResults, error) {
	ctx, cancel := context.WithTimeout(context.Background(), service.ContextTimeout)
	defer cancel()

	url := service.getBaseApiEndpoint("search/movie", map[string]string{
		"query": query,
		"page":  fmt.Sprint(page),
	})

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for search query '%s': %v.", query, err)
		return nil, err
	}

	resp, err := service.httpClient.Do(req)
	if err != nil {
		log.Printf("Error fetching search results for query '%s': %v.", query, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned status %d for search query '%s'.", resp.StatusCode, query)
		return nil, err
	}

	var searchResults SearchResults
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&searchResults)
	if err != nil {
		log.Printf("Error decoding response for search query '%s': %v.", query, err)
		return nil, err
	}

	return &searchResults, nil
}
