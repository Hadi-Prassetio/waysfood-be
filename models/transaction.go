package models

type Transaction struct {
	ID       int         `json:"id"  gorm:"primary_key:auto_increment"`
	BuyerID  int         `json:"-"`
	Buyer    UserProfile `json:"buyer"`
	SellerID int         `json:"-"`
	Seller   UserProfile `json:"seller"`
	CartID   int         `json:"-"`
	Cart     Cart        `json:"cart"`
	Total    int         `json:"total"`
	Status   string      `json:"status"`
}
