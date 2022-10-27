package models

type Transaction struct {
	ID       int    `json:"id"  gorm:"primary_key:auto_increment"`
	BuyerID  int    `json:"-"`
	SellerID int    `json:"-"`
	CartID   int    `json:"-"`
	Cart     Cart   `json:"cart"`
	Status   string `json:"status"`
}
