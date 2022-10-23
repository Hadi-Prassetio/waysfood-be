package userdto

type UpdateUser struct {
	Fullname string `json:"fullname" validate:"required"`
	Image    string `json:"image" validate:"reqired"`
	Email    string `json:"email" validate:"reqired"`
	Phone    string `json:"phone" validate:"reqired"`
	Location string `json:"location" validate:"reqired"`
}
