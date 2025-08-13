package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (service *Service) GetMovieById(id int) *MovieData {
	ctx, cancel := context.WithTimeout(context.Background(), service.ContextTimeout)
	defer cancel()

	url := service.getBaseApiEndpoint(fmt.Sprintf("movie/%d", id))

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for movie ID %d: %v", id, err)
		return nil
	}

	resp, err := service.httpClient.Do(req)
	if err != nil {
		log.Printf("Error fetching movie ID %d: %v", id, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned status %d for movie ID %d", resp.StatusCode, id)
		return nil
	}

	var movieData MovieData
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&movieData)
	if err != nil {
		log.Printf("Error decoding response for movie ID %d: %v", id, err)
		return nil
	}

	return &movieData
}

func (service *Service) SearchForMovie(query string, page int) *SearchResults {
	ctx, cancel := context.WithTimeout(context.Background(), service.ContextTimeout)
	defer cancel()

	url := service.getBaseApiEndpoint("search/movie", map[string]string{
		"query": query,
		"page":  fmt.Sprint(page),
	})

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Printf("Error creating request for search query '%s': %v", query, err)
		return nil
	}

	resp, err := service.httpClient.Do(req)
	if err != nil {
		log.Printf("Error fetching search results for query '%s': %v", query, err)
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("API returned status %d for search query '%s'", resp.StatusCode, query)
		return nil
	}

	var searchResults SearchResults
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&searchResults)
	if err != nil {
		log.Printf("Error decoding response for search query '%s': %v", query, err)
		return nil
	}

	return &searchResults
}
