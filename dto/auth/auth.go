package authdto

type RequestRegister struct {
	Email    string `json:"email" validate:"required"`
	Fullname string `json:"fullname" validate:"required"`
	Password string `json:"password" validate:"required"`
	Gender   string `json:"gender" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Role     string `json:"role" validate:"required"`
}

type RequestLogin struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ResponseLogin struct {
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Password string `json:"password" `
	Gender   string `json:"gender"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
	Token    string `json:"token"`
}
