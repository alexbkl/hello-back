package form

type RegisterUserRequest struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type LoginUserRequest struct {
	Name string `json:"name"`
	Role string `json:"role"`
}
