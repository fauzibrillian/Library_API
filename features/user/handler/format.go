package handler

type LoginRequest struct {
	Phone    string `json:"phone" form:"phone"`
	Password string `json:"password" form:"password"`
}

type LoginResponse struct {
	ID       uint   `json:"user_id"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Role     string `json:"role"`
}
