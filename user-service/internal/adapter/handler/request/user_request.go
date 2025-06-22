package request

type SignInRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type SignUpRequest struct {
	Name                 string `json:"name" validate:"required"`
	Email                string `json:"email" validate:"required,email"`
	Password             string `json:"password" validate:"required,min=8"`
	PasswordConfirmation string `json:"password_confirmation" validate:"required,min=8"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdatePasswordRequest struct {
	CurrentPassword string `json:"password,omitempty"`
	NewPassword     string `json:"password_new" validate:"required,min=8"`
	ConfirmPassword string `json:"password_confirmation" validate:"required,min=8"`
}
