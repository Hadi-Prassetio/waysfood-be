package models

type Cart struct {
	ID       int         `json:"id"  gorm:"primary_key:auto_increment"`
	UserID   int         `json:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	User     UserProfile `json:"user"`
	Order    []Order     `json:"order"`
	Qty      int         `json:"qty"`
	SubTotal int         `json:"sub_total"`
	Status   string      `json:"status"`
}
