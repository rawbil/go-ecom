package authutils

type UserRegisterParams struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password_format"`
}

type UserLoginParams struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}