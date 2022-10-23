package models

type Order struct {
	ID        int     `json:"id"  gorm:"primary_key:auto_increment"`
	CartID    int     `json:"-"`
	ProductID int     `json:"-"`
	Product   Product `json:"product"`
	Qty       int     `json:"qty"`
	SubAmount int     `json:"sub_amount"`
}
