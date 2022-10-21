package productdto

type CreateProduct struct {
	Title  string `json:"title" validate:"required"`
	Image  string `json:"image"`
	Price  int    `json:"price" validate:"required"`
	UserID int    `json:"user_id"`
}

type UpdateProduct struct {
	Title  string `json:"title"`
	Image  string `json:"image"`
	Price  int    `json:"price"`
	UserID int    `json:"user_id"`
}
