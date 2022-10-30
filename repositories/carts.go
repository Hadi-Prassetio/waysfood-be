package repositories

import (
	"waysfood/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	FindCarts() ([]models.Cart, error)
	GetCart(ID int) (models.Cart, error)
	CreateCart(Cart models.Cart) (models.Cart, error)
	DeleteCart(Cart models.Cart) (models.Cart, error)
	UpdateCart(Cart models.Cart) (models.Cart, error)
	FindbyIDCart(CartId int, Status string) (models.Cart, error)
	GetOneCart(ID int) (models.Cart, error)
}

func RepositoryCart(db *gorm.DB) *repository {
	return &repository{db}
}
func (r *repository) FindCarts() ([]models.Cart, error) {
	var Carts []models.Cart
	err := r.db.Preload("User").Preload("Order").Preload("Order.Product").Find(&Carts).Error
	return Carts, err
}

func (r *repository) GetCart(ID int) (models.Cart, error) {
	var Cart models.Cart
	err := r.db.Preload("User").Preload("Order").Preload("Order.Product").Find(&Cart, ID).Error
	return Cart, err
}

func (r *repository) CreateCart(Cart models.Cart) (models.Cart, error) {
	err := r.db.Create(&Cart).Error

	return Cart, err
}

func (r *repository) UpdateCart(Cart models.Cart) (models.Cart, error) {
	err := r.db.Save(&Cart).Error

	return Cart, err
}

func (r *repository) DeleteCart(Cart models.Cart) (models.Cart, error) {
	err := r.db.Delete(&Cart).Error

	return Cart, err
}

func (r *repository) FindbyIDCart(CartId int, Status string) (models.Cart, error) {
	var Cart models.Cart
	err := r.db.Preload("User").Preload("Order").Preload("Order.Product").Preload("Order.Product.User").Where("user_id = ? AND status = ?", CartId, Status).First(&Cart).Error

	return Cart, err
}

func (r *repository) GetOneCart(ID int) (models.Cart, error) {
	var Cart models.Cart
	err := r.db.Preload("Product").Preload("Product.User").Preload("Buyer").Preload("Seller").First(&Cart, "id = ?", ID).Error

	return Cart, err
}
