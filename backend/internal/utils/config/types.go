package config

type TMDBServiceConfig struct {
	TMDBAPIUrl   string `json:"tmdbapi_url"`
	TMDBAPIKey   string	`json:"tmdbapi_key,omitempty"`
}

type InvidiousServiceConfig struct {
	VideoAPIUrl string `json:"videoapi_url"`
	VideoAPIKey string `json:"videoapi_key,omitempty"`
}

type LocalServiceConfig struct {
	MediaPath string `json:"save_path"`
}

type MovieConfig struct {
	StreamAPIUrl string `json:"streamapi_url"`
}

type TVConfig struct {
	StreamAPIUrl string `json:"streamapi_url"`
}

type VideoConfig struct {
	StreamAPIUrl string `json:"streamapi_url"`
}