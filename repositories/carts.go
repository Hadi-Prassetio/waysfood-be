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
	UpdateCarts(status string, ID string) error
	FindbyIDCart(CartId int, Status string) (models.Cart, error)
	GetOneCart(ID int) (models.Cart, error)
	AllProductById(UserID int) ([]models.Cart, error)
}

func RepositoryCart(db *gorm.DB) *repository {
	return &repository{db}
}
func (r *repository) FindCarts() ([]models.Cart, error) {
	var Carts []models.Cart
	err := r.db.Preload("User").Preload("Carts").Preload("Carts.Product").Find(&Carts).Error
	return Carts, err
}

func (r *repository) GetCart(ID int) (models.Cart, error) {
	var Cart models.Cart
	err := r.db.Preload("User").Preload("Carts").Preload("Carts.Product").Find(&Cart, ID).Error
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

func (r *repository) UpdateCarts(status string, ID string) error {
	var Cart models.Cart
	r.db.Preload("User").Preload("Order").Preload("Carts.Product").First(&Cart, ID)

	// If is different & Status is "success" decrement product quantity
	if status != Cart.Status && status == "success" {
		var product models.Product
		r.db.First(&product, Cart.ID)
		// product.Qty = product.Qty - 1
		r.db.Save(&product)
	}

	Cart.Status = status

	err := r.db.Save(&Cart).Error

	return err
}

func (r *repository) FindbyIDCart(CartId int, Status string) (models.Cart, error) {
	var Cart models.Cart
	err := r.db.Preload("User").Preload("Carts").Preload("Carts.Product").Where("user_id = ? AND status = ?", CartId, Status).First(&Cart).Error

	return Cart, err
}

func (r *repository) GetOneCart(ID int) (models.Cart, error) {
	var Cart models.Cart
	err := r.db.Preload("Product").Preload("Product.User").Preload("Buyer").Preload("Seller").First(&Cart, "id = ?", ID).Error

	return Cart, err
}

func (r *repository) AllProductById(UserID int) ([]models.Cart, error) {
	var Cart []models.Cart
	err := r.db.Preload("User").Preload("Carts").Preload("Carts.Product").Find(&Cart, "user_id = ?", UserID).Error

	return Cart, err
}
