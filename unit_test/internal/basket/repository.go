package basket

import (
	"context"
	"gorm.io/gorm"
)

// Repository encapsulates the logic to access basket from the data source.
type Repository interface {
	// Get returns the basket with the specified basket Id.
	Get(ctx context.Context, id string) *Basket
	// GetByCustomerId returns the basket with the specified customer Id.
	GetByCustomerId(ctx context.Context, customerId string) *Basket
	// Create saves a new basket in the storage.
	Create(ctx context.Context, basket *Basket) error
	// Update updates the basket with given Is in the storage.
	Update(ctx context.Context, basket Basket) error
	// Delete removes the basket with given Is from the storage.
	Delete(ctx context.Context, basket Basket) error
}

type BasketRepository struct {
	db *gorm.DB
}

func NewBasketRepository(db *gorm.DB) *BasketRepository {
	return &BasketRepository{
		db: db,
	}
}

func (r *BasketRepository) Get(ctx context.Context, id string) *Basket {
	var basket *Basket
	r.db.WithContext(ctx).Where("Id = ?", id).Find(&basket)

	return basket
}
func (r *BasketRepository) GetByCustomerId(ctx context.Context, id string) *Basket {
	var basket *Basket
	r.db.WithContext(ctx).Where("CustomerId = ?", id).Find(&basket)

	return basket
}

// Eğer pointer verilmezse reflect.Value.Set using unaddressable value hatası alınır
func (r *BasketRepository) Create(ctx context.Context, b *Basket) error {
	result := r.db.WithContext(ctx).Create(b)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BasketRepository) Update(ctx context.Context, b Basket) error {
	result := r.db.WithContext(ctx).Save(b)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BasketRepository) Delete(ctx context.Context, b Basket) error {
	result := r.db.WithContext(ctx).Delete(b)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *BasketRepository) Migration() {
	r.db.AutoMigrate(&Basket{})
}
