package basket

import (
	"context"
	"fmt"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strconv"
	"testing"
	"time"
)

/*
https://betterprogramming.pub/how-to-unit-test-a-gorm-application-with-sqlmock-97ee73e36526
https://github.com/go-gorm/gorm/issues/3565
https://medium.com/@rosaniline/unit-testing-gorm-with-go-sqlmock-in-go-93cbce1f6b5b
*/
func TestService(usecase *testing.T) {

	usecase.Run("NewService", func(t *testing.T) {
		type args struct {
			repo Repository
		}
		tests := []struct {
			name string
			args args
			want Service
		}{
			{name: "WithValidArgs_ShouldSuccess", args: args{repo: &mockRepository{}}, want: &service{repo: &mockRepository{}}},
			{name: "WithNullRepo_ShouldReturnNil", args: args{nil}, want: nil},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {

				if got := newService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("NewService() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	usecase.Run("ReadMethods", func(t *testing.T) {

		givenBasket := Basket{
			Id:         "ID_1",
			CustomerId: "Customer",
			CreatedAt:  time.Now(),
		}

		mockRepo := &mockRepository{items: []Basket{givenBasket}}
		loadData(mockRepo)
		s := newService(mockRepo)

		t.Run("Get Method Tests", func(t *testing.T) {

			tests := []struct {
				name       string
				args       string
				wantBasket *Basket
				wantErr    bool
			}{
				{name: "WithEmptyBasket_ShouldNotFoundError", args: "INVALID_ID", wantBasket: nil, wantErr: false},
				{name: "WithEmptyBasket", args: "ID_1", wantBasket: &givenBasket, wantErr: false},
			}
			ctx := context.Background()
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					gotBasket, err := s.Get(ctx, tt.args)
					if (err != nil) != tt.wantErr {
						t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
						return
					}
					if diff := cmp.Diff(gotBasket, tt.wantBasket); diff != "" {
						t.Errorf("Get() mismatch (-want +got):\n%s", diff)
					}
				})
			}
		})
	})

	usecase.Run("Crud Operations", func(t *testing.T) {

		givenBasket := Basket{
			Id:         "CantDelete",
			CustomerId: "Buyer",
			CreatedAt:  time.Now(),
		}
		mockRepo := &mockRepository{items: []Basket{givenBasket}}
		loadData(mockRepo)
		s := newService(mockRepo)
		ctx := context.Background()

		t.Run("CreateBasket", func(t *testing.T) {
			t.Run("WithValidBuyer_ShouldBeSuccess", func(t *testing.T) {
				got, err := s.Create(ctx, "Buyer-X")
				if err != nil {
					t.Errorf("Count() error = %v", err)
					return
				}
				assert.NotNil(t, got)

			})
			t.Run("WithErrorBuyer_ShouldBeFailed", func(t *testing.T) {

				_, err := s.Create(ctx, "error")
				assert.Equal(t, errCRUD, errors.Cause(err))

			})
		})
		t.Run("DeleteBasket", func(t *testing.T) {
			t.Run("WithValidBasket_ShouldBeSuccsess", func(t *testing.T) {
				given, _ := s.Get(ctx, "ID_1")
				got, err := s.Delete(ctx, "ID_1")
				assert.NoError(t, err)
				if diff := cmp.Diff(got, given); diff != "" {
					t.Errorf("Get() mismatch (-want +got):\n%s", diff)
				}

			})
			t.Run("WithNotFoundBasketId_ShouldBeFailed", func(t *testing.T) {
				_, err := s.Delete(ctx, "NotFound")
				t.Log(err)
				assert.Equal(t, err.Error(), "Service: Basket not found")

			})
			t.Run("WithemptyBasketId_ShouldBeFailed", func(t *testing.T) {
				_, err := s.Delete(ctx, "")
				t.Log(err)
				assert.Equal(t, "Id cannot be nil or empty", errors.Cause(err).Error())

			})

			t.Run("WithExistBasketIdButCantDelete_ShouldBeFailed", func(t *testing.T) {
				_, err := s.Delete(ctx, "CantDelete")
				t.Log(err)
				assert.Equal(t, errCRUD, errors.Cause(err))

			})
		})
		t.Run("AddItem", func(t *testing.T) {

			tests := []struct {
				name    string
				args    []string
				want    string
				wantErr bool
				wantStr string
			}{
				{name: "WithValidItem_ShouldBeAdded", args: []string{"ID_5", "SKU_2", "5", "10"}, want: "", wantErr: false, wantStr: ""},
				{name: "WithExistItem_ShouldBeFailed", args: []string{"ID_5", "SKU_5", "5", "10"}, want: "", wantErr: true, wantStr: "Service: Item already added"},
				{name: "WithNonExistBasket_ShouldBeFailed", args: []string{"INVALID_BASKET_ID", "SKU_1", "5", "10"}, want: "", wantErr: true, wantStr: "Service: Get basket error. Basket Id : INVALID_BASKET_ID"},
				{name: "WithNonExistBasket_ShouldBeFailed", args: []string{"", "SKU_1", "5", "10"}, want: "", wantErr: true, wantStr: "Service: Get basket error. Basket Id : "},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					_, err := s.AddItem(ctx, tt.args[0], tt.args[1], castStrToInt(tt.args[2]), float64(castStrToInt(tt.args[3])))
					if (err != nil) != tt.wantErr {
						t.Errorf("AddItem() error = %v, wantErr %v", err, tt.wantErr)
					}
					if (err != nil) && tt.wantErr && errors.Cause(err).Error() != tt.wantStr {
						t.Errorf("AddItem() error = %v, wantErr %v", errors.Cause(err).Error(), tt.wantStr)
					}
				})
			}
		})
		t.Run("UpdateItem", func(t *testing.T) {

			tests := []struct {
				name    string
				args    []string
				wantErr bool
				wantStr string
			}{
				{name: "WithValidItemParameters_ShouldBeSuccess", args: []string{"ID_UPDATE", "ITEM_UPDATE", "8"}, wantErr: false, wantStr: ""},
				{name: "WithInvalidItemIdParameters_ShouldBeSuccess", args: []string{"ID_UPDATE", "INVALID_ITEM_ID", "3"}, wantErr: true, wantStr: "Item can not found. ItemId : INVALID_ITEM_ID"},
				{name: "WithInvalidBasketIdParameters_ShouldBeSuccess", args: []string{"ID_INVALID_ID", "ITEM_UPDATE", "5"}, wantErr: true, wantStr: "Service: Get basket error. Basket Id:ID_INVALID_ID"},
				{name: "WithInvalidBasketIdParameters_ShouldBeSuccess", args: []string{"", "ITEM_UPDATE", "4"}, wantErr: true, wantStr: "Service: Get basket error. Basket Id:"},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					err := s.UpdateItem(ctx, tt.args[0], tt.args[1], castStrToInt(tt.args[2]))
					if (err != nil) != tt.wantErr {
						t.Errorf("UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
					}
					if (err != nil) && tt.wantErr && errors.Cause(err).Error() != tt.wantStr {
						t.Errorf("AddItem() error = %v, wantErr %v", errors.Cause(err).Error(), tt.wantStr)
					}
				})
			}
		})
		t.Run("DeleteItem", func(t *testing.T) {

			tests := []struct {
				name    string
				args    []string
				wantErr bool
			}{
				{name: "WithValidItem_ShouldBeSuccess", args: []string{"ID_DELETE", "ITEM_DELETE"}, wantErr: false},
				{name: "WithValidItem_ShouldBeSuccess", args: []string{"ID_DELETE", "INVALID_ITEM_DELETE"}, wantErr: true},
				{name: "WithValidItem_ShouldBeSuccess", args: []string{"INVALID_ID_DELETE", "ITEM_DELETE"}, wantErr: true},
				{name: "WithValidItem_ShouldBeSuccess", args: []string{"", "ITEM_DELETE"}, wantErr: true},
			}
			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {

					if err := s.DeleteItem(ctx, tt.args[0], tt.args[1]); (err != nil) != tt.wantErr {
						t.Errorf("DeleteItem() error = %v, wantErr %v", err, tt.wantErr)
					}
				})
			}

		})
	})
}

/*
MockRepository here
*/
var errCRUD = errors.New("Mock: Error crud operation")

type mockRepository struct {
	items []Basket
}

func (m mockRepository) Get(ctx context.Context, id string) *Basket {
	if len(id) == 0 {
		return nil
	}

	for _, item := range m.items {
		if item.Id == id {
			return &item
		}
	}
	return nil
}
func (m mockRepository) GetByCustomerId(ctx context.Context, customerId string) *Basket {

	if len(customerId) == 0 {
		return nil
	}

	for _, item := range m.items {
		if item.CustomerId == customerId {
			return &item
		}
	}
	return nil
}
func (m *mockRepository) Create(ctx context.Context, basket *Basket) error {
	if basket.CustomerId == "error" {
		return errCRUD
	}
	m.items = append(m.items, *basket)
	return nil
}

func (m *mockRepository) Update(ctx context.Context, basket Basket) error {
	if basket.CustomerId == "error" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.Id == basket.Id {
			m.items[i] = basket
			break
		}
	}
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, basket Basket) error {
	if basket.Id == "CantDelete" {
		return errCRUD
	}
	for i, item := range m.items {
		if item.Id == basket.Id {
			m.items[i] = m.items[len(m.items)-1]
			m.items = m.items[:len(m.items)-1]
			break
		}
	}
	return nil
}

func loadData(repo *mockRepository) {
	ctx := context.TODO()

	for i := 1; i < 100; i++ {

		basket := &Basket{
			Id:         fmt.Sprintf("ID_%v", i),
			CustomerId: fmt.Sprintf("Customer_%v", i),
			CreatedAt:  time.Now(),
		}
		basket.AddItem(fmt.Sprintf("SKU_%v", i), 5, 1)
		repo.Create(ctx, basket)
	}

	//for UpdateItem
	repo.Create(ctx, &Basket{
		Id:         "ID_UPDATE",
		CustomerId: "Customer",
		Items: []Item{{
			Id:        "ITEM_UPDATE",
			Sku:       "SKU",
			UnitPrice: 5,
			Quantity:  2,
		}},
		CreatedAt: time.Now(),
	})
	//for DeleteItem
	repo.Create(ctx, &Basket{
		Id:         "ID_DELETE",
		CustomerId: "Customer",
		Items: []Item{{
			Id:        "ITEM_DELETE",
			Sku:       "SKU",
			UnitPrice: 5,
			Quantity:  2,
		}},
		CreatedAt: time.Now(),
	})

}

func castStrToInt(s string) int {
	if i, ok := strconv.Atoi(s); ok == nil {
		return i
	}
	return 0
}
