package repositories

import (
	"waysfood/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	FindOrder() ([]models.Order, error)
	GetOrder(ID int) (models.Order, error)
	CreateOrder(Order models.Order) (models.Order, error)
	DeleteOrder(Order models.Order, ID int) (models.Order, error)
	GetCartID(UserID int) (models.Cart, error)
	UpdateOrder(Order models.Order, ID int) (models.Order, error)
}

func RepositoryOrder(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) FindOrders() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Preload("Product").Find(&orders).Error

	return orders, err
}

func (r *repository) GetOrder(ID int) (models.Order, error) {
	var order models.Order
	err := r.db.Preload("Product").First(&order, ID).Error

	return order, err
}

func (r *repository) GetCartID(UserID int) (models.Cart, error) {
	var cart models.Cart
	err := r.db.Where("user_id = ? AND status = ?", UserID, "pending").First(&cart).Error
	return cart, err
}

func (r *repository) CreateOrder(order models.Order) (models.Order, error) {
	err := r.db.Create(&order).Error

	return order, err
}

func (r *repository) DeleteOrder(order models.Order, ID int) (models.Order, error) {
	err := r.db.Delete(&order).Error

	return order, err
}
func (r *repository) UpdateOrder(order models.Order, ID int) (models.Order, error) {
	err := r.db.Save(&order).Error

	return order, err
}
