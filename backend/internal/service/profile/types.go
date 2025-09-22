package profile

type Profile struct {
	Id        int    `json:"id"`
	Username  string `json:"username"`
	IsChild   bool   `json:"is_child"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	LastLogin string `json:"last_login"`
}
