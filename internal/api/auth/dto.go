package auth

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	User        struct {
		ID       int64   `json:"id"`
		Name     string  `json:"name"`
		Role     string  `json:"role"`
		Username string  `json:"username"`
		Email    string  `json:"email"`
		Phone    string  `json:"phone,omitempty"`
		DeptID   *int64  `json:"deptId,omitempty"`
		DeptName *string `json:"deptName,omitempty"`
	} `json:"user"`
}

type RefreshResponse struct {
	AccessToken string `json:"accessToken"`
}
