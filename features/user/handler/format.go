package handler

type LoginRequest struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type LoginResponse struct {
	ID    uint   `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
	Role  string `json:"role"`
}

type RegisterRequest struct {
	Name     string `json:"name" form:"name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}

type RegisterResponse struct {
	ID    uint   `json:"user_id"`
	Name  string `json:"name"`
	Email string `json:"email"`
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

type PutUserRequest struct {
	ID     uint   `json:"user_id" form:"user_id"`
	Name   string `json:"name" form:"name"`
	Email  string `json:"email" form:"email"`
	Avatar string `json:"avatar" form:"avatar"`
}

type PutUserResponse struct {
	ID     uint   `json:"user_id" form:"user_id"`
	Name   string `json:"name" form:"name"`
	Email  string `json:"email" form:"email"`
	Avatar string `json:"avatar" form:"avatar"`
}
