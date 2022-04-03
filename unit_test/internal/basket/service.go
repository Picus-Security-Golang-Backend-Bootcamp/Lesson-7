package basket

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"time"
)

type Service interface {
	Get(ctx context.Context, id string) (*Basket, error)
	GetByCustomerId(ctx context.Context, customerId string) (*Basket, error)
	Create(ctx context.Context, buyer string) (*Basket, error)
	Delete(ctx context.Context, id string) (*Basket, error)

	UpdateItem(ctx context.Context, basketId, itemId string, quantity int) error
	AddItem(ctx context.Context, basketId, sku string, quantity int, price float64) (string, error)
	DeleteItem(ctx context.Context, basketId, itemId string) error
}
type service struct {
	repo Repository
}

// NewService creates a new basket service.
func newService(repo Repository) Service {

	if repo == nil {
		return nil
	}
	return &service{repo}
}

func (s *service) Get(ctx context.Context, id string) (basket *Basket, err error) {
	if len(id) == 0 {
		return nil, errors.New("Id cannot be nil or empty")
	}
	basket = s.repo.Get(ctx, id)
	if err != nil {
		err = errors.Wrapf(err, "get basket error. Basket Id:%s", id)
	}

	return
}
func (s *service) GetByCustomerId(ctx context.Context, customerId string) (basket *Basket, err error) {

	basket = s.repo.GetByCustomerId(ctx, customerId)
	if err != nil {
		err = errors.Wrapf(err, "get basket error. Customer Id:%s", customerId)
	}

	return
}

// Create creates a new basket
func (s *service) Create(ctx context.Context, customerId string) (*Basket, error) {

	basket := &Basket{
		Id:         uuid.New().String(),
		CustomerId: customerId,
		Items:      nil,
		CreatedAt:  time.Now(),
	}
	err := s.repo.Create(ctx, basket)

	if err != nil {
		return nil, errors.Wrap(err, "Service:Failed to create basket")
	}
	return basket, nil

}

func (s *service) AddItem(ctx context.Context, basketId, sku string, quantity int, price float64) (string, error) {

	basket := s.repo.Get(ctx, basketId)
	if basket == nil {
		return "", errors.Errorf("Service: Get basket error. Basket Id : %s", basketId)
	}
	if basket == nil {
		return "", errors.New("Service: Basket not found")
	}
	item, err := basket.AddItem(sku, price, quantity)

	if err != nil {
		return "", errors.Wrap(err, "Service: Failed to item added to basket.")
	}
	if err := s.repo.Update(ctx, *basket); err != nil {
		return "", errors.Wrap(err, "Service: Failed to update basket in data storage.")
	}

	return item.Id, nil
}
func (s *service) UpdateItem(ctx context.Context, basketId, itemId string, quantity int) error {

	basket := s.repo.Get(ctx, basketId)
	if basket == nil {
		return errors.Errorf("Service: Get basket error. Basket Id:%s", basketId)
	}
	if basket == nil {
		return errors.New("Service: Basket not found")
	}
	err := basket.UpdateItem(itemId, quantity)

	if err != nil {
		return errors.Wrapf(err, "Service: Failed to update item")
	}
	if err := s.repo.Update(ctx, *basket); err != nil {
		return errors.Wrap(err, "Service: Failed to update basket in data storage.")
	}
	return nil
}

func (s *service) DeleteItem(ctx context.Context, basketId, itemId string) error {

	basket := s.repo.Get(ctx, basketId)
	if basket == nil {
		return errors.Errorf("Service: Get basket error. Basket Id:%s", basketId)
	}
	if basket == nil {
		return errors.New("Service: Basket not found")
	}
	err := basket.RemoveItem(itemId)
	if err != nil {
		return errors.Wrap(err, "Service: Basket Item not found.")
	}
	if err := s.repo.Update(ctx, *basket); err != nil {
		return errors.Wrap(err, "Service: Failed to update basket in data storage.")
	}
	return nil
}

//Delete deletes the basket with the spesified Id
func (s *service) Delete(ctx context.Context, id string) (*Basket, error) {
	basket, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if basket == nil {
		return nil, errors.New("Service: Basket not found")
	}
	if err = s.repo.Delete(ctx, *basket); err != nil {
		return nil, errors.Wrap(err, "Service:Failed to delete basket")
	}
	return basket, nil
}
