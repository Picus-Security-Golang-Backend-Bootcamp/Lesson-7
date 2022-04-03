package basket

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
)

var (
	basePath = "basket"
)

func RegisterHandlers(instance *gin.Engine, repo Repository) {

	res := &resource{
		service: newService(repo),
	}

	instance.GET(fmt.Sprintf("%s/:id", basePath), res.getBasket)
	instance.POST(fmt.Sprintf("%s", basePath), res.createBasket)
	instance.DELETE(fmt.Sprintf("%s/:id", basePath), res.deleteBasket)

	instance.POST(fmt.Sprintf("%s/item", basePath), res.addItem)
	instance.DELETE(fmt.Sprintf("%s/:id/item/:itemId", basePath), res.deleteItem)
	instance.PUT(fmt.Sprintf("%s/:id/item/:item/quantity/:quantity", basePath), res.updateItem)

	if gin.Mode() == gin.TestMode {
		instance.GET("/test", func(c *gin.Context) {
			c.String(418, "I don't exist in production")
		})
	}

}

type resource struct {
	service Service
}

func (r *resource) getBasket(g *gin.Context) {
	id := g.Param("id")
	result, err := r.service.Get(g.Request.Context(), id)
	if err != nil {
		g.JSON(http.StatusInternalServerError, err.Error())
	}
	if result == nil {
		g.JSON(http.StatusNotFound, "")
	}
	g.JSON(http.StatusOK, result)

}
func (r *resource) createBasket(g *gin.Context) {
	entity := new(CreateBasketRequest)

	if err := g.Bind(entity); err != nil {
		g.JSON(http.StatusBadRequest, err.Error())
	}
	if b, err := r.service.Create(g.Request.Context(), entity.Buyer); err != nil {
		g.JSON(http.StatusBadRequest, err.Error())
	} else {

		g.JSON(http.StatusCreated, map[string]string{"id": b.Id})
	}
}

func (r *resource) deleteBasket(g *gin.Context) {
	id := g.Param("id")
	_, err := r.service.Delete(g.Request.Context(), id)

	if errors.Cause(err) == sql.ErrNoRows {
		g.JSON(http.StatusNotFound, err.Error())

	}
	if err != nil {
		g.JSON(http.StatusInternalServerError, err.Error())
	}

	g.JSON(http.StatusAccepted, "")

}
func (r *resource) addItem(g *gin.Context) {
	req := new(AddItemRequest)

	if err := g.Bind(req); err != nil {
		g.JSON(http.StatusBadRequest, err.Error())
	}
	if itemId, err := r.service.AddItem(g.Request.Context(), req.BasketId, req.Sku, req.Quantity, req.Price); err != nil {
		g.JSON(http.StatusBadRequest, err.Error())
	} else {
		g.JSON(http.StatusCreated, map[string]string{"id": itemId})
	}

}
func (r *resource) updateItem(g *gin.Context) {

	id := g.Param("id")
	itemId := g.Param("itemId")
	quantity, err := strconv.Atoi(g.Param("quantity"))

	if len(id) == 0 || len(itemId) == 0 || err != nil || quantity <= 0 {
		g.JSON(http.StatusBadRequest, "Failed to update item. BasketId or BasketItem Id is null or empty.")
	}
	if err := r.service.UpdateItem(g.Request.Context(), id, itemId, quantity); err != nil {
		g.JSON(http.StatusBadRequest, err.Error())
	}
	g.JSON(http.StatusAccepted, "")

}

func (r *resource) deleteItem(g *gin.Context) {

	id := g.Param("id")
	itemId := g.Param("itemId")

	if len(id) == 0 || len(itemId) == 0 {
		g.JSON(http.StatusBadRequest, "Failed to delete item. BasketId or BasketItem Id is null or empty.")
	}
	if err := r.service.DeleteItem(g.Request.Context(), id, itemId); err != nil {
		g.JSON(http.StatusBadRequest, err.Error())
	}
	g.JSON(http.StatusAccepted, "")
}

type (
	CreateBasketRequest struct {
		Buyer string `json:"buyer" validate:"required"`
	}

	AddItemRequest struct {
		BasketId string  `json:"basketId"  validate:"required"`
		Sku      string  `json:"sku"  validate:"required"`
		Quantity int     `json:"quantity" validate:"required,gte=0,lte=20"`
		Price    float64 `json:"price" validate:"required,gte=0"`
	}
)
