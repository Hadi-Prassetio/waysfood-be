package userdto

type UpdateUser struct {
	Fullname string `json:"fullname" validate:"required"`
	Image    string `json:"image"`
	Email    string `json:"email" validate:"required"`
	Phone    string `json:"phone" validate:"required"`
	Location string `json:"location" validate:"required"`
}
