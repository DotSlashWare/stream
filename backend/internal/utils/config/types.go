package config

type TMDBServiceConfig struct {
	TMDBAPIUrl string `json:"tmdb_api_url"`
	TMDBAPIKey string `json:"tmdb_api_key"`
}

type InvidiousServiceConfig struct {
	VideoAPIUrl string `json:"video_api_url"`
	VideoAPIKey string `json:"video_api_key"`
}

type LocalServiceConfig struct {
	MediaPath string `json:"media_path"`
}
