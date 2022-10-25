package cartdto

type CreateCart struct {
	ID       int `json:"id"`
	UserID   int `json:"user_id"`
	SubTotal int `json:"sub_total"`
	Qty      int `json:"qty"`
}

type UpdateCart struct {
	Status   string `json:"status"`
	UserID   int    `json:"user_id"`
	SubTotal int    `json:"sub_total"`
	Qty      int    `json:"qty"`
}
