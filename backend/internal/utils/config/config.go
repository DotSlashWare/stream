package config

type Config struct {
    TMDBService      TMDBServiceConfig      `json:"tmdb_service"`
    InvidiousService InvidiousServiceConfig `json:"invidious_service"`
    LocalService     LocalServiceConfig     `json:"local_service"`
    MovieConfig      MovieConfig            `json:"movie_config"`
    TVConfig         TVConfig               `json:"tv_config"`
    VideoConfig      VideoConfig            `json:"video_config"`
}

func NewConfig() *Config {
	return &Config{
		TMDBService: TMDBServiceConfig{
			TMDBAPIUrl: "https://api.themoviedb.org/3",
			TMDBAPIKey: "",
		},
		InvidiousService: InvidiousServiceConfig{
			VideoAPIUrl: "https://invidious.example.com/api/v1",
			VideoAPIKey: "",
		},
		LocalService: LocalServiceConfig{
			MediaPath: "/var/media/stream",
		},
		MovieConfig: MovieConfig{
			StreamAPIUrl: "https://stream.example.com/movie",
		},
		TVConfig: TVConfig{
			StreamAPIUrl: "https://stream.example.com/tv",
		},
		VideoConfig: VideoConfig{
			StreamAPIUrl: "https://stream.example.com/video",
		},
	}
}
