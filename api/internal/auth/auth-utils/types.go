package authutils

type UserRegisterParams struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password_format,max=12"`
}

type UserLoginParams struct {
	Email string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}


type PasswordResetParams struct {
	NewPassword string `json:"new_password" validate:"required,password_format,min=8,max=12"`
	OldPassword string `json:"old_password" validate:"required"`
}