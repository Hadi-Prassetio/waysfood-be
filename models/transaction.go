package models

type Transaction struct {
	ID      int         `json:"id"  gorm:"primary_key:auto_increment"`
	BuyerID int         `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Buyer   UserProfile `json:"buyer"`
	CartID  int         `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Cart    Cart        `json:"cart"`
	Total   int         `json:"total"`
	Status  string      `json:"status"`
}
