package authutils

type UserRegisterParams struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,password_format,max=12"`
}

type UserLoginParams struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type PasswordResetParams struct {
	NewPassword string `json:"new_password" validate:"required,password_format,min=8,max=12"`
	OldPassword string `json:"old_password" validate:"required"`
}

type RefreshTokenParam struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type UpdateUsernameParams struct {
	Username string `json:"username" validate:"required,min=3"`
}

type UpdateEmailParams struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdateProductParams struct {
	ProductId int64 `json:"product_id" validate:"required"`
	Price     int64 `json:"price" validate:"min=0"`
	Quantity  int32 `json:"quantity" validate:"min=0"`
}
