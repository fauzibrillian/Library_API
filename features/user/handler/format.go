package handler

type LoginRequest struct {
	Phone    string `json:"phone" form:"phone"`
	Password string `json:"password" form:"password"`
}

type LoginResponse struct {
	ID    uint   `json:"user_id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Token string `json:"token"`
	Role  string `json:"role"`
}

type RegisterRequest struct {
	Name     string `json:"name" form:"name"`
	Phone    string `json:"phone" form:"phone"`
	Password string `json:"password" form:"password"`
}

type RegisterResponse struct {
	ID    uint   `json:"user_id"`
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type ResetPasswordRequest struct {
	ID          uint   `json:"user_id" form:"user_id"`
	Password    string `json:"password" form:"password"`
	NewPassword string `json:"newpassword" form:"newpassword"`
}

type ResetPasswordResponse struct {
	ID   uint   `json:"user_id" form:"user_id"`
	Name string `json:"name" form:"name"`
}
