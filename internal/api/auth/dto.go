package auth

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	User        struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Role     string `json:"role"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"user"`
}

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}
