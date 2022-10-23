package models

type Transaction struct {
	ID     int    `json:"id"  gorm:"primary_key:auto_increment"`
	UserID int    `json:"-"`
	CartID int    `json:"-"`
	Cart   Cart   `json:"cart"`
	Status string `json:"status"`
}
